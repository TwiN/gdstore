package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io/ioutil"
	"os"
)

var (
	ErrNotFound = errors.New("no data was found for the given key")
)

type GDStore struct {
	FilePath string
	data     map[string][]byte
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

func (store *GDStore) Get(key string) ([]byte, error) {
	if value, ok := store.data[key]; ok {
		return value, nil
	}
	return nil, ErrNotFound
}

func (store *GDStore) Put(key string, value []byte) error {
	store.data[key] = value
	return store.saveToDisk()
}

func (store *GDStore) Delete(key string) error {
	delete(store.data, key)
	return store.saveToDisk()
}

func (store *GDStore) Count() int {
	return len(store.data)
}

func (store *GDStore) saveToDisk() error {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(store.data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(store.FilePath, buffer.Bytes(), os.ModePerm)
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
