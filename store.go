package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"time"
)

type StoreKeyValuePair struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type StoreKeyValuePairs struct {
	StoreKeyValuePairs []StoreKeyValuePairAll `json:"elements"`
}

type StoreKeyValuePairAll struct {
	Data     []StoreKeyValuePair `json:"Data"`
	LastTime time.Time           `json:"LastTime"`
}

type byLastTime []StoreKeyValuePairAll

func (t byLastTime) Len() int {
	return len(t)
}

func (t byLastTime) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t byLastTime) Less(i, j int) bool {
	return t[i].LastTime.Sub(t[j].LastTime) > 0
}

func WriteStore() {
	file := OpenFile(StoreFileName)
	defer file.Close()
	allRecords := ReadStore(false)	
	data := convertToDbObject()
	allRecords.StoreKeyValuePairs = append(allRecords.StoreKeyValuePairs, data)
	byteArray, err := json.Marshal(allRecords)
	if err != nil {
		panic(err)
	}	
	if _, err := file.Write(byteArray); err != nil {
		panic(err)
	}
}

func ReadStore(initialize bool) *StoreKeyValuePairs {
	if initialize {
		customMap = &InMemoryMap{KeyValuePair: make(map[string]string)}	
	}	
	file := OpenFile(StoreFileName)
	defer file.Close()
	data, err := ioutil.ReadFile(file.Name())
	if err != nil {
		fmt.Println(err)
	}
	var all StoreKeyValuePairs
	json.Unmarshal(data, &all)
	sort.Sort(byLastTime(all.StoreKeyValuePairs))
	return &all
}

func ReadStoreFirst(){
	allRecords := ReadStore(true)
	if len(allRecords.StoreKeyValuePairs) > 0 {
		lastMapObj := allRecords.StoreKeyValuePairs[0]
		for i := 0; i < len(lastMapObj.Data); i++ {
			customMap.KeyValuePair[lastMapObj.Data[i].Key] = lastMapObj.Data[i].Value
		}
	}
}

func convertToDbObject() StoreKeyValuePairAll {
	list := make([]StoreKeyValuePair, len(customMap.KeyValuePair))

	index := 0
	for k, v := range customMap.KeyValuePair {
		list[index] = StoreKeyValuePair{Key: k, Value: v}
		index++
	}

	dataObj := StoreKeyValuePairAll{Data: list, LastTime: time.Now()}
	return dataObj
}
