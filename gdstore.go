package gdstore

import (
	"bufio"
	"os"
	"strconv"
	"sync"
)

type GDStore struct {
	// FilePath is the path to the file used to persist
	FilePath string

	// useBuffer lets the user define if GDStore should use a buffer, or write directly to the file.
	//
	// Writing to a buffer is much faster, but failure to close to the store (GDStore.Close()) will
	// result in a buffer that hasn't been flushed, meaning some entries may be lost.
	// You can manually flush the buffer by using GDStore.Flush().
	//
	// In contrast, writing to the file without a buffer is slower, but more reliable if your
	// application is prone to suddenly crashing
	//
	// Defaults to false
	useBuffer bool

	// persistence lets the user set whether to persist data to the file or not. Meant to be used for testing without
	// having to clean up the file.
	//
	// (default) If set to true, data is available both in-memory and persisted to the file found at FilePath.
	// If set to false, data is available only in-memory, and upon destruction, all data will be lost.
	//
	// Defaults to true
	persistence bool

	file   *os.File
	writer *bufio.Writer
	data   map[string][]byte
	mux    sync.Mutex
}

// New creates a new GDStore
func New(filePath string) *GDStore {
	store := &GDStore{
		FilePath:    filePath,
		data:        make(map[string][]byte),
		persistence: true,
	}
	err := store.loadFromDisk()
	if err != nil {
		panic(err)
	}
	return store
}

// WithBuffer sets GDStore's useBuffer parameter to the value passed as parameter
//
// The default value for useBuffer is false
func (store *GDStore) WithBuffer(useBuffer bool) *GDStore {
	store.useBuffer = useBuffer
	return store
}

// WithPersistence sets GDStore's persistence parameter to the value passed as parameter
//
// The ability to set persistence to false is there mainly for testing purposes
//
// The default value for persistence is true
func (store *GDStore) WithPersistence(persistence bool) *GDStore {
	store.persistence = persistence
	return store
}

// Get returns the value of a key as well as a bool that indicates whether an entry exists for that key.
//
// The bool is particularly useful if you want to differentiate between a key that has a nil value, and a
// key that doesn't exist
func (store *GDStore) Get(key string) (value []byte, ok bool) {
	value, ok = store.data[key]
	return
}

// GetString does the same thing as Get, but converts the value to a string
func (store *GDStore) GetString(key string) (valueAsString string, ok bool) {
	var value []byte
	value, ok = store.Get(key)
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
