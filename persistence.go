package gdstore

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

// Consolidate combines all entries recorded in the file and re-saves only the necessary entries.
// The function is executed on creation, but can also be executed manually if storage space is a concern.
// The original file is backed up.
func (store *GDStore) Consolidate() error {
	store.mux.Lock()
	defer store.mux.Unlock()
	// Close the file because we need to rename it
	store.Close()
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
	err = file.Close()
	if err != nil {
		panic(err)
	}
	// Close store AFTER appending all entries to the new file (hence defer)
	// to make sure all the data is definitely in the new file
	defer store.Close()
	return store.appendEntriesToFile(newBulkEntries(ActionPut, store.data))
}

// loadFromDisk loads the store from the disk and consolidates the entries, or creates an empty file if there is no file
func (store *GDStore) loadFromDisk() error {
	store.data = make(map[string][]byte)
	file, err := os.Open(store.FilePath)
	if err != nil {
		// Check if the file exists, if it doesn't, then create it and return.
		if os.IsNotExist(err) {
			file, err := os.Create(store.FilePath)
			if err != nil {
				return err
			}
			return file.Close()
		} else {
			return err
		}
	}

	// File doesn't exist, so we need to read it.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entry, err := newEntryFromLine(scanner.Text())
		if err != nil {
			continue
		}
		if entry.Action == ActionPut {
			store.data[entry.Key] = entry.Value
		} else if entry.Action == ActionDelete {
			delete(store.data, entry.Key)
		}
	}
	_ = file.Close()
	return store.Consolidate()
}

// appendEntryToFile appends an entry to the store's file
func (store *GDStore) appendEntryToFile(entry *Entry) error {
	return store.appendEntriesToFile([]*Entry{entry})
}

// appendEntriesToFile appends a list of entries to the store's file
func (store *GDStore) appendEntriesToFile(entries []*Entry) (err error) {
	if store.file == nil {
		store.file, err = os.OpenFile(store.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return
		}
		store.writer = bufio.NewWriter(store.file)
	}
	for _, entry := range entries {
		if store.useBuffer {
			_, err = store.writer.Write(entry.toLine())
		} else {
			_, err = store.file.Write(entry.toLine())
		}
		if err != nil {
			return
		}
	}
	return
}
