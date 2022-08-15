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
	valide := desambigueTabJeu(tj)


	valide.AfficheAvecComptes()
	desambigueTabJeu(valide)
}







func desambigueTabJeu(tj tabjeu.TabJeu) (valide tabjeu.TabJeu) {

	taille := len(tj)

	nbCPU := runtime.NumCPU()
	_ = nbCPU

	nbWorkers := nbCPU
	txt := fmt.Sprintf("%d workers", nbWorkers)

	startTime := time.Now()
	var nbSol int
	var prev tabjeu.TabJeu
	var diff tabjeu.Diff = tabjeu.NewDiff(taille)
	for solution := range solver.SolveConcurrent(tj, nbWorkers, true) {

		if prev != nil {
			solution.Compare(prev, diff)
		}
		prev = solution
		solution.AfficheAvecComptes()
		//fmt.Println(sol)
		nbSol++
	}

	fmt.Printf("%v cellules ambigues\n", diff.Count())

	duree := time.Since(startTime)
	fmt.Printf("%s : %d solutions trouvées en %v\n", txt, nbSol, duree)
	
	valide = tj.Copy()
	valide.ReveleVides(diff)

	return valide
}
