package main

import (
	"errors"
)

// ErrNotFound : Not found error
var ErrNotFound = errors.New("Key not found in memory map")

// LogFile : name of the file
// NotFoundID : message of the missing parameter
const (
	// LogFile : File name of log
	LogFile                    string = "development.log"
	NotFoundID                 string = "id is missing in parameters"
	SuccessfullyFlushedMessage string = "Successfully flushed"
	StoreFileName              string = "TIMESTAMP-data.json"
)
