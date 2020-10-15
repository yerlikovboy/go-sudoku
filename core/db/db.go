package db

import "go-sudoku/core/types"

type SudokuDB interface {
	Solution() types.Board
	StorePuzzle(s types.Board)
	PickPuzzle() types.Puzzle
}

type POTD interface {
	//	setPOTD(string, types.Puzzle) error
	GetPOTD(string) types.Puzzle
}
