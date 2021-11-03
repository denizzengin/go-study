package main

import (
	"fmt"
	"testing"
)

func TestSet(t *testing.T) {
	var val interface{} = 5
	var key string = "x"
	set(key, val)

	if newVal, ok := customMap.KeyValuePair[key]; !ok {
		t.Errorf("Error cannot set %v = %v", newVal, val)
	} else if ok {
		if newVal != fmt.Sprint(val) {
			t.Errorf("Error not equal %v = %v", newVal, val)
		}
	}
}

func TestGet(t *testing.T) {
	var key string = "x"
	var curValue interface{} = 5
	set(key, curValue)
	val, err := get(key)
	if err != nil {
		t.Errorf("Error %v", err)
	}
	if val != fmt.Sprint(curValue) {
		t.Errorf("Error %v", err)
	}
}

func TestFlush(t *testing.T) {
	flush()
	if count := len(customMap.KeyValuePair); count != 0 {
		t.Errorf("Error map isn't empty, count : %v", count)
	}
}

func TestStoreReadFirst(t *testing.T) {
	readStoreFirst()
	if len(customMap.KeyValuePair) <= 0 {
		t.Errorf("Error len must be bigger than 0")
	}
}

func TestStoreWrite(t *testing.T) {
	customMap = &InMemoryMap{KeyValuePair: make(map[string]string)}
	writeStore()
}
