package nonogram

import (
	"sync"
	"time"
)

type WorkQueue struct {
	ch chan func()     // queue d'entree du worker pool
	wg *sync.WaitGroup // attend la fin de la recursion
}

type IdxCols []int
type IdxColsSet []IdxCols

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
	workQueue.wg = new(sync.WaitGroup)
	// recursion monothread classique si workQueue.ch = nil

	if nbWorkers > 0 {
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

	perf := NewPerfCounter(time.Second)

	// prepare la liste initiale des colonnes valides = toutes les combinaisons possibles
	allCols := allBlocs.cols
	colonnes := make(IdxColsSet, prob.taille)
	for numCol := range allCols {
		colsN := make(IdxCols, len(allCols[numCol]))
		for n := range allCols[numCol] {
			colsN[n] = n
		}
		colonnes[numCol] = colsN
	}

	// lance la recherche
	if workQueue.ch == nil {
		workQueue.wg.Add(1)
		go func() {
			// bloque jusqu'à la fin de la recherche, car workqueue.ch == nil
			solveRecursif(&allBlocs, nil, colonnes, workQueue, solutions, perf)
			workQueue.wg.Done()
		}()
	} else {
		// lance la recherche. nonbloquant car workQueue.ch != nil
		solveRecursif(&allBlocs, nil, colonnes, workQueue, solutions, perf)
	}

	// attend la fin du traitment pour fermer le channel
	go func() {
		workQueue.wg.Wait()
		close(solutions)
		perf.Stop()
	}()

	return solutions
}

func solveRecursif(allBlocs *allPossibleBlocs,
	tjPartiel []int, // index (dans allBlocs) des lignes déja placées
	colonnes IdxColsSet, // index dans allBlocs des colonnes encore valides
	wq WorkQueue,
	solutions chan *TabJeu,
	perf *PerfCounter) {

	taille := len(allBlocs.rows)
	numLigneCourante := len(tjPartiel)
	tryLines := (*allBlocs).rows[numLigneCourante]

	nextTasks := make([]func(), 0, len(tryLines))

	// essaye toutes le combinaisons possibles pour la ligne courante
	for n, nextLigne := range tryLines {

		perf.Inc() // pour mesurer la vitesse

		// parmi les colonnes reçues du parent, elimine celles incompatibles avec nextLine
		nextColonnes, ok := filtreColonnes(allBlocs, nextLigne, numLigneCourante, colonnes)

		// abandonne nextLigne s'il n'y a plus de colonnes valides
		if !ok {
			continue
		}

		// copie le tj partiel reçu du parent et ajoute la ligne courante
		tjNext := make([]int, numLigneCourante+1)
		copy(tjNext, tjPartiel)
		tjNext[numLigneCourante] = n // index d'une ligne dans allBlocs.rows

		// on a trouvé une solution si toutes les lignes sont remplies
		if len(tjNext) == taille-1 {
			solutions <- tabJeuFromIndex(allBlocs, tjNext)
			continue
		}

		// prepare la recherche sur les lignes restantes
		nextTask := func() {
			solveRecursif(allBlocs, tjNext, nextColonnes, wq, solutions, perf)
		}

		if wq.ch == nil {
			// execution immediate (bloquante) : recursion classique mono-thread
			nextTask()
		} else {
			// execution différée pour recursion concurrente : ajoute à la liste des tâches à traiter
			nextTasks = append(nextTasks, nextTask)
		}
	}

	// envoi des tâches au pool de workers
	wq.wg.Add(len(nextTasks))
	go func() {
		for _, task := range nextTasks {
			wq.ch <- task
		}
	}()
}

func filtreColonnes(allBlocs *allPossibleBlocs,
	ligne LigneJeu,
	numLigne int,
	colonnes IdxColsSet) (filteredCols IdxColsSet, ok bool) {

	taille := len(colonnes)

	// filteredCols va recevoir les colonnes valides avec la ligne courante
	filteredCols = make(IdxColsSet, taille)
	for numCol := 0; numCol < taille; numCol++ {

		// validCols va contenir toutes possibilités encore valides pour la colonne iCol
		validCols := make(IdxCols, 0, len(colonnes[numCol]))

		cellLignePlein := ligne[numCol].estPlein()

		for _, n := range colonnes[numCol] {
			col := (*allBlocs).cols[numCol][n]
			cellColPlein := col[numLigne].estPlein()
			if cellLignePlein == cellColPlein {
				validCols = append(validCols, n)
			}
		}

		// arrete dès qu'une colonne est incompatible avec la ligne ajoutée
		if len(validCols) == 0 {
			return nil, false
		}
		filteredCols[numCol] = validCols
	}
	return filteredCols, true
}

// tabJeuFromIndex construit un tabJeu à partir des index de lignes
func tabJeuFromIndex(allBlocs *allPossibleBlocs, index []int) *TabJeu {
	taille := len(index)
	tj := make(TabJeu, taille)
	for n, i := range index {
		tj[n] = (*allBlocs).rows[n][i]
	}
	return &tj
}
