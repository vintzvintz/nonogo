package nonogram

import (
	"time"
)

type IdxCols []int
type IdxColsSet []IdxCols

func makeConcurrentSolver(nbWorkers int, showPerf bool) SolverFunc {
	solver := func(prob Probleme) chan *TabJeu {
		return SolveConcurrent(prob, nbWorkers, showPerf)
	}
	return solver
}

func SolveConcurrent(prob Probleme, nbWorkers int, showPerf bool) chan *TabJeu {
	solutions := make(chan *TabJeu)
	allBlocs := allPossibleBlocs{
		rows: buildAllSequences(prob.taille, prob.seqLignes),
		cols: buildAllSequences(prob.taille, prob.seqColonnes),
	}

	var workerPool *WorkerPool = NewWorkQueue(nbWorkers)

	var perf *PerfCounter
	if showPerf {
		perf = NewPerfCounter(time.Second)
	}

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

	termine := func() {
		close(solutions)
		if perf != nil {
			perf.Stop()
		}
	}

	// lance la recherche
	if workerPool != nil {
		// recursion concurrente
		solveRecursif(&allBlocs, nil, colonnes, workerPool, solutions, perf)
		// attend la fin du traitment pour fermer le channel
		go func() {
			workerPool.Wait()
			termine()
		}()
	} else {
		go func(){
			// recursion bloquante classique
			solveRecursif(&allBlocs, nil, colonnes, nil, solutions, perf)
			termine()
		}()
	}

	return solutions
}

func solveRecursif(allBlocs *allPossibleBlocs,
	tjPartiel []int, // index (dans allBlocs) des lignes déja placées
	colonnes IdxColsSet, // index dans allBlocs des colonnes encore valides
	wp *WorkerPool,
	solutions chan *TabJeu,
	perf *PerfCounter) {

	taille := len(allBlocs.rows)
	numLigneCourante := len(tjPartiel)
	tryLines := (*allBlocs).rows[numLigneCourante]

	// alloue un tableau de fonctions pour poursuivre la recherche
	var nextTasks []func()
	if wp != nil {
		nextTasks = make([]func(), 0, len(tryLines))
	}

	// essaye toutes le combinaisons possibles pour la ligne courante
	for n, nextLigne := range tryLines {

		if perf != nil {
			perf.Inc() // pour mesurer la vitesse
		}

		// parmi les colonnes reçues du parent, elimine celles incompatibles avec nextLine
		nextColonnes, ok := filtreColonnes(allBlocs, nextLigne, numLigneCourante, colonnes)

		// abandonne nextLigne s'il n'y a plus de colonnes valides
		if !ok {
			continue
		}

		// copie le tj partiel reçu du parent et ajoute la ligne courante
		// TODO make/copy en une seule fois
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
			solveRecursif(allBlocs, tjNext, nextColonnes, wp, solutions, perf)
		}

		if wp != nil {
			// execution différée pour recursion concurrente : ajoute à la liste des tâches à traiter
			nextTasks = append(nextTasks, nextTask)
		} else {
			// execution immediate (bloquante) : recursion classique mono-thread
			nextTask()
		}
	}

	// envoi des tâches au pool de workers
	if wp != nil {
		wp.AddTasks(nextTasks)
	}

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
		// TODO allouer en une seule fois
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
