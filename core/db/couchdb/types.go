package couchdb

import "go-sudoku/core/types"

type response struct {
	TotalRows uint32 `json:"total_rows,omitempty"`
	Rows      []grid `json:"rows"`
}

type grid struct {
	ID        string  `json:"id"`
	Timestamp uint64  `json:"key"`
	Value     []uint8 `json:"value"`
}

type DBInfo struct {
	Name     string `json:"db_name"`
	DocCount uint32 `json:"doc_count"`
}

type puzzle_selector struct {
	NClues int `json:"n_clues"`
}

type puzzle_request struct {
	Selector puzzle_selector `json:"selector"`
	Limit    uint8           `json:"limit"`
	Skip     uint32          `json:"skip"`
}

type puzzle_request_result struct {
	TotalRows uint32 `json:"total_rows,omitempty"`
	Rows      []struct {
		Doc types.Puzzle `json:"doc"`
	} `json:"rows"`
}

func NewPuzzleRequest(nl uint32) puzzle_request {
	return puzzle_request{
		Selector: puzzle_selector{
			NClues: 38,
		},
		Limit: 1,
		Skip:  nl,
	}

}
