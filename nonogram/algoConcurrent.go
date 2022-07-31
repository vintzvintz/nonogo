package nonogram

import (
	"sync"
	"time"
)

type WorkQueue struct {
	ch chan func()     // queue d'entree du worker pool
	wg *sync.WaitGroup // attend la fin de la recursion
}

func makeConcurrentSolver(nbWorkers int) SolverFunc {
	solver := func(prob Probleme) chan *TabJeu {
		return SolveConcurrent(prob, nbWorkers)
	}
	return solver
}

func SolveConcurrent(prob Probleme, nbWorkers int) chan *TabJeu {
	solutions := make(chan *TabJeu)
	allBlocs := allPossibleBlocs{
		rows: buildAllSequences(prob.taille, prob.seqLignes),
		cols: buildAllSequences(prob.taille, prob.seqColonnes),
	}

	var workQueue WorkQueue

	if nbWorkers > 0 {
		workQueue.wg = new(sync.WaitGroup)
		workQueue.ch = make(chan func())

		// lance les workers
		for i := 0; i < nbWorkers; i++ {
			go func(id int) {
				for tache := range workQueue.ch {
					tache()
					workQueue.wg.Done()
				}
			}(i)
		}
	}

	perf := NewPerfCounter(time.Second / 10)

	termine := func() {
		close(solutions)
		perf.Stop()
	}

	// lance la recherche
	if nbWorkers == 0 {
		go func() {
			// bloquante jusqu'à la fin de la recherche (car workqueue est vide)
			solveRecursifConcurrent(&allBlocs, nil, allBlocs.cols, workQueue, solutions, perf)
			termine()
		}()
	} else {
		// retourne immediatement
		solveRecursifConcurrent(&allBlocs, nil, allBlocs.cols, workQueue, solutions, perf)

		// attend la fin du traitment pour fermer le channel
		go func() {
			workQueue.wg.Wait()
			termine()
		}()
	}

	return solutions
}

func solveRecursifConcurrent(allBlocs *allPossibleBlocs,
	tjPartiel TabJeu,
	colonnes lineListSet,
	wq WorkQueue,
	solutions chan *TabJeu,
	perf *PerfCounter) {

	taille := len(allBlocs.rows)
	idxLigneCourante := len(tjPartiel)
	tryLines := (*allBlocs).rows[idxLigneCourante]

	nextTasks := make([]func(), 0, taille)

	// essaye toutes le combinaisons possibles pour la ligne courante
	for n, nextLigne := range tryLines {

		// pour mesurer la vitesse
		perf.Inc()

		// cree une copie du tj partiel reçu du parent, et ajoute la ligne courante
		tjNext := make(TabJeu, idxLigneCourante+1)
		copy(tjNext, tjPartiel)
		tjNext[idxLigneCourante] = nextLigne

		// copie seulement les colonnes compatibles avec la ligne ajoutée
		nextColonnes, ok := filtreColonnes(tjNext, colonnes)

		_ = n
		//fmt.Printf("ligne %d, combinaison %d/%d, resultat %v\n", idxLigneCourante, n+1, len(tryLines), ok)

		// on arrete la récursion si nextLigne est incompatible avec les colonnes possibles
		if !ok {
			continue
		}
		// on a trouvé une solution si toutes les lignes sont remplies
		if len(tjNext) == taille {
			solutions <- (*TabJeu)(&tjNext)
			continue
		}
		// sinon on continue sur la ligne suivante
		nextTask := func() {
			solveRecursifConcurrent(
				allBlocs,
				tjNext,
				nextColonnes,
				wq,
				solutions,
				perf)
		}
		nextTasks = append(nextTasks, nextTask)
	}

	// recherche concurrente
	if wq.wg != nil {
		// envoi non bloquant des taches aux workers
		wq.wg.Add(len(nextTasks))
		go func() {
			for _, task := range nextTasks {
				wq.ch <- task
			}
		}()
		return
	}

	// recherche recursive classique (depth-first)
	for _, task := range nextTasks {
		task() // traite immédiatement / bloquant
	}
}

func filtreColonnes(tjPartiel TabJeu, colonnes lineListSet) (filteredCols lineListSet, ok bool) {

	taille := len(colonnes)

	/// dernière ligne placée ( = ligne courante de cette fonction )
	idxLine := len(tjPartiel) - 1
	lastLine := tjPartiel[idxLine]

	// filteredCols va recevoir les colonnes valides avec la ligne courante
	filteredCols = make(lineListSet, taille)
	for iCol := 0; iCol < taille; iCol++ {

		// toutes les possibilités valides pour la colonne iCol
		validCols := make(lineList, 0, len(colonnes[iCol]))

		cell_ligne := lastLine[iCol].estPlein()

		for _, col := range colonnes[iCol] {
			cell_col := col[idxLine].estPlein()
			colOk := (cell_ligne == cell_col)
			if colOk {
				validCols = append(validCols, col)
			}
		}

		// inutile de continuer dès qu'une colonne est incompatible avec la ligne ajoutée
		if len(validCols) == 0 {
			return nil, false
		}
		filteredCols[iCol] = validCols
	}
	return filteredCols, true
}
