package nonogram

import (
	"sync"
)

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

	// waitGroup pour attendre la fin de la recursion
	wg := new(sync.WaitGroup)

	// queue d'entree du worker pool
	workQueue := make(chan func())

	// lance les workers
	for i := 0; i < nbWorkers; i++ {
		go func(id int) {
			for tache := range workQueue {
				tache()
				wg.Done()
			}
		}(i)
	}

	// lance la recherche
	solveRecursifConcurrent(&allBlocs, nil, allBlocs.cols, workQueue, wg, solutions)

	// attend la fin du traitment et ferme le channel
	go func() {
		wg.Wait()
		close(solutions)
	}()

	return solutions
}

func solveRecursifConcurrent(allBlocs *allPossibleBlocs,
	tjPartiel TabJeu,
	colonnes lineListSet,
	workQueue chan func(),
	wg *sync.WaitGroup,
	solutions chan *TabJeu) {

	taille := len(allBlocs.rows)
	idxLigneCourante := len(tjPartiel)
	tryLines := (*allBlocs).rows[idxLigneCourante]

	nextTasks := make([]func(), 0, taille)

	// essaye toutes le combinaisons possibles pour la ligne courante
	for n, nextLigne := range tryLines {

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
				workQueue,
				wg,
				solutions)
		}
		nextTasks = append(nextTasks, nextTask)
	}
	wg.Add(len(nextTasks))
	go func() {
		for _, task := range nextTasks {
			workQueue <- task // ajoute la tache avec une coroutine
		}
		// // fmt.Printf("%d tâches ajoutées à la queue\n", len(nextTasks))
	}()
}
