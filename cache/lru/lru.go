package lru

import (
	"container/list"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Cache is the base struct to store elements
type Cache struct {
	mu sync.Mutex

	// list & table of *entry objects
	list  *list.List
	table map[string]*list.Element

	// Our current size, in bytes. Obviously a gross simplification and low-grade
	// approximation.
	size uint64

	// How many bytes we are limiting the cache to.
	capacity uint64
}

// Value that go into LRUCache need to satisfy this interface.
type Value interface {
	Size() int
}

// Item that store key and value
type Item struct {
	Key   string
	Value Value
}

type entry struct {
	key          string
	value        Value
	size         int
	timeAccessed time.Time
}

// NewCache returns LRUCache with a capacity set
func NewCache(capacity uint64) *Cache {
	return &Cache{
		list:     list.New(),
		table:    make(map[string]*list.Element),
		capacity: capacity,
	}
}

// Get item with key
func (lru *Cache) Get(key string) (v Value, ok bool) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	element := lru.table[key]
	if element == nil {
		return nil, false
	}
	lru.moveToFront(element)
	return element.Value.(*entry).value, true
}

// Set value with key, if exits replace it
func (lru *Cache) Set(key string, value Value) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if element := lru.table[key]; element != nil {
		lru.updateInplace(element, value)
	} else {
		lru.addNew(key, value)
	}
}

// SetIfAbsent if item exits, move to front
func (lru *Cache) SetIfAbsent(key string, value Value) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if element := lru.table[key]; element != nil {
		lru.moveToFront(element)
	} else {
		lru.addNew(key, value)
	}
}

// Delete removes item with key
func (lru *Cache) Delete(key string) bool {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	element := lru.table[key]
	if element == nil {
		return false
	}

	lru.list.Remove(element)
	delete(lru.table, key)
	lru.size -= uint64(element.Value.(*entry).size)
	return true
}

// Clear clean all items
func (lru *Cache) Clear() {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	lru.list.Init()
	lru.table = make(map[string]*list.Element)
	lru.size = 0
}

// SetCapacity set capacity of lru
func (lru *Cache) SetCapacity(capacity uint64) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	lru.capacity = capacity
	lru.checkCapacity()
}

// Stats returns all fields of lru
func (lru *Cache) Stats() (length, size, capacity uint64, oldest time.Time) {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	if lastElem := lru.list.Back(); lastElem != nil {
		oldest = lastElem.Value.(*entry).timeAccessed
	}
	return uint64(lru.list.Len()), lru.size, lru.capacity, oldest
}

// StatsJSON returns json serialized of lru.
// Performance better than JSON.dumps
func (lru *Cache) StatsJSON() string {
	if lru == nil {
		return "{}"
	}
	l, s, c, o := lru.Stats()
	return fmt.Sprintf("{\"Length\": %v, \"Size\": %v, \"Capacity\": %v, \"OldestAccess\": \"%v\"}", l, s, c, o)
}

// Keys returns all of keys
func (lru *Cache) Keys() []string {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	keys := make([]string, 0, lru.list.Len())
	for e := lru.list.Front(); e != nil; e = e.Next() {
		keys = append(keys, e.Value.(*entry).key)
	}
	return keys
}

// Items returns all of items
func (lru *Cache) Items() []Item {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	items := make([]Item, 0, lru.list.Len())
	for e := lru.list.Front(); e != nil; e = e.Next() {
		v := e.Value.(*entry)
		items = append(items, Item{Key: v.key, Value: v.value})
	}
	return items
}

// SaveItems try to dumps items to io.Writer
func (lru *Cache) SaveItems(w io.Writer) error {
	items := lru.Items()
	encoder := gob.NewEncoder(w)
	return encoder.Encode(items)
}

// SaveItemsToFile try to dumps items to disk
func (lru *Cache) SaveItemsToFile(path string) error {
	var (
		wr  *os.File
		err error
	)
	if wr, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644); err != nil {
		return err
	}
	defer wr.Close()
	return lru.SaveItems(wr)
}

// LoadItems loads items from io.Reader
func (lru *Cache) LoadItems(r io.Reader) error {
	items := make([]Item, 0)
	decoder := gob.NewDecoder(r)
	if err := decoder.Decode(&items); err != nil {
		return err
	}

	lru.mu.Lock()
	defer lru.mu.Unlock()
	for _, item := range items {
		// TODO: copied from Set()
		if element := lru.table[item.Key]; element != nil {
			lru.updateInplace(element, item.Value)
		} else {
			lru.addNew(item.Key, item.Value)
		}
	}

	return nil
}

// LoadItemsFromFile loads cache data from disk
func (lru *Cache) LoadItemsFromFile(path string) error {
	var (
		rd  *os.File
		err error
	)
	if rd, err = os.Open(path); err != nil {
		return err
	}
	defer rd.Close()
	return lru.LoadItems(rd)
}

func (lru *Cache) updateInplace(element *list.Element, value Value) {
	valueSize := value.Size()
	sizeDiff := valueSize - element.Value.(*entry).size
	element.Value.(*entry).value = value
	element.Value.(*entry).size = valueSize
	lru.size += uint64(sizeDiff)
	lru.moveToFront(element)
	lru.checkCapacity()
}

func (lru *Cache) moveToFront(element *list.Element) {
	lru.list.MoveToFront(element)
	element.Value.(*entry).timeAccessed = time.Now()
}

func (lru *Cache) addNew(key string, value Value) {
	newEntry := &entry{key, value, value.Size(), time.Now()}
	element := lru.list.PushFront(newEntry)
	lru.table[key] = element
	lru.size += uint64(newEntry.size)
	lru.checkCapacity()
}

func (lru *Cache) checkCapacity() {
	// Partially duplicated from Delete
	for lru.size > lru.capacity {
		delElem := lru.list.Back()
		delValue := delElem.Value.(*entry)
		lru.list.Remove(delElem)
		delete(lru.table, delValue.key)
		lru.size -= uint64(delValue.size)
	}
}
