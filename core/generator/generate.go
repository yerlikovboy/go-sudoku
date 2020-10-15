package generator

import (
	"math/rand"
	"time"

	"go-sudoku/core/types"
)

const (
	ByRowOffset    = 1
	ByColumnOffset = 9
)

type GroupFinder func([]uint8) []ClueGroup

// ClueGroup is a data structure that denotes areas in a row or column where
// there are more than one values next to each other.
// For example, in the following rows, there are two contiguous parts:
//
// 	first ClueGroup would start at index 1 and end at index 5 and have length of 3
// _  1 2 4 _ _ 7 8 9
type ClueGroup struct {
	Start  uint8
	Length uint8
	Offset uint8
}

func NewClueGroup(s, l, offset int) ClueGroup {
	return ClueGroup{
		Start:  uint8(s),
		Length: uint8(l),
		Offset: uint8(offset),
	}
}

func IsEqual(l, r *ClueGroup) bool {
	return l.Start == r.Start &&
		l.Length == r.Length &&
		l.Offset == r.Offset
}

func scan_rows(r []uint8) []ClueGroup {
	// -1 means we currently are not on a contiguous set of clues
	s := -1

	var res []ClueGroup
	for i, v := range r {
		if v != 0 {
			if s == -1 {
				// found a clue
				s = i
			}
		}

		if v == 0 && s != -1 {
			res = append(res, NewClueGroup(s, i-s, ByRowOffset))
			s = -1
		}
	}
	// if its at the end
	if s != -1 {
		res = append(res, NewClueGroup(s, len(r)-s, ByRowOffset))
	}
	return res
}

func scan_columns(g []uint8) []ClueGroup {

	var res []ClueGroup

	for c := 0; c < 9; c++ {
		s := -1
		for i := 0; i < 9; i++ {
			// index equals
			idx := i*9 + c
			v := g[idx]
			if v != 0 {
				// found a new clue
				if s == -1 {
					s = idx
				}
			} else if v == 0 && s != -1 {
				// end of clue group
				len := (idx - s) / 9
				res = append(res, NewClueGroup(s, len, ByColumnOffset))
				s = -1
			}
		}
		if s != -1 {
			len := (81 + c - s) / 9
			res = append(res, NewClueGroup(s, len, ByColumnOffset))
		}
	}
	return res
}

func candidates(g types.Grid, finder GroupFinder) []uint8 {

	var r []uint8

	cgs := finder(g[:])

	for _, cg := range cgs {

		multiplier := cg.Length
		if multiplier > 9 {
			multiplier = 9
		}
		var i uint8

		for i = 0; i < multiplier; i++ {
			// need to make a slice of numbers which holds the index
			tmp := make([]uint8, cg.Length)
			var j uint8
			for j = 0; j < cg.Length; j++ {
				tmp[j] = cg.Start + j*cg.Offset
			}
			r = append(r, tmp[:]...)
		}
	}

	return r
}

func Make(s types.Board, num_clues uint8) types.Board {

	rand.Seed(time.Now().UnixNano())
	var lc types.Grid
	cells := s.Cells()
	copy(lc[:], cells[:])

	var i uint8
	num_zeros := 81 - num_clues
	var fnMap = map[uint8]GroupFinder{
		0: scan_rows,
		1: scan_columns,
	}
	for i = 0; i < num_zeros; i++ {

		c := candidates(lc, fnMap[i%2])

		pick_idx := rand.Intn(len(c))
		pick := c[pick_idx]
		lc[pick] = 0
	}

	return types.NewBoard(lc).
		WithDerivedFromID(s.ID).
		WithCreatedTS(uint64(time.
			Now().
			UnixNano()))
}
