package main

import (
	"vintz.fr/nonogram/solver"
	"vintz.fr/nonogram/tabjeu"
)

func main() {
	tj := tabjeu.NewTabJeu(15, 40, 1003)
	tj.AfficheAvecComptes()
	prob := tj.MakeProbleme()
	solver.Bench(prob, false)
}
