package main

import (
	"errors"
	"fmt"
	"os"
	"sync"
)

type GDStore struct {
	FilePath string
	data     map[string][]byte
	mux      sync.Mutex
}

func New(filePath string) *GDStore {
	store := &GDStore{
		FilePath: filePath,
		data:     make(map[string][]byte),
	}
	err := store.loadFromDisk()
	if err != nil {
		panic(err)
	}
	return store
}

func (store *GDStore) Get(key string) (value []byte, ok bool) {
	value, ok = store.data[key]
	return
}

func (store *GDStore) Put(key string, value []byte) error {
	store.mux.Lock()
	defer store.mux.Unlock()
	store.data[key] = value
	return store.appendEntry(newEntry(ActionPut, key, value))
}

func (store *GDStore) PutAll(entries map[string][]byte) error {
	store.mux.Lock()
	defer store.mux.Unlock()
	for key, value := range entries {
		store.data[key] = value
	}
	return store.appendEntries(newBulkEntries(ActionPut, entries))
}

func (store *GDStore) Delete(key string) error {
	store.mux.Lock()
	defer store.mux.Unlock()
	delete(store.data, key)
	return store.appendEntry(newEntry(ActionDelete, key, nil))
}

func (store *GDStore) Count() int {
	store.mux.Lock()
	defer store.mux.Unlock()
	return len(store.data)
}

// Consolidate combines all entries recorded in the file and re-saves only the necessary entries
// The original file is backed up.
//
// The function is automatically executed on creation, but can also be ran during operation if
// storage space is a concern
func (store *GDStore) Consolidate() error {
	// Back up the old file before doing the consolidation
	err := os.Rename(store.FilePath, fmt.Sprintf("%s.bak", store.FilePath))
	if err != nil {
		return errors.New(fmt.Sprintf("unable to rename %s to %s.bak during consolidation: %s", store.FilePath, store.FilePath, err.Error()))
	}
	// Create a new empty file
	file, err := os.Create(store.FilePath)
	if err != nil {
		return errors.New(fmt.Sprintf("unable to create new empty file at %s during consolidation: %s", store.FilePath, err.Error()))
	}
	_ = file.Close()
	return store.appendEntries(newBulkEntries(ActionPut, store.data))
}
