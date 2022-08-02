package nonogram

import (
	"fmt"
	"time"
)

// lineList est une liste de lignes
// contrairement à TabJeu la largeur n'est pas forcément égale à la hauteur
type lineList []LigneJeu

// allPossibleLines est un slice de la même longueur que TabJeu
// chaque élement est l'ensemble des valeurs possibles pour la ligne correspondante
type lineListSet []lineList

// structure intermédiaire pour construire l'ensemble des lignes possibles
type indexedLineSet struct {
	num   int
	lines lineList
}

// cellVide et cellPlein sont les éléments de base pour construire les sequences possibles
// on aurait aussi bien pu utiliser des booleéns à la place des pointeurs
var cellVide = &cellule{base: vide, joué: blanc}
var cellPlein = &cellule{base: plein, joué: blanc}

type allPossibleBlocs struct {
	rows lineListSet
	cols lineListSet
}

type SolverFunc func(Probleme) chan *TabJeu

func Bench(prob Probleme, showPerf bool) {

	//BenchSolver(SolveBourrin, prob, "Bourrin")

	tests := []int{0, 1,2,3,4,5,6,12,6,12,6,12,6,12,6,12,24,24,96,96}
	for _, nb := range tests {
		txt := fmt.Sprintf("%d workers", nb)
		solver := makeConcurrentSolver(nb, showPerf)
		BenchSolver(solver, prob, txt)
	}
}

func BenchSolver(solver SolverFunc, prob Probleme, txt string) {
	solutions := solver(prob)
	var nb int
	startTime := time.Now()
	for sol := range solutions {
		nb++
		//fmt.Printf("Solution n°%d\n", nb)
		//fmt.Print(*sol)
		_ = sol
	}
	duree := time.Since(startTime)
	fmt.Printf("%s : %d solutions trouvées en %v\n", txt, nb, duree)
}
