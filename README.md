# gdstore

**gdstore**, short for **G**o **D**isk store, is a key-value library in Go for persisting data to disk.

This library does not have speed as its main purpose, but rather, ease of use.

If you're looking for a high-performance key-value store/database/cache, there are definitely better
alternatives, but if you're searching for a simple way to persist key-value entries to disk, then 
this is definitely what you're looking for.



## Table of Contents

- [Motivation](#motivation)
- [Usage](#usage)
    - [Write](#write)
    - [Read](#read)
    - [Delete](#delete)


## Motivation

Why does this library exist? Because the numerous other options that currently exists
were very overkill for simple use cases.


## Usage

```go
store := New("store.db")
```


### Write

```go
err := store.Put("key", []byte("value"))
```

Writes are slow; if you can write in bulk, use `PutAll` instead.

```go
entries := map[string][]byte{
	"1": []byte("apple"),
	"2": []byte("banana"),
	"3": []byte("orange"),
}
err := store.PutAll(entries)
```


### Read

```go
value, exists := store.Get("key")
```

Even though writes are slow, the data is stored in-memory, so read operations are fast.


### Delete

```go
err := store.Delete("key")
```
