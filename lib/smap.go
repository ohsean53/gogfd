package lib

import (
	"sync"
)

type syncType int
const (
	RW_MUTEX syncType = 1
	CHANNEL syncType = 2
)

// A thread safe map(type: `map[int64]interface{}`).
type SharedMap interface {
	// Sets the given value under the specified key
	Set(k int64, v interface{})

	// Retrieve an item from map under given key.
	Get(k int64) (interface{}, bool)

	// Remove an item from the map.
	Remove(k int64)

	// Return the number of item within the map.
	Count() int

	// Return the map object
	Map() map[int64]interface{}

	// Return all the keys or a subset of the keys of an Map
	GetKeys() []int64
}


/**
 Using RWMutex
 */
type sharedMapRWMutex struct {
	sync.RWMutex
	m map[int64]interface{}
}

func (sm sharedMapRWMutex) Map() map[int64]interface{} {
	return sm.m
}

// Sets the given value under the specified key
func (sm sharedMapRWMutex) Set(k int64, v interface{}) {
	sm.Lock()
	sm.m[k] = v
	sm.Unlock()
}

// Retrieve an item from map under given key.
func (sm sharedMapRWMutex) Get(k int64) (interface{}, bool) {
	sm.RLock()
	defer sm.RUnlock()
	v, ok := sm.m[k]
	return v, ok
}

// Remove an item from the map.
func (sm sharedMapRWMutex) Remove(k int64) {
	sm.Lock()
	delete(sm.m, k)
	sm.Unlock()
}

// Return the number of item within the map.
func (sm sharedMapRWMutex) Count() int {
	sm.RLock()
	defer sm.RUnlock()
	return len(sm.m)
}

// Return all the keys or a subset of the keys of an Map (추가함)
func (sm sharedMapRWMutex) GetKeys() []int64 {
	sm.RLock()
	defer sm.RUnlock()
	keys := make([]int64, 0, 1024)
	for key, _ := range sm.m {
		keys = append(keys, key)
	}
	return keys
}


/**
 Using Channel
 */
type sharedMapChannel struct {
	m map[int64]interface{}
	c chan command
}

type command struct {
	action int
	key    int64
	value  interface{}
	result chan <- interface{}
}

const (
	set = iota
	get
	remove
	count
	keys
)

func (sm sharedMapChannel) Map() map[int64]interface{} {
	return sm.m
}

// Sets the given value under the specified key
func (sm sharedMapChannel) Set(k int64, v interface{}) {
	sm.c <- command{action: set, key: k, value: v}
}

// Retrieve an item from map under given key.
func (sm sharedMapChannel) Get(k int64) (interface{}, bool) {
	callback := make(chan interface{})
	sm.c <- command{action: get, key: k, result: callback}
	result := (<-callback).([2]interface{})
	return result[0], result[1].(bool)
}

// Remove an item from the map.
func (sm sharedMapChannel) Remove(k int64) {
	sm.c <- command{action: remove, key: k}
}

// Return the number of item within the map.
func (sm sharedMapChannel) Count() int {
	callback := make(chan interface{})
	sm.c <- command{action: count, result: callback}
	return (<-callback).(int)
}

// Return all the keys or a subset of the keys of an Map (추가함)
func (sm sharedMapChannel) GetKeys() []int64 {
	callback := make(chan interface{})
	sm.c <- command{action: keys, result: callback}
	return (<-callback).([]int64)
}

func (sm sharedMapChannel) run() {
	for cmd := range sm.c {
		switch cmd.action {
		case set:
			sm.m[cmd.key] = cmd.value
		case get:
			v, ok := sm.m[cmd.key]
			cmd.result <- [2]interface{}{v, ok}
		case remove:
			delete(sm.m, cmd.key)
		case count:
			cmd.result <- len(sm.m)
		case keys:
			keys := make([]int64, 0, 1024)
			for key, _ := range sm.m {
				keys = append(keys, key)
			}
			cmd.result <- keys
		}
	}
}


// Create a new shared map with sync type
func NewSMap(t syncType) SharedMap {

	var sm SharedMap
	if (t == RW_MUTEX) {
		sm = sharedMapRWMutex{
			m: make(map[int64]interface{}),
		}
	}

	if (t == CHANNEL) {
		sm = sharedMapChannel{
			m: make(map[int64]interface{}),
		}
	}

	return sm
}
