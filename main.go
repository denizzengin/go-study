package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Go study in-memory rest-api")
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var input KeyValuePair
	json.Unmarshal(reqBody, &input)
	set(input.Key, input.Value)
	json.NewEncoder(w).Encode(input)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, ok := vars["id"]
	if !ok {
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
	fmt.Fprintf(w, "%+v", SuccessfullyFlushedMessage)
}

// See : https://stackoverflow.com/questions/38443889/logging-http-responses-in-addition-to-requests
func logHandler(logfile *os.File, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer logfile.Close()
		x, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		logfile.WriteString(fmt.Sprintf("%q \n", x))
		rec := httptest.NewRecorder()
		fn(rec, r)
		logfile.WriteString(fmt.Sprintf("%q \n", rec.Body))
		rec.Body.WriteTo(w)
	}
}

func handleRequests() {
	logFile := openFile(LogFile)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/set", logHandler(logFile, setHandler)).Methods("POST")
	router.HandleFunc("/get/{id}", logHandler(logFile, getHandler))
	router.HandleFunc("/flush", logHandler(logFile, flushHandler)).Methods("DELETE")
	e := http.ListenAndServe(":8084", router)
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
