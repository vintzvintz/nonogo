package solver

import (
	"fmt"
	"testing"
	"time"

	TJ "vintz.fr/nonogram/tabjeu"
)

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

func checkSolution(prob TJ.Probleme, sol *TJ.TabJeu) bool {
	lignesOk := compareBlocList(prob.BlocsLignes, sol.CompteBlocs(TJ.LIGNE))
	colsOk := compareBlocList(prob.BlocsColonnes, sol.CompteBlocs(TJ.COLONNE))
	return colsOk && lignesOk
}

func TestConcurrent(t *testing.T) {

	//tj := TJ.NewTabJeu(15, 45, 1003)
	tj := TJ.NewTabJeu(15, 45, 1003)
	tj.AfficheAvecComptes()
	prob := tj.MakeProbleme()

	//nbWorkers := []int{12, 0, 1, 2, 4, 6, 12}

	nbWorkers := []int{1, 2}

	// retient le nombre de solutions pour chaque nombre de workers
	nbSolutions := make([]int, len(nbWorkers))

	// teste l'algorithme avec différents nombre de workers
	for i := range nbWorkers {

		startTime := time.Now()
		txt := fmt.Sprintf("%d workers", nbWorkers[i])
		solver := makeConcurrentSolver(nbWorkers[i], true)

		// verifie toutes les solutions renvoyées
		var nbBad, nbTotal int
		for sol := range solver(prob) {

			//sol.AfficheAvecComptes()
			if !checkSolution(prob, sol) {
				nbBad++
			}
			nbTotal++
		}
		duree := time.Since(startTime)
		t.Logf("%s : %d solutions trouvées en %v\n", txt, nbTotal, duree)
		if nbBad > 0 {
			t.Errorf("%d/%d solutions erronées", nbBad, nbTotal)
		}
		nbSolutions[i] = nbTotal
	}

	// vérifie que le nombre de solutions ne dépend pas du nb de workers
	for i := 0; i < len(nbWorkers)-1; i++ {
		if nbSolutions[i] != nbSolutions[i+1] {
			t.Errorf("Nombre de solutions différent selon le nb de workers")
			break
		}
	}
}
