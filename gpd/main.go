package main

import (
	"go-sudoku/core/db/couchdb"
	"go-sudoku/core/generator"
	"net/http"

	"flag"
)

func app(isDaemon bool, nClues uint8) {
	db := couchdb.NewDB(&http.Client{})
	for {
		g := db.Solution()
		p := generator.Make(g, nClues)
		db.StorePuzzle(p)

		if !isDaemon {
			break
		}

	}
}

func main() {
	numClues := flag.Uint("n", 38, "number of clues (default 38)")
	isDaemon := flag.Bool("d", false, "run as daemon")
	flag.Parse()
	app(*isDaemon, uint8(*numClues))
}
