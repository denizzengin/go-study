package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func setHandler(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var input KeyValuePair
	json.Unmarshal(reqBody, &input)
	set(input.Key, input.Value)
	w.WriteHeader(http.StatusCreated)
	s, _ := json.Marshal(input)
	w.Write(s)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/get/")
	if key == "" {
		fmt.Fprintf(w, "%+v", NotFoundID)
	}
	pair, err := get(key)
	if err != nil {
		fmt.Fprintf(w, "%+v", string(err.Error()))
	} else {
		json.NewEncoder(w).Encode(KeyValuePair{Key: key, Value: pair})
	}
}

func flushHandler(w http.ResponseWriter, r *http.Request) {
	flush()
	w.WriteHeader(http.StatusNoContent)
}

func logger(w http.ResponseWriter, r *http.Request) {
	// Log request.
	c, e := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(c))
	if e != nil {
		panic(e)
	}
	log.Printf("%s\t\t%s\t\t%s", r.Method, r.RequestURI, c)

	switch {
	case r.Method == http.MethodGet && strings.TrimPrefix(r.URL.Path, "/get/") != "":
		getHandler(w, r)
	case r.Method == http.MethodPost && strings.TrimPrefix(r.URL.Path, "/set") == "":
		setHandler(w, r)
	case r.Method == http.MethodDelete && strings.TrimPrefix(r.URL.Path, "/flush") == "":
		flushHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}

func handleRequests() {
	http.HandleFunc("/", logger)
	e := http.ListenAndServe(":8084", nil)
	log.Fatal(e)
}

func writeAsync(quit chan bool) {
	// Store every two minutes
	ticker := time.NewTicker(2 * time.Minute)
	for {
		select {
		case <-quit:
			ticker.Stop()
			return
		case <-ticker.C:
			writeStore()
		}
	}
}

func main() {
	fmt.Println("Waiting request...")
	readStoreFirst() // Moved here from init because of wrong use.
	quit := make(chan bool)
	defer close(quit)
	go writeAsync(quit)
	handleRequests()
}
