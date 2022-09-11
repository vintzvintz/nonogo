package solver

import (
	"fmt"
	"testing"
	"time"

	TJ "vintz.fr/nonogram/tabjeu"
)

const DEFAULT_SIZE = 10
const DEFAULT_SEED = 1006
const DEFAULT_RATIO = 0.45


func compareBlocList(blocs1, blocs2 []TJ.BlocCount) bool {

	if len(blocs1) != len(blocs2) {
		return false
	}
	for b := range blocs1 {
		b1 := blocs1[b]
		b2 := blocs2[b]
		if len(b1) != len(b2) {
			return false
		}
		for i := range b1 {
			if b1[i] != b2[i] {
				return false
			}
		}
	}
	return true
}


func checkSolution(sol TJ.TabJeu, blocsLignes, blocsColonnes TJ.BlocCountList) bool {
	lignesOk := compareBlocList(blocsLignes, sol.CompteBlocs(TJ.LIGNE))
	colsOk := compareBlocList(blocsColonnes, sol.CompteBlocs(TJ.COLONNE))
	return colsOk && lignesOk
}


func TestConcurrent(t *testing.T) {

	size := 15
	nbPlein := int(DEFAULT_RATIO * float32(size))
	tj := TJ.NewTabJeu(size, nbPlein, DEFAULT_SEED)
	fmt.Println(tj)

	bcLigne := tj.CompteBlocs(TJ.LIGNE)
	bcCol := tj.CompteBlocs(TJ.COLONNE)

	nbWorkers := []int{0, 1, 2, 4, 6, 12, 100}
	//nbWorkers := []int{0}

	// retient le nombre de solutions pour chaque nombre de workers
	nbSolutions := make([]int, len(nbWorkers))

	// teste l'algorithme avec différents nombre de workers
	for i := range nbWorkers {

		startTime := time.Now()
		txt := fmt.Sprintf("%d workers", nbWorkers[i])
		solver := makeConcurrentSolver(nbWorkers[i], true)

		// verifie toutes les solutions renvoyées
		var nbExact, nbBad, nbTotal int
		for sol:= range solver(tj) {
			
			//sol.AfficheAvecComptes()
			// est-ce la solution exacte ?
			if tj.Compare(sol,nil)==0 {
				nbExact++
			}
			// est-ce une solution valide différente ?
			if !checkSolution( tj, bcLigne, bcCol) {
				nbBad++
			}
			nbTotal++
		}
		duree := time.Since(startTime)
		t.Logf("%s : %d solutions trouvées en %v\n", txt, nbTotal, duree)
		if nbBad>0 {
			t.Errorf("%d/%d solutions invalides\n", nbBad, nbTotal)
		}
		if nbExact != 1 {
			t.Errorf("Solution exacte non trouvée\n")
		}
		nbSolutions[i] = nbTotal
	}

	// vérifie que le nombre de solutions ne dépend pas du nb de workers
	for i := 0; i < len(nbWorkers)-1; i++ {
		if nbSolutions[i] != nbSolutions[i+1] {
			t.Errorf("Nombre de solutions variable selon le nb de workers")
			break
		}
	}
}

func makeConcurrentSolver(nbWorkers int, showPerf bool) SolverFunc {
	solver := func(prob TJ.TabJeu) chan TJ.TabJeu {
		return SolveConcurrent(prob, nbWorkers, showPerf)
	}
	return solver
}
