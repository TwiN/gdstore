# gdstore

**gdstore**, short for **G**o **D**isk store, is a thread-safe (goroutine-safe) key-value library in Go for 
persisting data to disk.

This library does not have speed as its main purpose, but rather, ease of use.
As such, [the configuration required](#usage) is minimal.

If you're looking for a high-performance key-value store/database/cache, there are definitely better alternatives, but if you're searching for a simple way to persist key-value entries to disk, then this is definitely what you're looking for.


## Table of Contents

- [Motivation](#motivation)
- [Features](#features)
- [Usage](#usage)
    - [Write](#write)
    - [Read](#read)
    - [Delete](#delete)
- [FAQ](#faq)
    - [How is data persisted?](#how-is-data-persisted)


## Motivation

Why does this library exist? Because the numerous other options that currently exists
were very overkill for simple use cases.


## Features

The main features are as follow:
- **Simple to use**
- **In-memory**
- **Persistence**: Every entry is persisted to a single file
- **Thread safe/Goroutine safe**: You can call the same store concurrently


## Usage

```go
store := gdstore.New("store.db")
defer store.Close()
```

**NOTE:** You do not have to close the store every time you write in it. Also, the store is automatically opened on write. Closing a store that is already closed has no effect.


### Write

```go
err := store.Put("key", []byte("value"))
```

If prefer to write in bulk, you can use `PutAll`.

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

While the data is always persisted on disk, the data is also stored in-memory, so read operations are fast.


### Delete

```go
err := store.Delete("key")
```


## FAQ

### How is data persisted?

In order to improve the write speed as much as possible, rather than storing the data 
that currently exists in the store, the _actions_ are stored incrementally.

For instance, say you're creating a key `john` with the value `100` and another key `bob` with the value `500`, and then deleting the key `john`, the file would contain the following events:
```
SET john 100
SET bob 500
DEL john
``` 

On one hand, this has the advantage of not requiring to search in the file for the key `john` and then removing it, which could take an extended period of time, or re-creating a new file with the current data.

On the other hand, if you're storing a large amount of data, depending on your use case, the file where the data is persisted could become large.

With this in mind, an extra function called `Consolidate` has been introduced to the `GDStore` struct.

`Consolidate` combines all entries recorded in the file and re-saves only the necessary entries. Using the same example
as earlier, before `Consolidate` is executed, you'd have a file that would look like this:

```
SET john 100
SET bob 500
DEL john
```

After `Consolidate` is executed, your file would look like the following:

```
SET bob 500
```

This function is automatically executed every time a store is loaded (through `gdstore.New(...)`), but can be manually 
called if necessary.

The reason why it's not periodically executed in the background is because this library should fit both major use cases
for a persistent map, which are:
- **Long-lived**: You need to perform operations over a long period of time. A good use case would be a web server that needs to store some data.
- **Short-lived**: You need to perform operations over a short period of time. A good use case would be a CLI application.

A long-lived application that receives a constant stream of requests which may leverage gdstore can benefit from 
periodically executing `Consolidate` to reduce the size of the store file in the long run, but a short-lived 
application like a CLI tool might not benefit from it.
