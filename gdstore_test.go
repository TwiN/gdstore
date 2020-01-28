package main

import (
	"os"
	"testing"
)

const (
	TestStoreFile = "gdstore.data"
)

func TestNew(t *testing.T) {
	store := New(TestStoreFile)
	defer deleteStoreFile(store)
	if store == nil {
		t.Error("Store shouldn't have returned nil")
	}
}

func TestNewWithExistingStoreFile(t *testing.T) {
	// Create a store file
	store := New(TestStoreFile)
	defer deleteStoreFile(store)
	if store.Count() != 0 {
		t.Errorf("Expected to have 0 entries, but got %d instead", store.Count())
	}
	_ = store.Put("test1", []byte("..."))
	_ = store.Put("test2", []byte("..."))
	// Check if the previous store was persisted to the file
	store = New(TestStoreFile)
	if store.Count() != 2 {
		t.Errorf("Expected to have 2 entries, but got %d instead", store.Count())
	}
}

func TestCount(t *testing.T) {
	store := New(TestStoreFile)
	defer deleteStoreFile(store)
	_ = store.Put("test1", []byte("hello"))
	_ = store.Put("test2", []byte("hey"))
	_ = store.Put("test3", []byte("hi"))
	numberOfEntries := store.Count()
	expectedNumberOfEntries := 3
	if numberOfEntries != expectedNumberOfEntries {
		t.Errorf("Expected to have %d entries, but got %d", expectedNumberOfEntries, numberOfEntries)
	}
}

func TestPut(t *testing.T) {
	store := New(TestStoreFile)
	defer deleteStoreFile(store)
	_ = store.Put("test1", []byte("hello"))
	checkValueForKey(t, store, "test1", []byte("hello"))
	_ = store.Put("test2", []byte("hey"))
	checkValueForKey(t, store, "test2", []byte("hey"))
	_ = store.Put("test3", []byte("hi"))
	checkValueForKey(t, store, "test3", []byte("hi"))
}

func TestPutNilValue(t *testing.T) {
	store := New(TestStoreFile)
	defer deleteStoreFile(store)
	_ = store.Put("test", nil)
	checkValueForKey(t, store, "test", nil)
}

//func TestPutPerf(t *testing.T) {
//	store := New(TestStoreFile)
//	defer deleteStoreFile(store)
//	start := time.Now()
//	for i := 0; i < 1000; i++ {
//		_ = store.Put(fmt.Sprintf("test_%d", i), []byte("hello"))
//	}
//	end := time.Since(start)
//	t.Errorf("Took %s", end)
//}

func TestPutAll(t *testing.T) {
	store := New(TestStoreFile)
	defer deleteStoreFile(store)
	entries := map[string][]byte{
		"1": []byte("apple"),
		"2": []byte("banana"),
		"3": []byte("orange"),
	}
	checkKeyNotExists(t, store, "1")
	checkKeyNotExists(t, store, "2")
	checkKeyNotExists(t, store, "3")
	_ = store.PutAll(entries)
	checkValueForKey(t, store, "1", []byte("apple"))
	checkValueForKey(t, store, "2", []byte("banana"))
	checkValueForKey(t, store, "3", []byte("orange"))
}

//func TestPutAllPerf(t *testing.T) {
//	store := New(TestStoreFile)
//	defer deleteStoreFile(store)
//	entries := make(map[string][]byte)
//	for i := 0; i < 10000; i++ {
//		entries[fmt.Sprintf("test_%d", i)] = []byte("hello")
//	}
//	start := time.Now()
//	_ = store.PutAll(entries)
//	end := time.Since(start)
//	if end > 500 * time.Millisecond {
//		t.Errorf("Took too long (%s)", end)
//	}
//}

func TestPutThenDelete(t *testing.T) {
	store := New(TestStoreFile)
	defer deleteStoreFile(store)
	checkKeyNotExists(t, store, "key")
	_ = store.Put("key", []byte("value"))
	checkValueForKey(t, store, "key", []byte("value"))
	_ = store.Delete("key")
	checkKeyNotExists(t, store, "key")
}

func checkValueForKey(t *testing.T, store *GDStore, key string, expectedValue []byte) {
	value, exists := store.Get(key)
	if !exists {
		t.Errorf("[%s] Expected key '%s' to exist", t.Name(), key)
	}
	if string(value) != string(expectedValue) {
		t.Errorf("[%s] Expected key '%s' to have value '%v', but had '%v' instead", t.Name(), key, expectedValue, value)
	}
}

func checkKeyNotExists(t *testing.T, store *GDStore, key string) {
	if _, exists := store.Get(key); exists {
		t.Errorf("[%s] Expected key '%s' to not exist", t.Name(), key)
	}
}

func deleteStoreFile(store *GDStore) {
	_ = os.Remove(store.FilePath)
}
