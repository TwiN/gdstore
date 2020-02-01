package main

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
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
	return store.saveToDisk()
}

func (store *GDStore) PutAll(entries map[string][]byte) error {
	store.mux.Lock()
	defer store.mux.Unlock()
	for key, value := range entries {
		store.data[key] = value
	}
	return store.saveToDisk()
}

func (store *GDStore) Delete(key string) error {
	store.mux.Lock()
	defer store.mux.Unlock()
	delete(store.data, key)
	return store.saveToDisk()
}

func (store *GDStore) Count() int {
	store.mux.Lock()
	defer store.mux.Unlock()
	return len(store.data)
}

func (store *GDStore) saveToDisk() error {
	// TODO: Save entry by entry instead of saving the entire map object?
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(store.data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(store.FilePath, buffer.Bytes(), 0644)
}

func (store *GDStore) loadFromDisk() error {
	b, err := ioutil.ReadFile(store.FilePath)
	if os.IsNotExist(err) {
		// A new file will be created on first persist, so we can ignore this error
		store.data = make(map[string][]byte)
		return nil
	} else {
		buffer := new(bytes.Buffer)
		buffer.Write(b)
		decoder := gob.NewDecoder(buffer)
		return decoder.Decode(&store.data)
	}
}
