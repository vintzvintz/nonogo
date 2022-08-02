package main

import (
	"vintz.fr/gotest/nonogram"
)

func main() {
	tj := nonogram.NewTabJeu(15, 45, 1003)
	tj.AfficheAvecComptes()
	prob := tj.MakeProbleme()
	nonogram.Bench(prob, false)
}
