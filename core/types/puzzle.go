package types

type Puzzle struct {
	ID              string `json:"_id,omitempty"`
	NumClues        uint8  `json:"n_clues"`
	Cells           Grid   `json:"grid"`
	SolutionID      string `json:"solution_id"`
	GeneratedMillis uint64 `json:"generated_millis"`
}

func FromBoard(b Board) Puzzle {
	return Puzzle{
		NumClues:        b.NumClues(),
		Cells:           b.Cells(),
		SolutionID:      b.DerivedFromID(),
		GeneratedMillis: b.CreatedTS(),
	}
}
