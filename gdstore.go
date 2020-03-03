package gdstore

import (
	"os"
	"strconv"
	"sync"
)

type GDStore struct {
	FilePath string
	file     *os.File
	data     map[string][]byte
	mux      sync.Mutex
}

// New creates a new GDStore
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

// Get returns the value of a key as well as a bool that indicates whether an entry exists for that key.
// The bool is particularly useful if you want to differentiate between a key that has a nil value, and a
// key that doesn't exist
func (store *GDStore) Get(key string) (value []byte, ok bool) {
	value, ok = store.data[key]
	return
}

// GetString does the same thing as Get, but converts the value to a string
func (store *GDStore) GetString(key string) (valueAsString string, ok bool) {
	var value []byte
	value, ok = store.data[key]
	if ok {
		valueAsString = string(value)
	}
	return
}

// GetInt does the same thing as Get, but converts the value to an int
func (store *GDStore) GetInt(key string) (valueAsInt int, ok bool, err error) {
	var value string
	value, ok = store.GetString(key)
	if ok {
		valueAsInt, err = strconv.Atoi(value)
	}
	return
}

// Put creates an entry or updates the value of an existing key
func (store *GDStore) Put(key string, value []byte) error {
	store.mux.Lock()
	defer store.mux.Unlock()
	store.data[key] = value
	return store.appendEntryToFile(newEntry(ActionPut, key, value))
}

// PutAll creates or updates a map of entries
func (store *GDStore) PutAll(entries map[string][]byte) error {
	store.mux.Lock()
	defer store.mux.Unlock()
	for key, value := range entries {
		store.data[key] = value
	}
	return store.appendEntriesToFile(newBulkEntries(ActionPut, entries))
}

// Delete removes a key from the store
func (store *GDStore) Delete(key string) error {
	store.mux.Lock()
	defer store.mux.Unlock()
	delete(store.data, key)
	return store.appendEntryToFile(newEntry(ActionDelete, key, nil))
}

// Count returns the total number of entries in the store
func (store *GDStore) Count() int {
	store.mux.Lock()
	defer store.mux.Unlock()
	return len(store.data)
}

// Keys returns a list of all keys
func (store *GDStore) Keys() []string {
	store.mux.Lock()
	defer store.mux.Unlock()
	keys := make([]string, len(store.data))
	i := 0
	for k := range store.data {
		keys[i] = k
		i++
	}
	return keys
}

// Values returns a list of all values
func (store *GDStore) Values() [][]byte {
	store.mux.Lock()
	defer store.mux.Unlock()
	values := make([][]byte, len(store.data))
	i := 0
	for _, v := range store.data {
		values[i] = v
		i++
	}
	return values
}

// Close closes the store's file if it isn't already closed.
// Note that any write actions, such as the usage of Put and PutAll, will automatically re-open the store.
func (store *GDStore) Close() {
	if store.file != nil {
		err := store.file.Close()
		store.file = nil
		if err != nil {
			panic(err)
		}
	}
}
