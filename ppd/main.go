package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"go-sudoku/core/db/couchdb"
	"go-sudoku/core/types"
)

const iso8601DateFormat = "2006-01-02"
const defaultPuzzleID = "2bcbe1df9fa759dd48624772f20675bb"

func getID(v url.Values) string {
	elems, ok := v["id"]
	if !ok || len(elems) == 0 {
		return defaultPuzzleID
	}
	return elems[0]
}

func getPuzzle(clnt *http.Client) func(http.ResponseWriter, *http.Request) {
	db := couchdb.NewDatabase("puzzles")

	return func(w http.ResponseWriter, req *http.Request) {
		id := getID(req.URL.Query())
		log.Printf("puzzle id: %v", id)

		p := types.Puzzle{}
		db.GetDocByID(id, &p)
		log.Printf("puzzle: %v", p)

		js, err := json.Marshal(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func corsHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wPtr := &w
		(*wPtr).Header().Set("Access-Control-Allow-Origin", "*")
		(*wPtr).Header().Set("Access-Control-Allow-Credentials", "true")
		(*wPtr).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		(*wPtr).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if r.Method == "OPTIONS" {
			// handle pre-flight
			log.Printf("Handling pre-flight (OPTIONS)")
		} else {
			handler(w, r)
		}
	}
}

//type HandleFn func(http.ResponseWriter, *http.Request)

func main() {
	clnt := &http.Client{}
	http.HandleFunc("/potd", corsHandler(getPuzzle(clnt)))
	log.Printf("starting server")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
