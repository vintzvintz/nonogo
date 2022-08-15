package level

import (
	"fmt"
	"runtime"
	"time"

	"vintz.fr/nonogram/solver"
	"vintz.fr/nonogram/tabjeu"
)

type Level tabjeu.TabJeu

func NewLevel(size int, ratio float32, seed int64) Level {

	tj := tabjeu.NewTabJeu(size, ratio, seed)
	valide := desambigueTabJeu(tj, -1)
	return Level(valide)
}

func NewDefault() Level {
	return NewLevel(tabjeu.DEFAULT_SIZE, tabjeu.DEFAULT_RATIO, tabjeu.DEFAULT_SEED)
}

func desambigueTabJeu(tj tabjeu.TabJeu, nbWorkers int) (valide tabjeu.TabJeu) {

	taille := len(tj)
	startTime := time.Now()
	if nbWorkers < 0 {
		nbWorkers = runtime.NumCPU()
	}

	var nbSol int
	var prev tabjeu.TabJeu
	var diff tabjeu.Diff = tabjeu.NewDiff(taille)
	for solution := range solver.SolveConcurrent(tj, nbWorkers, true) {
		if prev != nil {
			solution.Compare(prev, diff)
		}
		prev = solution
		nbSol++
	}

	duree := time.Since(startTime)
	fmt.Printf(" %d solutions et %d cellules ambigües. %d workers en %v\n", nbSol, diff.Count(), nbWorkers, duree)
	
	valide = tj.Copy()
	valide.ReveleVides(diff)

	return valide
}

func (lvl Level) String() string {
	return tabjeu.TabJeu(lvl).String()
}