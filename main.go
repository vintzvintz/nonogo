package main

import (
	"vintz.fr/gotest/nonogram"
)

func main() {
	tj := nonogram.NewTabJeu(14, 42)
	tj.AfficheAvecComptes()
	prob := tj.MakeProbleme()
	nonogram.Bench(prob)
}
