package gdstore

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

const (
	TestStoreFile = "gdstore.data"
)

func TestNew(t *testing.T) {
	store := New(TestStoreFile)
	defer deleteTestStoreFile()
	if store == nil {
		t.Error("Store shouldn't have returned nil")
	} else {
		store.Close()
	}
}

func TestNewWithExistingStoreFile(t *testing.T) {
	// Create a store file
	store := New(TestStoreFile)
	defer deleteTestStoreFile()
	if store.Count() != 0 {
		t.Errorf("Expected to have 0 entries, but got %d instead", store.Count())
	}
	_ = store.Put("test1", []byte("..."))
	_ = store.Put("test2", []byte("..."))
	_ = store.Put("test3", []byte("..."))
	_ = store.Delete("test3")
	store.Close()

	// Check if the previous store was persisted to the file
	store = New(TestStoreFile)
	if store.Count() != 2 {
		t.Errorf("Expected to have 2 entries, but got %d instead", store.Count())
	}
	store.Close()
}

func TestGDStore_Count(t *testing.T) {
	store := New(TestStoreFile)
	defer deleteTestStoreFile()
	_ = store.Put("test1", []byte("hello"))
	_ = store.Put("test2", []byte("hey"))
	_ = store.Put("test3", []byte("hi"))
	numberOfEntries := store.Count()
	expectedNumberOfEntries := 3
	if numberOfEntries != expectedNumberOfEntries {
		t.Errorf("Expected to have %d entries, but got %d", expectedNumberOfEntries, numberOfEntries)
	}
	store.Close()
}

func TestGDStore_Put(t *testing.T) {
	store := New(TestStoreFile)
	defer deleteTestStoreFile()
	_ = store.Put("key", []byte("value"))
	checkValueForKey(t, store, "key", []byte("value"))
	store.Close()
}

func TestGDStore_PutMultiple(t *testing.T) {
	store := New(TestStoreFile)
	defer deleteTestStoreFile()
	_ = store.Put("test1", []byte("hello"))
	checkValueForKey(t, store, "test1", []byte("hello"))
	_ = store.Put("test2", []byte("hey"))
	checkValueForKey(t, store, "test2", []byte("hey"))
	_ = store.Put("test3", []byte("hi"))
	checkValueForKey(t, store, "test3", []byte("hi"))
	store.Close()
}

func TestGDStore_PutNilValue(t *testing.T) {
	store := New(TestStoreFile)
	defer deleteTestStoreFile()
	_ = store.Put("test", nil)
	checkValueForKey(t, store, "test", nil)
	store.Close()
}

//func TestPutPerf(t *testing.T) {
//	store := New(TestStoreFile)
//	defer deleteTestStoreFile()
//	start := time.Now()
//	for i := 0; i < 1000; i++ {
//		_ = store.Put(fmt.Sprintf("test_%d", i), []byte("hello"))
//	}
//	end := time.Since(start)
//	t.Logf("Took %s", end)
//  store.Close()
//}

func TestGDStore_PutAll(t *testing.T) {
	store := New(TestStoreFile)
	defer deleteTestStoreFile()
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
	store.Close()
}

//func TestGDStore_PutAllPerf(t *testing.T) {
//	store := New(TestStoreFile)
//	defer deleteTestStoreFile()
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

func TestGDStore_PutThenDelete(t *testing.T) {
	store := New(TestStoreFile)
	defer deleteTestStoreFile()
	checkKeyNotExists(t, store, "key")
	_ = store.Put("key", []byte("value"))
	checkValueForKey(t, store, "key", []byte("value"))
	_ = store.Delete("key")
	checkKeyNotExists(t, store, "key")
	store.Close()
}

func TestGDStore_PutConcurrent(t *testing.T) {
	store := New(TestStoreFile)
	defer deleteTestStoreFile()
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			_ = store.Put(fmt.Sprintf("k%d", i), nil)
		}(i)
	}
	wg.Wait()
	store.Close()
}

func TestGDStore_Keys(t *testing.T) {
	store := New(TestStoreFile)
	defer deleteTestStoreFile()
	_ = store.Put("1", nil)
	_ = store.Put("2", nil)
	_ = store.Put("3", nil)
	keys := store.Keys()
	if len(keys) != 3 {
		t.Errorf("[%s] Expected 3 keys, got %d instead", t.Name(), len(keys))
	}
	for _, key := range keys {
		if key != "1" && key != "2" && key != "3" {
			t.Errorf("[%s] Expected keys to be '1', '2' or '3', but got '%s' instead", t.Name(), key)
		}
	}
	store.Close()
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

func deleteTestStoreFile() {
	_ = os.Remove(TestStoreFile)
	_ = os.Remove(fmt.Sprintf("%s.bak", TestStoreFile))
}
