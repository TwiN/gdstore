package gdstore

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

const (
	TestStoreFile = "gdstore.data"
)

func TestNew(t *testing.T) {
	// Make sure the test store file doesn't exist
	_, err := os.Stat(TestStoreFile)
	if !os.IsNotExist(err) {
		t.Error("Store file shouldn't exist yet")
	}

	// Create a new store
	store := New(TestStoreFile)
	defer deleteTestStoreFile()
	if store == nil {
		t.Error("Store shouldn't have returned nil")
	} else {
		store.Close()
	}

	// Check if the test store file exists
	_, err = os.Stat(TestStoreFile)
	if os.IsNotExist(err) {
		t.Error("Store file should exist")
	}
}

func TestNewWithExistingStoreFile(t *testing.T) {
	// Create a store file
	store := New(TestStoreFile)
	defer deleteTestStoreFile()
	if store.Count() != 0 {
		t.Errorf("Expected to have 0 entries, but got %d instead", store.Count())
	}
	_ = store.Put("test1", []byte("..."))
	_ = store.Put("test2", []byte("..."))
	_ = store.Put("test3", []byte("..."))
	_ = store.Delete("test3")
	store.Close()

	// Check if the previous store was persisted to the file
	store = New(TestStoreFile)
	if store.Count() != 2 {
		t.Errorf("Expected to have 2 entries, but got %d instead", store.Count())
	}
	store.Close()
}

func TestGDStore_Count(t *testing.T) {
	store := New(TestStoreFile)
	defer deleteTestStoreFile()
	_ = store.Put("test1", []byte("hello"))
	_ = store.Put("test2", []byte("hey"))
	_ = store.Put("test3", []byte("hi"))
	numberOfEntries := store.Count()
	expectedNumberOfEntries := 3
	if numberOfEntries != expectedNumberOfEntries {
		t.Errorf("Expected to have %d entries, but got %d", expectedNumberOfEntries, numberOfEntries)
	}
	store.Close()
}

func TestGDStore_Put(t *testing.T) {
	store := New(TestStoreFile)
	defer deleteTestStoreFile()
	_ = store.Put("key", []byte("value"))
	checkValueForKey(t, store, "key", []byte("value"))
	store.Close()
}

func TestGDStore_PutMultiple(t *testing.T) {
	store := New(TestStoreFile)
	defer deleteTestStoreFile()
	_ = store.Put("test1", []byte("hello"))
	checkValueForKey(t, store, "test1", []byte("hello"))
	_ = store.Put("test2", []byte("hey"))
	checkValueForKey(t, store, "test2", []byte("hey"))
	_ = store.Put("test3", []byte("hi"))
	checkValueForKey(t, store, "test3", []byte("hi"))
	store.Close()
}

func TestGDStore_PutNilValue(t *testing.T) {
	store := New(TestStoreFile)
	defer deleteTestStoreFile()
	_ = store.Put("test", nil)
	checkValueForKey(t, store, "test", nil)
	store.Close()
}

//func TestPutPerf(t *testing.T) {
//	store := New(TestStoreFile)
//	defer deleteTestStoreFile()
//	start := time.Now()
//	for i := 0; i < 1000; i++ {
//		_ = store.Put(fmt.Sprintf("test_%d", i), []byte("hello"))
//	}
//	end := time.Since(start)
//	t.Logf("Took %s", end)
//  store.Close()
//}

func TestGDStore_PutAll(t *testing.T) {
	store := New(TestStoreFile)
	defer deleteTestStoreFile()
	entries := map[string][]byte{
		"1": []byte("apple"),
		"2": []byte("banana"),
		"3": []byte("orange"),
	}
	checkKeyNotExists(t, store, "1")
	checkKeyNotExists(t, store, "2")
	checkKeyNotExists(t, store, "3")
	_ = store.PutAll(entries)
	checkValueForKey(t, store, "1", []byte("apple"))
	checkValueForKey(t, store, "2", []byte("banana"))
	checkValueForKey(t, store, "3", []byte("orange"))
	store.Close()
}

//func TestGDStore_PutAllPerf(t *testing.T) {
//	store := New(TestStoreFile)
//	defer deleteTestStoreFile()
//	entries := make(map[string][]byte)
//	for i := 0; i < 10000; i++ {
//		entries[fmt.Sprintf("test_%d", i)] = []byte("hello")
//	}
//	start := time.Now()
//	_ = store.PutAll(entries)
//	end := time.Since(start)
//	if end > 500 * time.Millisecond {
//		t.Errorf("Took too long (%s)", end)
//	}
//}

func TestGDStore_PutThenDelete(t *testing.T) {
	store := New(TestStoreFile)
	defer deleteTestStoreFile()
	checkKeyNotExists(t, store, "key")
	_ = store.Put("key", []byte("value"))
	checkValueForKey(t, store, "key", []byte("value"))
	_ = store.Delete("key")
	checkKeyNotExists(t, store, "key")
	store.Close()
}

func TestGDStore_PutConcurrent(t *testing.T) {
	store := New(TestStoreFile)
	defer deleteTestStoreFile()
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			_ = store.Put(fmt.Sprintf("k%d", i), nil)
		}(i)
	}
	wg.Wait()
	store.Close()
}

func TestGDStore_Keys(t *testing.T) {
	store := New(TestStoreFile)
	defer deleteTestStoreFile()
	_ = store.Put("1", nil)
	_ = store.Put("2", nil)
	_ = store.Put("3", nil)
	keys := store.Keys()
	if len(keys) != 3 {
		t.Errorf("[%s] Expected 3 keys, got %d instead", t.Name(), len(keys))
	}
	for _, key := range keys {
		if key != "1" && key != "2" && key != "3" {
			t.Errorf("[%s] Expected keys to be '1', '2' or '3', but got '%s' instead", t.Name(), key)
		}
	}
	store.Close()
}

func TestGDStore_Values(t *testing.T) {
	store := New(TestStoreFile)
	defer deleteTestStoreFile()
	_ = store.Put("1", []byte("one"))
	_ = store.Put("2", []byte("two"))
	_ = store.Put("3", []byte("three"))
	values := store.Values()
	if len(values) != 3 {
		t.Errorf("[%s] Expected 3 values, got %d instead", t.Name(), len(values))
	}
	for _, value := range values {
		valueAsString := string(value)
		if valueAsString != "one" && valueAsString != "two" && valueAsString != "three" {
			t.Errorf("[%s] Expected values to be '1', '2' or '3', but got '%s' instead", t.Name(), value)
		}
	}
	store.Close()
}

func TestGDStore_GetInt(t *testing.T) {
	store := New(TestStoreFile)
	defer deleteTestStoreFile()
	_ = store.Put("test", []byte("42"))
	number, ok, err := store.GetInt("test")
	if !ok || err != nil {
		t.Errorf("[%s] Expected key 'test' to exist and not be return an error", t.Name())
	}
	if number != 42 {
		t.Errorf("[%s] Expected key 'test' to have value '42', but got '%d' instead", t.Name(), number)
	}
	store.Close()
}

func TestGDStore_GetIntWithNegativeNumber(t *testing.T) {
	store := New(TestStoreFile)
	defer deleteTestStoreFile()
	_ = store.Put("test", []byte("-42"))
	number, ok, err := store.GetInt("test")
	if !ok || err != nil {
		t.Errorf("[%s] Expected key 'test' to exist and not be return an error", t.Name())
	}
	if number != -42 {
		t.Errorf("[%s] Expected key 'test' to have value '-42', but got '%d' instead", t.Name(), number)
	}
	store.Close()
}

func TestGDStore_GetIntWithNonInt(t *testing.T) {
	store := New(TestStoreFile)
	defer deleteTestStoreFile()
	_ = store.Put("test", []byte("NaN"))
	_, ok, err := store.GetInt("test")
	if !ok {
		t.Errorf("[%s] Expected key 'test' to exist", t.Name())
	}
	if err == nil {
		t.Errorf("[%s] Expected key 'test' to return an error because the value is not an int", t.Name())
	}
	store.Close()
}

////////////////
// BENCHMARKS //
////////////////

func BenchmarkGDStore_Put(b *testing.B) {
	store := New(TestStoreFile)
	defer deleteTestStoreFile()
	for n := 0; n < b.N; n++ {
		_ = store.Put(fmt.Sprintf("test_%d", n), []byte("value"))
	}
	store.Close()
}

func BenchmarkGDStore_PutWithLargeValue(b *testing.B) {
	store := New(TestStoreFile)
	defer deleteTestStoreFile()
	for n := 0; n < b.N; n++ {
		_ = store.Put(fmt.Sprintf("test_%d", n), []byte("large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value"))
	}
	store.Close()
}

func BenchmarkGDStore_PutWithBuffer(b *testing.B) {
	store := New(TestStoreFile).WithBuffer(true)
	defer deleteTestStoreFile()
	for n := 0; n < b.N; n++ {
		_ = store.Put(fmt.Sprintf("test_%d", n), []byte("value"))
	}
	store.Close()
}

func BenchmarkGDStore_PutWithBufferAndLargeValue(b *testing.B) {
	store := New(TestStoreFile).WithBuffer(true)
	defer deleteTestStoreFile()
	for n := 0; n < b.N; n++ {
		_ = store.Put(fmt.Sprintf("test_%d", n), []byte("large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value_large_value"))
	}
	store.Close()
}

//func BenchmarkMap(b *testing.B) {
//	m := make(map[string][]byte)
//	for n := 0; n < b.N; n++ {
//		m[fmt.Sprintf("test_%d", n)] = []byte("value")
//	}
//}

///////////////////////
// Utility functions //
///////////////////////

func checkValueForKey(t *testing.T, store *GDStore, key string, expectedValue []byte) {
	value, exists := store.Get(key)
	if !exists {
		t.Errorf("[%s] Expected key '%s' to exist", t.Name(), key)
	}
	if string(value) != string(expectedValue) {
		t.Errorf("[%s] Expected key '%s' to have value '%v', but had '%v' instead", t.Name(), key, expectedValue, value)
	}
}

func checkKeyNotExists(t *testing.T, store *GDStore, key string) {
	if _, exists := store.Get(key); exists {
		t.Errorf("[%s] Expected key '%s' to not exist", t.Name(), key)
	}
}

func deleteTestStoreFile() {
	_ = os.Remove(TestStoreFile)
	_ = os.Remove(fmt.Sprintf("%s.bak", TestStoreFile))
}
