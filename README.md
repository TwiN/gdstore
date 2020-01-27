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

Why does this library exist? Because a lot of other options exists, but they felt overkill
for persisting 


## Usage

```go
store := New("store.db")
```


### Write

```go
err := store.Put("key", []byte("value"))
```


### Read

```go
value, exists := store.Get("key")
```


### Delete

```go
err := store.Delete("key")
```
