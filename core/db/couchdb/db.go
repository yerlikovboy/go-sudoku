package couchdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-sudoku/core/types"
	"log"
	"net/http"
)

const dbURL = "http://hostname:5984/%s"
const docURL = "http://hostname:5984/%s/%s"

//DB ...
type DB struct {
	Name string
	clnt *couchDBClient
}

//NewDatabase ..
func NewDatabase(dbname string) DB {
	return DB{
		Name: dbname,
		clnt: &couchDBClient{
			clnt: &http.Client{},
			cfg:  defaultConfig(),
		},
	}
}

//DocCount ...
func (db DB) DocCount() uint32 {
	urlStr := fmt.Sprintf(dbURL, db.Name)
	log.Printf("getting doc count for db %s", db.Name)
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := db.clnt.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(res.Body)
	var dbInfo DBInfo
	err = decoder.Decode(&dbInfo)
	if err != nil {
		log.Fatal(err)
	}
	return dbInfo.DocCount
}

//GetDocByID ..
func (db DB) GetDocByID(id string, docPtr interface{}) {

	urlStr := fmt.Sprintf(docURL, db.Name, id)

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := db.clnt.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(docPtr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("doc: %v", docPtr)
}

func (db DB) StoreDoc(b types.Board) {
	urlStr := fmt.Sprintf(dbURL, db.Name)

	payload, err := getPayload(b)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("error creating request: %v", err)
		return
	}

	resp, err := db.clnt.Do(req)
	if err != nil {
		log.Printf("error storing document: %v", err)
		return
	}

	if resp.StatusCode != http.StatusCreated {
		log.Printf("unexpected response from doc upload: %v", resp)
	} else {
		log.Printf("document was stored successfully")
	}

}

func getPayload(b types.Board) ([]byte, error) {

	type puzzle struct {
		ID              string     `json:"_id,omitempty"`
		NumClues        uint8      `json:"n_clues"`
		Cells           types.Grid `json:"grid"`
		SolutionID      string     `json:"solution_id"`
		GeneratedMillis uint64     `json:"generated_millis"`
	}

	p := puzzle{
		NumClues:        b.NumClues(),
		Cells:           b.Cells(),
		SolutionID:      b.DerivedFromID(),
		GeneratedMillis: b.CreatedTS(),
	}

	return json.Marshal(p)

}
