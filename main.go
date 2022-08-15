package main

import (
	"fmt"

	"vintz.fr/nonogram/level"
)

func main() {
	real_main()
}

func real_main() {
	lvl := level.NewDefault()
	fmt.Println(lvl)
}


/*


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
*/
