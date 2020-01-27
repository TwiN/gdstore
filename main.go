package main

import "errors"

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
	store.loadFromDisk()
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
	return nil
}

func (store *GDStore) loadFromDisk() {

}
