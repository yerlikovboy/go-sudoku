package main

import (
	"fmt"
	"go-sudoku/core/db/couchdb"
	"go-sudoku/core/types"
	"net/http"
)

func main() {

	db := couchdb.NewDatabase("puzzles", &http.Client{})
	fmt.Printf("Doc Count for db puzzles: %v\n", db.DocCount())
	p := types.Puzzle{}
	db.GetDocByID("0e122d16f058318b3e06a555830032b3", &p)
	fmt.Printf("Sample Doc: %v", p)
}
