package main

import (
	"errors"
)

var ErrNotFound = errors.New("Key not found in memory map")

const (
	LogFile                    string = "development.log"
	NotFoundId                 string = "id is missing in parameters"
	SuccessfullyFlushedMessage string = "Successfully flushed"
	StoreFileName              string = "TIMESTAMP-data.json"
)
