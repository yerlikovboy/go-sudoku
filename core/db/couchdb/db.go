package couchdb

import (
	"encoding/json"
	"fmt"
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
