package main

import (
	"fmt"
	"runtime"
	"time"

	"vintz.fr/nonogram/solver"
	"vintz.fr/nonogram/tabjeu"
)

func main() {
	essai()
}

func essai() {
	tj := tabjeu.NewTabJeu(5, 42, 1003)
	tj.AfficheAvecComptes()
	prob := tj.MakeProbleme()

	nbCPU := runtime.NumCPU()
	_ = nbCPU

	nbWorkers := 2
	txt := fmt.Sprintf("%d workers", nbWorkers)

	for iter := 0; iter < 1; iter++ {
		startTime := time.Now()
		var nbSol int
		for sol := range solver.SolveConcurrent(prob, nbWorkers, true) {
			_ = sol
			//sol.AfficheAvecComptes()
			//fmt.Println(sol)
			nbSol++
		}
		duree := time.Since(startTime)
		fmt.Printf("%s : %d solutions trouvées en %v\n", txt, nbSol, duree)
	}
}
