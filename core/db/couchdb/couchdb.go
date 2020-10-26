package couchdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"go-sudoku/core/db"
	"go-sudoku/core/types"
)

//CouchSudokuDB ...
type CouchSudokuDB struct {
	clnt *http.Client
	cfg  config
}

//NewDB ...
func NewDB(clnt *http.Client) db.SudokuDB {
	return CouchSudokuDB{
		clnt: &http.Client{},
		cfg:  defaultConfig(),
	}
}

func (s CouchSudokuDB) puzzleCount() uint32 {

	req, _ := http.NewRequest("GET", "http://hostname:5984/grids/_design/puzzles/_view/completed", nil)
	s.cfg.SetupRequest(req)

	// the way query is set is so lame ...
	q := req.URL.Query()
	q.Add("limit", "1")
	req.URL.RawQuery = q.Encode()

	resp, err := s.clnt.Do(req)
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

func (s CouchSudokuDB) nthGrid(n uint32) grid {

	// log.Printf("pick #%v from view", n)

	req, _ := http.NewRequest("GET", "http://localhost:5984/grids/_design/puzzles/_view/completed", nil)
	s.cfg.SetupRequest(req)

	// the way query is set is bullshit!
	qry := req.URL.Query()
	qry.Add("limit", "1")
	qry.Add("skip", fmt.Sprint(n))
	req.URL.RawQuery = qry.Encode()

	// log.Printf("request: %v", req)
	resp, err := s.clnt.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatal("houston, we have a problem: ", resp.StatusCode)

	}

	decoder := json.NewDecoder(resp.Body)
	var r response
	err = decoder.Decode(&r)
	if err != nil {
		log.Fatal(err)
	}
	// log.Printf("response: %v", r)
	if len(r.Rows) == 0 {
		log.Fatal("unable to retrieve puzzle (rows len == 0)")
	}
	return r.Rows[0]
}

//Solution implementation for CouchDB
func (s CouchSudokuDB) Solution() types.Board {
	rowCount := s.puzzleCount()
	pick := uint32(rand.Int31n(int32(rowCount)))
	grid := s.nthGrid(pick)
	var c types.Grid
	copy(c[:], grid.Value[0:81])

	return types.NewBoard(c).
		WithCreatedTS(grid.Timestamp).
		WithID(grid.ID)
}

func (s CouchSudokuDB) getPuzzleCount() uint32 {
	req, _ := http.NewRequest("GET", "http://localhost:5984/puzzles", nil)
	s.cfg.SetupRequest(req)

	resp, err := s.clnt.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("houston, we have a problem: ", resp.StatusCode)

	}

	decoder := json.NewDecoder(resp.Body)
	var r DBInfo
	err = decoder.Decode(&r)
	if err != nil {
		log.Fatal(err)
	}
	return r.DocCount
}

//PickPuzzle implementation for CouchDB
func (s CouchSudokuDB) PickPuzzle() types.Puzzle {
	pick := rand.Int31n(int32(s.getPuzzleCount()))

	reqBody, reqErr := json.Marshal(NewPuzzleRequest(uint32(pick)))
	if reqErr != nil {
		log.Fatal(reqErr)
	}
	log.Printf("pick puzzle req: %s", reqBody)
	req, _ := http.NewRequest("POST", "http://localhost:5984/puzzles/_all_docs", bytes.NewBuffer(reqBody))
	s.cfg.SetupRequest(req)
	qry := req.URL.Query()
	qry.Add("limit", "1")
	qry.Add("skip", fmt.Sprint(pick))
	qry.Add("include_docs", "true")
	req.URL.RawQuery = qry.Encode()

	log.Printf("request: %v", req)
	puzzleRes, resErr := s.clnt.Do(req)
	if resErr != nil {
		log.Fatal(resErr)
	}
	decoder := json.NewDecoder(puzzleRes.Body)
	var r puzzle_request_result
	err := decoder.Decode(&r)
	if err != nil {
		log.Fatal(err)
	}
	if len(r.Rows) != 1 {
		log.Fatal("unexpected result from db: ", r)
	}
	p := r.Rows[0]
	return p.Doc
}

//StorePuzzle implementation for CouchDB
func (s CouchSudokuDB) StorePuzzle(b types.Board) {
	// for now, just print it to console
	p := types.FromBoard(b)
	raw, _ := json.Marshal(p)
	fmt.Println(string(raw))

	req, _ := http.NewRequest("POST", "http://localhost:5984/puzzles", bytes.NewBuffer(raw))
	s.cfg.SetupRequest(req)

	resp, err := s.clnt.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 201 {
		log.Printf("puzzle upload status code: %v", resp)
	}

}
