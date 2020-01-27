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

func checkValueForKey(t *testing.T, store *GDStore, key string, expectedValue []byte) {
	value, exists := store.Get(key)
	if !exists {
		t.Errorf("[%s] Expected key '%s' to exist", t.Name(), key)
	}
	if string(value) != string(expectedValue) {
		t.Errorf("[%s] Expected key '%s' to have value '%v', but had '%v' instead", t.Name(), key, expectedValue, value)
	}
}

func deleteStoreFile(store *GDStore) {
	_ = os.Remove(store.FilePath)
}
