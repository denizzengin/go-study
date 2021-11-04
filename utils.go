package main

import (
	"fmt"
	"os"
	"sync"
)

// KeyValuePair : response model
type KeyValuePair struct {
	Key   string      `json:"Key"`
	Value interface{} `json:"Value"`
}

// InMemoryMap : data struct
type InMemoryMap struct {
	KeyValuePair map[string]string
	Mutex        *sync.Mutex
}

var customMap *InMemoryMap

func set(key string, val interface{}) {
	customMap.Mutex.Lock()
	customMap.KeyValuePair[key] = fmt.Sprint(val)
	customMap.Mutex.Unlock()
}

func get(key string) (string, error) {
	customMap.Mutex.Lock()
	defer customMap.Mutex.Unlock()
	if val, ok := customMap.KeyValuePair[key]; ok {
		return val, nil
	}
	return "", ErrNotFound
}

func flush() {
	customMap = &InMemoryMap{KeyValuePair: make(map[string]string), Mutex: &sync.Mutex{}}
}

func openFile(fileName string) *os.File {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	return file
}
