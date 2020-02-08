package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func TestGDStore_Consolidate(t *testing.T) {
	store := New(TestStoreFile)
	defer deleteTestStoreFile()

	fileContent := getStoreFileContent(store)
	if len(fileContent) != 0 {
		t.Errorf("New store should've been empty, but instead, had: %s", fileContent)
	}

	// Add two new keys
	_ = store.Put("key1", []byte("value1"))
	_ = store.Put("key2", []byte("value2"))

	// The file should contain 2 entries, so two lines (SET)
	fileContent = getStoreFileContent(store)
	if numberOfLines := len(strings.Split(fileContent, "\n")); numberOfLines != 2 {
		t.Errorf("Store file should've had 2 lines, but had %d instead", numberOfLines)
	}

	// Delete a key that already exists
	_ = store.Delete("key1")

	// Even though a key was deleted, it still counts as an entry, so there should be 2 SET, 1 DEL, for a total of 3
	fileContent = getStoreFileContent(store)
	if numberOfLines := len(strings.Split(fileContent, "\n")); numberOfLines != 3 {
		t.Errorf("Store file should've had 3 lines, but had %d instead", numberOfLines)
	}

	_ = store.Consolidate()

	// Because we've consolidated the store, all unnecessary entries should've been removed, so since the only
	// remaining key is key2, there should be 1 entry left.
	fileContent = getStoreFileContent(store)
	if numberOfLines := len(strings.Split(fileContent, "\n")); numberOfLines != 1 {
		t.Errorf("Store file should've had 1 lines, but had %d instead", numberOfLines)
	}

	// The primary store file has been consolidated, but there should still be a backup of the old store file
	// that should still have 3 entries
	backupFileContent := getStoreBackupFileContent(store)
	if numberOfLines := len(strings.Split(backupFileContent, "\n")); numberOfLines != 3 {
		t.Errorf("Store backup file should've had 3 lines, but had %d instead", numberOfLines)
	}

	oldStoreCount := store.Count()
	oldStoreKey2Value, _ := store.Get("key2")
	store.Close()

	// Check if the new consolidated file results in a store that has the same content as the original store
	newStore := New(TestStoreFile)
	if newStore.Count() != oldStoreCount {
		t.Error("New store should have the same content as the old one")
	}
	if newStoreKey2Value, _ := newStore.Get("key2"); string(newStoreKey2Value) != string(oldStoreKey2Value) {
		t.Error("New store should have the same content as the old one")
	}
}

func getStoreFileContent(store *GDStore) string {
	store.Close()
	raw, _ := ioutil.ReadFile(store.FilePath)
	return strings.TrimSpace(string(raw))
}

func getStoreBackupFileContent(store *GDStore) string {
	store.Close()
	raw, _ := ioutil.ReadFile(fmt.Sprintf("%s.bak", store.FilePath))
	return strings.TrimSpace(string(raw))
}
