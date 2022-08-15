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
	tj := tabjeu.NewTabJeu(tabjeu.DEFAULT_SIZE, tabjeu.DEFAULT_RATIO, tabjeu.DEFAULT_SEED)
	tj.AfficheAvecComptes()
	valideProbleme(tj)
}

func valideProbleme(tj tabjeu.TabJeu) {

	taille := len(tj)

	nbCPU := runtime.NumCPU()
	_ = nbCPU

	nbWorkers := nbCPU
	txt := fmt.Sprintf("%d workers", nbWorkers)

	for iter := 0; iter < 1; iter++ {
		startTime := time.Now()
		var nbSol int
		var prev tabjeu.TabJeu
		var diff tabjeu.Diff = tabjeu.NewDiff(taille)
		for solution := range solver.SolveConcurrent(tj, nbWorkers, true) {

			if prev != nil {
				solution.Compare(prev, diff)
			}
			prev = solution
			//sol.AfficheAvecComptes()
			//fmt.Println(sol)
			nbSol++
		}

		fmt.Printf("%v cellules ambigues\n", diff.Count())

		duree := time.Since(startTime)
		fmt.Printf("%s : %d solutions trouvées en %v\n", txt, nbSol, duree)
	}

}
