package solver

import (
	//"fmt"
	"os"
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

	startRecursion := func() {
		tj := make(TjPartiel, 0, prob.Taille)
		solveRecursif(&allBlocs, tj, colonnes, nil, workerPool, solutions, pc)
	}

	// lance la récursion
	if nbWorkers > 0 {
		workerPool.Exec(func() {
			startRecursion()
		})
		// goroutine nécessaire pour attendre la fin du traitment
		go func() {
			workerPool.Wait()
			termine()
		}()
	} else {
		go func() {
			startRecursion()
			termine() // pas la peine d'attendre dans une goroutine distincte du lancement
		}()
	}

	return solutions
}

type TjPartiel []int
type CopyDoneChan chan bool

func solveRecursif(allBlocs *allPossibleBlocs,
	tjPartiel TjPartiel, // index dans allBlocs des lignes déja placées (longueur variable entre 0 et [taille] éléments)
	initialCols IdxColsSet, // index dans allBlocs des colonnes encore valides
	lockCopy CopyDoneChan,
	wp *WorkerPool,
	solutions chan *TJ.TabJeu,
	pc *perf.PerfCounter) {
	//fmt.Printf("#%v tjPartiel ligne %v  %v (lockCopy %v)\n", goid(), numLigneCourante, tjPartiel, lockCopy)

	taille := len(allBlocs.rows)
	numLigneCourante := len(tjPartiel)

	// copie l'etat si cela est demandé par l'appelant (recursion parallele, le parent est une goroutine différente )
	// inutile si le parent est la même goroutine (recursion classique, même goroutine )	
	if lockCopy != nil {
		initialCols = initialCols.Copy()
		tjPartiel = tjPartiel.Copy(taille)
		lockCopy <- true // débloque l'execution du parent
	}

	// met à jour le compteur de vitesse
	if pc != nil {
		pc.Inc1()
	}

	// combinaisons de blocs (en ligne) à essayer sur la ligne courante
	tryLines := allBlocs.rows[numLigneCourante] 

	// allocation d'espace pour recevoir les colonnes valides restantes avec chaque ligne à essayer 
	nextCols := initialCols.AllocEmpty()
	// augmente la longueur de tjPartiel pour recevoir l'index de la ligne en cours d'essai
	tjPartiel = append(tjPartiel,-1)

	// pour attendre la copie de l'état lors de la récursion concurrente (par une autre goroutine)
	waitChild := make(CopyDoneChan) 
	
	// essaye toutes les combinaisons encore valides pour la ligne courante
	for n, nextLigne := range tryLines {

		// parmi les colonnes reçues du parent, elimine celles incompatibles avec nextLine
		ok := filtreColonnesInplace(allBlocs, nextLigne, numLigneCourante, initialCols, nextCols)

		// condition d'arrêt : nextLigne est incompatible avec les colonnes encore valides à ce point de la recursion
		if !ok {
			continue
		}

		// si la ligne est valide, on l'inscrit dans tjPartiel pour continuer la recherche
		tjPartiel[numLigneCourante] = n

		// condition d'arrêt : on a trouvé une solution si toutes les lignes sont remplies
		if numLigneCourante == taille-1 {
			solutions <- tabJeuFromIndexArray(allBlocs, tjPartiel, taille)
			continue
		}
		// Prepare deux versions de l'appel récursif
		recurseNoCopy := func() {
			// pour execution par la même goroutine, sans copie des paramètres (lockCopy=nil) 
			solveRecursif(allBlocs, tjPartiel, nextCols, nil, wp, solutions, pc)
		}
		recurseWithCopy := func() {
			//  pour execution dans une autre goroutine avec copie des parametres
			solveRecursif(allBlocs, tjPartiel, nextCols, waitChild, wp, solutions, pc)
		}

		// Tente la recursion concurrente
		accepted := wp.TryExec(recurseWithCopy)
		if accepted {
			// attend que recurseWithCopy ait fini de copier ses paramètres
			<-waitChild
			continue
		}
		// continue avec une récursion classique si la tâche a été refusée par le workerpool
		recurseNoCopy()
	}
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

		// vérifie si chaque element de initialCols est compatible avec la ligne fournie
		cellLigne = ligne[numCol] // lookup en dehors de la boucle
		for _, idxCol := range initialCols[numCol] {
			cellCol = allBlocs.cols[numCol][idxCol][numLigne]

			// si out, ajoute l'index de la colonne dans le slice de destination
			if cellLigne == cellCol {
				filteredCols[numCol] = append(filteredCols[numCol], idxCol)
			}
		}

		// termine dès qu'une colonne est incompatible avec la ligne testée
		if len(filteredCols[numCol]) == 0 {
			return false
		}
	}
	return true
}

// tabJeuFromIndex construit un tabJeu à partir d'un array d'index de lignes
func tabJeuFromIndexArray(allBlocs *allPossibleBlocs, index TjPartiel, taille int) *TJ.TabJeu {
	tj := make(TJ.TabJeu, taille)
	for n := 0; n < taille; n++ {
		rows := allBlocs.rows[n]
		idx := index[n]
		tj[n] = rows[idx]
	}
	return &tj
}


func (src IdxColsSet) AllocEmpty() (dst IdxColsSet) {
	dst = allocfilteredCols(len(src))
	for numCol := range src {
		dst[numCol] = allocIdxCols(0, len(src[numCol])) // allocate capacity with lenght=0
	}
	return dst
}

func (src IdxColsSet) Copy() (dst IdxColsSet) {
	dst = allocfilteredCols(len(src))
	for numCol := range src {
		dst[numCol] = allocIdxCols(len(src[numCol]), len(src[numCol])) // allocate same length as src
		copy(dst[numCol], src[numCol])
	}
	return dst
}

func (src TjPartiel) Copy(capa int) (dst TjPartiel) {
	dst = allocTjPartiel( len(src), capa )
	copy( dst, src)
	return dst
}


// wrappers pour mieux identifier les allocations avec pprof (devraient être inlinés)
func allocIdxCols(long, capa int) IdxCols {
	return make(IdxCols, long, capa)
}
func allocfilteredCols(n int) IdxColsSet {
	return make(IdxColsSet, n)
}
func allocTjPartiel(long, capa int) TjPartiel {
	return make(TjPartiel, long, capa)
}