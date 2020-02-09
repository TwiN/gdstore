package gdstore

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrCannotDecodeElement = errors.New("failed to decode element")
	ErrBadLine             = errors.New("bad line")
)

type Entry struct {
	Action Action
	Key    string
	Value  []byte
}

func (e Entry) toLine() []byte {
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

	if len(elements) != 3 {
		return nil, ErrBadLine
	}

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
