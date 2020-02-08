package main

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Action string

var (
	ActionPut    Action = "SET"
	ActionDelete Action = "DEL"

	ErrCannotDecodeElement = errors.New("failed to decode element")
)

type Entry struct {
	Action Action
	Key    string
	Value  []byte
}

func (e Entry) ToLine() []byte {
	encodedKeyChannel := make(chan string)
	encodedValueChannel := make(chan string)
	go func() {
		encodedKeyChannel <- base64.StdEncoding.EncodeToString([]byte(e.Key))
	}()
	go func() {
		encodedValueChannel <- base64.StdEncoding.EncodeToString(e.Value)
	}()
	return []byte(fmt.Sprintf("%s,%s,%s\n", e.Action, <-encodedKeyChannel, <-encodedValueChannel))
}

func newEntry(action Action, key string, value []byte) *Entry {
	return &Entry{
		Action: action,
		Key:    key,
		Value:  value,
	}
}

func newBulkEntries(action Action, keyValues map[string][]byte) []*Entry {
	var entries []*Entry
	for k, v := range keyValues {
		entries = append(entries, &Entry{
			Action: action,
			Key:    k,
			Value:  v,
		})
	}
	return entries
}

func newEntryFromLine(line string) (*Entry, error) {
	elements := strings.Split(line, ",")

	keyAsBytes, err := base64.StdEncoding.DecodeString(elements[1])
	if err != nil {
		return nil, ErrCannotDecodeElement
	}
	key := string(keyAsBytes)

	value, err := base64.StdEncoding.DecodeString(elements[2])
	if err != nil {
		return nil, ErrCannotDecodeElement
	}

	return &Entry{
		Action: Action(elements[0]),
		Key:    key,
		Value:  value,
	}, nil
}

// appendEntry appends an action to the file
func (store GDStore) appendEntry(entry *Entry) error {
	return store.appendEntries([]*Entry{entry})
}

func (store GDStore) appendEntries(entries []*Entry) error {
	file, err := os.OpenFile(store.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, entry := range entries {
		_, err = file.Write(entry.ToLine())
		if err != nil {
			return err
		}
	}
	return nil
}

// loadFromDisk loads the store from the disk and consolidates the entries, or creates an empty file if there is no file
func (store *GDStore) loadFromDisk() error {
	store.data = make(map[string][]byte)
	file, err := os.Open(store.FilePath)
	if os.IsNotExist(err) {
		file, err := os.Create(store.FilePath)
		if err != nil {
			return err
		}
		return file.Close()
	}
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
	// Close it now - even if defer will close it again.
	// The file needs to be closed for store.Consolidate()
	_ = file.Close()
	return store.Consolidate()
}
