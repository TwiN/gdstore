package main

import (
	"os"
	"sync"
)

type GDStore struct {
	FilePath string
	file     *os.File
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
	return store.appendEntryToFile(newEntry(ActionPut, key, value))
}

func (store *GDStore) PutAll(entries map[string][]byte) error {
	store.mux.Lock()
	defer store.mux.Unlock()
	for key, value := range entries {
		store.data[key] = value
	}
	return store.appendEntriesToFile(newBulkEntries(ActionPut, entries))
}

func (store *GDStore) Delete(key string) error {
	store.mux.Lock()
	defer store.mux.Unlock()
	delete(store.data, key)
	return store.appendEntryToFile(newEntry(ActionDelete, key, nil))
}

func (store *GDStore) Count() int {
	store.mux.Lock()
	defer store.mux.Unlock()
	return len(store.data)
}

func (store *GDStore) Close() {
	if store.file != nil {
		err := store.file.Close()
		store.file = nil
		if err != nil {
			panic(err)
		}
	}
}
