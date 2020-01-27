package main

import "errors"

var (
	ErrNotFound = errors.New("no data was found for the given key")
)

type GDStore struct {
	FilePath string
}

func New(filePath string) *GDStore {
	return &GDStore{
		FilePath: filePath,
	}
}

func (store GDStore) Get(key string) ([]byte, error) {
	return nil, ErrNotFound
}

func (store GDStore) Put(key string, value []byte) error {
	return nil
}

func (store GDStore) Delete(key string) error {
	return nil
}

func (store GDStore) Count() int {
	return 0
}
