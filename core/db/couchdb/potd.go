package couchdb

import (
	"net/http"

	"go-sudoku/core/db"
	"go-sudoku/core/types"
)

type CouchPOTD struct {
	clnt *http.Client
	cfg  Config
}

func New(clnt *http.Client) db.POTD {
	return CouchPOTD{
		clnt: clnt,
		cfg:  DefaultConfig(),
	}
}

func (p CouchPOTD) GetPOTD(s string) types.Puzzle {
	return types.Puzzle{}
}
