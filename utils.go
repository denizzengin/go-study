package main

import (
	"fmt"
	"os"
	"sync"
)

type KeyValuePair struct {
	Key   string      `json:"Key"`
	Value interface{} `json:"Value"`
}

type InMemoryMap struct {
	KeyValuePair map[string]string
}

var mutex = &sync.Mutex{}
var customMap *InMemoryMap

func Set(key string, val interface{}) {
	mutex.Lock()
	customMap.KeyValuePair[key] = fmt.Sprint(val)
	mutex.Unlock()
}

func Get(key string) (string, error) {
	mutex.Lock()
	defer mutex.Unlock()
	if val, ok := customMap.KeyValuePair[key]; ok {
		return val, nil
	}
	return "", ErrNotFound
}

func Flush() {
	mutex.Lock()
	customMap = &InMemoryMap{KeyValuePair: make(map[string]string)}
	mutex.Unlock()
}

func OpenFile(fileName string) *os.File {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	return file
}

func init() {
	ReadStoreFirst()
}
