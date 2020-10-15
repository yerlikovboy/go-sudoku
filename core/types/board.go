package types

//Grid holds the values on the sudoku board
type Grid [81]uint8

//Dim holds either a row or column of values
type Dim [9]uint8

//Board holds values on the board as well as info about when the
// board was created, its solution, etc.
type Board struct {
	ID            string
	derivedFromID string
	createdTS     uint64
	// TODO: Create a singleton which contains these. These values are static and do not
	// mutate in individual instances.
	rowIdx map[uint8]Dim
	colIdx map[uint8]Dim
	cells  Grid
}

//NewBoard returns a new Sudoku board with provided values
func NewBoard(g Grid) Board {

	rowMap := make(map[uint8]Dim, 9)
	colMap := make(map[uint8]Dim, 9)

	var i uint8
	for i = 0; i < 9; i++ {
		rowMap[i] = getRow(i)
		colMap[i] = getColumn(i)
	}

	return Board{
		rowIdx: rowMap,
		colIdx: colMap,
		cells:  g,
	}
}

//WithDerivedFromID returns a copy of the board with the updated derivedFromID field
func (b Board) WithDerivedFromID(id string) Board {
	newObj := b
	newObj.derivedFromID = id
	return newObj
}

func (b Board) WithCreatedTS(ts uint64) Board {
	newObj := b
	newObj.createdTS = ts
	return newObj
}

func (b Board) WithID(id string) Board {
	newObj := b
	newObj.ID = id
	return newObj
}

func (b Board) DerivedFromID() string {
	return b.derivedFromID
}

func (b Board) CreatedTS() uint64 {
	return b.createdTS
}

func (b Board) Cells() Grid {
	return b.cells
}

func (b Board) NumClues() uint8 {
	var c uint8
	c = 0
	for _, v := range b.cells {
		if v != 0 {
			c++
		}
	}
	return c
}

func (b Board) IsSolved() bool {
	for _, v := range b.cells {
		if v == 0 {
			return false
		}
	}
	return true
}

func (b Board) RowIdx(rowNum uint8) []uint8 {
	rv, _ := b.rowIdx[rowNum]
	return rv[:]
}

func (b Board) ColIdx(colNum uint8) []uint8 {
	rv, _ := b.rowIdx[colNum]
	return rv[:]
}

func getRow(rownum uint8) [9]uint8 {
	var s [9]uint8
	start_idx := rownum * 9
	var i uint8
	for i = 0; i < 9; i++ {
		s[i] = start_idx + i
	}
	return s
}

func getColumn(colNum uint8) [9]uint8 {
	var s [9]uint8
	for i := 0; i < 9; i++ {
		s[i] = colNum + uint8(i*9)
	}
	return s
}
