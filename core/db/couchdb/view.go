package couchdb

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const urlTemplate = "http://hostname:5984/%s/_design/%s/_view/%s"

// View is a CoucDB database view
type View struct {
	dbName    string
	designDoc string
	viewName  string

	url  string
	clnt *http.Client
	cfg  config
}

//NewView constructor. See CouchDB docs for more info on database name,
// design docs and views.
func NewView(dbName, designDoc, viewName string, client *http.Client) View {
	return View{
		dbName:    dbName,
		designDoc: designDoc,
		viewName:  viewName,
		url:       fmt.Sprintf(urlTemplate, dbName, designDoc, viewName),
		clnt:      client,
		cfg:       defaultConfig(),
	}
}

//DocCount ...
func (v View) DocCount() uint32 {
	req, err := http.NewRequest("GET", v.url, nil)
	if err != nil {
		log.Fatal(err)
	}

	v.cfg.SetupRequest(req)

	q := req.URL.Query()
	q.Add("limit", "1")
	req.URL.RawQuery = q.Encode()

	resp, err := v.clnt.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatal("houston, we have a problem: ", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	var val response
	err = decoder.Decode(&val)
	if err != nil {
		log.Fatal(err)
	}

	return val.TotalRows
}

// GetDoc still needs work ... specifically a nice way to return the result
func (v View) GetDoc(limit, skip uint32, val interface{}) {
	req, err := http.NewRequest("GET", v.url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// not the greatest way to set up the query
	qry := req.URL.Query()
	qry.Add("limit", fmt.Sprint(limit))
	qry.Add("skip", fmt.Sprint(skip))
	req.URL.RawQuery = qry.Encode()

	// log.Printf("request: %v", req)
	resp, err := v.clnt.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatal("houston, we have a problem: ", resp.StatusCode)

	}

	decoder := json.NewDecoder(resp.Body)
	var r response
	err = decoder.Decode(val)
	if err != nil {
		log.Fatal(err)
	}
	// log.Printf("response: %v", r)
	if len(r.Rows) == 0 {
		log.Fatal("unable to retrieve puzzle (rows len == 0)")
	}
}
