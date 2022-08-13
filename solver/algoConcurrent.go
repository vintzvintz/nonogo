package solver

import (
	//"fmt"

	"os"
	"sync"
	"time"

	perf "vintz.fr/nonogram/perf"
	TJ "vintz.fr/nonogram/tabjeu"
)

type IdxCols []int
type IdxColsSet []IdxCols

func makeConcurrentSolver(nbWorkers int, showPerf bool) SolverFunc {
	solver := func(prob TJ.Probleme) chan *TJ.TabJeu {
		return SolveConcurrent(prob, nbWorkers, showPerf)
	}
	return solver
}

func SolveConcurrent(prob TJ.Probleme, nbWorkers int, showPerf bool) chan *TJ.TabJeu {
	solutions := make(chan *TJ.TabJeu)
	allBlocs := allPossibleBlocs{
		rows: buildAllSequences(prob.Taille, prob.BlocsLignes),
		cols: buildAllSequences(prob.Taille, prob.BlocsColonnes),
	}

	var workerPool *WorkerPool = NewWorkerPool(nbWorkers)

	var pc *perf.PerfCounter
	if showPerf {
		pc = perf.NewPerfCounter(time.Second, os.Stdout)
	}

	// prepare la liste initiale des colonnes valides = toutes les combinaisons possibles
	allCols := allBlocs.cols
	colonnes := make(IdxColsSet, prob.Taille)
	for numCol := range allCols {
		colsN := make(IdxCols, len(allCols[numCol]))
		for n := range allCols[numCol] {
			colsN[n] = n
		}
		colonnes[numCol] = colsN
	}

	termine := func() {
		if pc != nil {
			pc.Stop()
		}
		close(solutions)
	}

	// lance la récursion
	if nbWorkers > 0 {

		workerPool.Exec(func() {
			solveRecursif(&allBlocs, TjPartiel{}, 0, colonnes, nil, workerPool, solutions, pc)
		})
		// goroutine nécessaire pour attendre la fin du traitment
		go func() {
			workerPool.Wait()
			termine()
		}()
	} else {
		go func() {
			solveRecursif(&allBlocs, TjPartiel{}, 0, colonnes, nil, workerPool, solutions, pc)
			termine() // pas la peine d'attendre dans une goroutine distincte du lancement
		}()
	}

	return solutions
}

type TjPartiel [TJ.TAILLE_MAX]int

func solveRecursif(allBlocs *allPossibleBlocs,
	tjPartiel TjPartiel, // index dans allBlocs des lignes déja placées (longueur variable entre 0 et [taille] éléments)
	numLigneCourante int,
	initialCols IdxColsSet, // index dans allBlocs des colonnes encore valides
	lockCopy *sync.Mutex,
	wp *WorkerPool,
	solutions chan *TJ.TabJeu,
	pc *perf.PerfCounter) {


	// copie l'etat si demandé par l'appelant
	if lockCopy != nil {
		initialCols = initialCols.AllocNew(true)
		lockCopy.Unlock()   // débloque l'execution du parent
	}

	taille := len(allBlocs.rows)
	tryLines := allBlocs.rows[numLigneCourante]

	// met à jour le compteur de vitesse
	if pc != nil {
		pc.Inc(1)
	}

	// alloue un espace pour recevoir les colonnes valides restantes après chaque ligne
	nextCols := initialCols.AllocNew(false)

	lockNext := new(sync.Mutex) // protege l'acces aux parametre de recursion
	var waitChildCopy bool = true

	// essaye toutes les combinaisons encore valides pour la ligne courante
	for n, nextLigne := range tryLines {

		// attend que la récursion lancée au tour précédent ait fini de copier ses paramètres
		if waitChildCopy {
			lockNext.Lock()
			waitChildCopy = false
		}

		// parmi les colonnes reçues du parent, elimine celles incompatibles avec nextLine
		ok := filtreColonnesInplace(allBlocs, nextLigne, numLigneCourante, initialCols, nextCols)

		// abandonne définitivement nextLigne s'il n'y a plus assez de colonnes compatibles
		if !ok {
			continue
		}

		// si la ligne est valide, on l'inscrit dans tjPartiel pour continuer la recherche
		tjPartiel[numLigneCourante] = n

		// condition d'arrêt : on a trouvé une solution si :
		// toutes les lignes sont remplies (et on n'a pas abandonné la ligne en cours)
		if numLigneCourante == taille-1 {
			solutions <- tabJeuFromIndexArray(allBlocs, tjPartiel, taille)
			continue
		}
		// Prepare appel récursif bloquant sans copie des paramètres ( lockCopy=nil )
		// -> pour execution par la même goroutine (recursion classique )
		recurseNoCopy := func() {
			solveRecursif(allBlocs, tjPartiel, numLigneCourante+1, nextCols, nil, wp, solutions, pc )
		}

		// Prepare appel recursif avec copie  des paramètres (lockCopy != nil)
		// -> pour execution par un autre worker dans une autre goroutine
		recurseWithCopy := func() {
			solveRecursif(allBlocs, tjPartiel, numLigneCourante+1, nextCols, lockNext, wp, solutions, pc)
		}

		// essaie de continuer avec un worker, sinon recursion classique sans copie
		accepted := wp.TryExec(recurseWithCopy)
		if accepted {
			// on met un flag pour attendre que l'appelé ait fini de copier ses paramètres au début du prochain tour
			// lockNext est libere par l'appelé dès qu'il a fini
			waitChildCopy = true
			continue
		}
		// si le workerpool a refusé la tâche on continue avec une récursion classique
		recurseNoCopy()
	}
}

// filtreColonnesAlloc alloue et renvoie une nouvelle structure
// filteredCols rçoit les colonnes (désignées par leurs index dans allBlocs), compatibles avec la ligne indiquée
func filtreColonnesAlloc(allBlocs *allPossibleBlocs,
	ligne TJ.LigneJeu,
	numLigne int,
	colonnes IdxColsSet) (filteredCols IdxColsSet, ok bool) {

	taille := len(colonnes)

	// filteredCols va recevoir les colonnes valides avec la ligne courante
	// filteredCols = make(IdxColsSet, taille)
	filteredCols = allocfilteredCols(taille)

	for numCol := 0; numCol < taille; numCol++ {

		// validCols va contenir toutes possibilités encore valides pour la colonne iCol
		// TODO allouer en une seule fois
		//validCols := make(IdxCols, 0, len(colonnes[numCol]))
		validCols := allocValidCols(len(colonnes[numCol]))

		cellLigne := ligne[numCol]

		for _, n := range colonnes[numCol] {
			col := (*allBlocs).cols[numCol][n]
			cellCol := col[numLigne]
			if cellLigne == cellCol {
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

// filtreColonnesInplace est une versison sans allocation de filtreColonnes
// filteredCols rçoit les colonnes (désignées par leurs index dans allBlocs), compatibles avec la ligne indiquée
func filtreColonnesInplace(allBlocs *allPossibleBlocs,
	ligne TJ.LigneJeu,
	numLigne int,
	initialCols, filteredCols IdxColsSet) (ok bool) {

	var cellCol, cellLigne TJ.Cellule

	for numCol := range initialCols {

		// reset de la destination
		filteredCols[numCol] = filteredCols[numCol][:0]

		cellLigne = ligne[numCol] // lookup en dehors de la boucle

		// vérifie si chaque element de initialCols est compatible avec la ligne fournie
		for _, idxCol := range initialCols[numCol] {
			cellCol = allBlocs.cols[numCol][idxCol][numLigne]

			// si out, ajoute l'index de la colonne dans le slice de destination
			if cellLigne == cellCol {
				filteredCols[numCol] = append(filteredCols[numCol], idxCol)
			}
		}

		// arrete dès qu'une colonne est incompatible avec la ligne ajoutée
		if len(filteredCols[numCol]) == 0 {
			return false
		}
	}
	return true
}

// tabJeuFromIndex construit un tabJeu à partir des index de lignes
func tabJeuFromIndexSlice(allBlocs *allPossibleBlocs, index []int) *TJ.TabJeu {
	taille := len(index)
	tj := make(TJ.TabJeu, taille)
	for n, i := range index {
		tj[n] = allBlocs.rows[n][i]
	}
	return &tj
}

// tabJeuFromIndex construit un tabJeu à partir d'un array d'index de lignes
func tabJeuFromIndexArray(allBlocs *allPossibleBlocs, index TjPartiel, taille int) *TJ.TabJeu {
	//taille := len(index)
	tj := make(TJ.TabJeu, taille)
	for n := 0; n < taille; n++ {
		rows := allBlocs.rows[n]
		idx := index[n]
		tj[n] = rows[idx]
	}
	return &tj
}

// AllocNew() alloue des slices de la même taille avec copie optionnelle du contenu
func (src IdxColsSet) AllocNew(withCopy bool) (dst IdxColsSet) {
	dst = allocfilteredCols(len(src))
	for numCol := range src {
		if withCopy {
			dst[numCol] = make(IdxCols, len(src[numCol])) // allocate same length as src
			copy(dst[numCol], src[numCol])
		} else {
			dst[numCol] = make(IdxCols, 0, len(src[numCol])) // allocate capacity with lenght=0
		}
	}
	return dst
}

func allocValidCols(capacité int) IdxCols {
	return make(IdxCols, 0, capacité)
}

func allocfilteredCols(n int) IdxColsSet {
	return make(IdxColsSet, n)
}
