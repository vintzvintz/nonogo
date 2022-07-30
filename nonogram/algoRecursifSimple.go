package nonogram

import (
	"fmt"
	"sync"
	"time"
)

//"fmt"
//"sync"
//"time"

func SolveSimple(prob Probleme) chan *TabJeu {
	solutions := make(chan *TabJeu)
	allBlocs := allPossibleBlocs{
		rows: buildAllSequences(prob.taille, prob.seqLignes),
		cols: buildAllSequences(prob.taille, prob.seqColonnes),
	}

	go func() {
		solveRecursifSimple(&allBlocs, nil, allBlocs.cols, solutions)
		close(solutions)
	}()

	return solutions
}

func solveRecursifSimple(allBlocs *allPossibleBlocs,
	tjPartiel TabJeu,
	colonnes lineListSet,
	//perf PerfCounter,
	solutions chan *TabJeu) {

	taille := len(allBlocs.rows)
	idxLigneCourante := len(tjPartiel)
	tryLines := (*allBlocs).rows[idxLigneCourante]

	for n, nextLigne := range tryLines {

		// cree une copie du tj partiel avec une ligne supplémentaire
		tjNext := make(TabJeu, idxLigneCourante+1)
		copy(tjNext, tjPartiel)
		tjNext[idxLigneCourante] = nextLigne

		// elimine les colonnes incompatibles avec la ligne ajoutée
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
		// sinon on continue avec la ligne suivante
		solveRecursifSimple(allBlocs, tjNext, nextColonnes, solutions)

	}
}

func filtreColonnes(tjPartiel TabJeu, colonnes lineListSet) (filteredCols lineListSet, ok bool) {

	perfCounter.Inc()

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
	// implicit return
	return filteredCols, true
}

type PerfCounter struct {
	lock  sync.Mutex
	count int64
	//	lastCount int64
	//	lastTime time.Time
}

func (pc *PerfCounter) Inc() {
	pc.lock.Lock()
	defer pc.lock.Unlock()
	pc.count++
}

func (pc *PerfCounter) Get() int64 {
	pc.lock.Lock()
	defer pc.lock.Unlock()
	return pc.count
}

func (pc *PerfCounter) Start(interval time.Duration) {
	go func() {
		lastTime := time.Now()
		lastCount := pc.Get()
		for {
			time.Sleep(interval)
			now := time.Now()
			count := pc.Get()
			rate := (count - lastCount) * 1e6 / now.Sub(lastTime).Nanoseconds()
			lastTime, lastCount = now, count
			fmt.Printf("Vitesse %d k/s\n", rate)
		}
	}()
}
