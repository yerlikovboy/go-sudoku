package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"go-sudoku/core/db/couchdb"
)

const ISO8601_DATE_FORMAT = "2006-01-02"

func getKey(v url.Values) (string, error) {
	elems, ok := v["date"]
	if !ok || len(elems) == 0 {
		return time.Now().Format(ISO8601_DATE_FORMAT), nil
	}
	val := elems[0]
	if _, err := time.Parse(ISO8601_DATE_FORMAT, val); err != nil {
		return "", err
	}
	return val, nil
}

func GetPOTD(clnt *http.Client) func(http.ResponseWriter, *http.Request) {
	db := couchdb.NewDB(&http.Client{})

	return func(w http.ResponseWriter, req *http.Request) {
		dt, err := getKey(req.URL.Query())
		if err != nil {
			log.Printf("error retrieving date: ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Printf("date: %v", dt)

		p := db.PickPuzzle()
		log.Printf("puzzle pick: %v", p)

		js, err := json.Marshal(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

// GetHandler sets up db access object and retuns handler. This way handler
// does not have to create new db object each time it is called
func GetHandler() func(http.ResponseWriter, *http.Request) {
	db := couchdb.NewDB(&http.Client{})

	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf("req: %v", req)
		p := db.PickPuzzle()
		log.Printf("puzzle pick: %v", p)
		js, err := json.Marshal(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func corsHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w_ptr := &w
		(*w_ptr).Header().Set("Access-Control-Allow-Origin", "*")
		(*w_ptr).Header().Set("Access-Control-Allow-Credentials", "true")
		(*w_ptr).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		(*w_ptr).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if r.Method == "OPTIONS" {
			// handle pre-flight
			log.Printf("Handling pre-flight (OPTIONS)")
		} else {
			handler(w, r)
		}
	}
}

type HandleFn func(http.ResponseWriter, *http.Request)

func main() {
	clnt := &http.Client{}
	http.HandleFunc("/puzzle", GetHandler())
	http.HandleFunc("/potd", corsHandler(GetPOTD(clnt)))
	log.Printf("starting server")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
