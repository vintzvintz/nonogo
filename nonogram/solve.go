package nonogram

import (
	"fmt"
	"sync"
	"time"
)

// lineList est une liste de lignes
// contrairement à TabJeu la largeur n'est pas forcément égale à la hauteur
type lineList []LigneJeu

// allPossibleLines est un slice de la même longueur que TabJeu
// chaque élement est l'ensemble des valeurs possibles pour la ligne correspondante
type lineListSet []lineList

// structure intermédiaire pour construire l'ensemble des lignes possibles
type indexedLineSet struct {
	num   int
	lines lineList
}

// cellVide et cellPlein sont les éléments de base pour construire les sequences possibles
// on aurait aussi bien pu utiliser des booleéns à la place des pointeurs
var cellVide = &cellule{base: vide, joué: blanc}
var cellPlein = &cellule{base: plein, joué: blanc}

type allPossibleBlocs struct {
	rows lineListSet
	cols lineListSet
}

type SolverFunc func(Probleme) chan *TabJeu

var perfCounter PerfCounter

func Bench(prob Probleme) {

	//BenchSolver(SolveBourrin, prob, "Bourrin")

	perfCounter.Start(time.Second)

	tests := []int{8, 6, 4, 2, 1}
	for _, nb := range tests {
		txt := fmt.Sprintf("%d workers", nb)
		solver := makeConcurrentSolver(nb)
		BenchSolver(solver, prob, txt)
	}
	BenchSolver(SolveSimple, prob, "Récursif simple")
}

func BenchSolver(solver SolverFunc, prob Probleme, txt string) {
	solutions := solver(prob)
	var nb int
	startTime := time.Now()
	for sol := range solutions {
		nb++
		//fmt.Printf("Solution n°%d\n", nb)
		//fmt.Print(*sol)
		_ = sol
	}
	duree := time.Since(startTime)
	fmt.Printf("%s : %d solutions trouvées en %v\n", txt, nb, duree)
}

// buildAllSequences construit la liste des ensembles de lignes (ou colonnes)
// à partir d'une liste de listes de longueurs de blocs
func buildAllSequences(taille int, seqs []seqCount) lineListSet {
	result := make(lineListSet, taille)

	ch := make(chan indexedLineSet)
	wg := new(sync.WaitGroup)

	// construit les ensembles de lignes possibles en parallèle
	for i := range seqs {
		wg.Add(1)
		go func(taille int, sc seqCount, idx int) {
			lines := buildSequences(taille, sc)
			ch <- indexedLineSet{num: idx, lines: lines}
		}(taille, seqs[i], i)
	}

	// concatene les resultats pour chaque sequence de longueurs de blocs
	go func() {
		for lines := range ch {
			result[lines.num] = lines.lines
			wg.Done()
		}
	}()

	wg.Wait()
	return result
}

// buildSequences construit l'ensemble des lignes (ou colonnes) possibles
// à partir d'une liste de longueurs de blocs et de la taille de la ligne
func buildSequences(taille int, seqs seqCount) lineList {
	// cas particulier des lignes complètement vides
	if len(seqs) == 0 {
		ligneVide := make(LigneJeu, taille)
		for i := range ligneVide {
			ligneVide[i] = cellVide
		}
		return lineList{ligneVide}
	}

	// cas particulier pour le dernier bloc (pas de séparateur à la fin)
	lastBloc := len(seqs) == 1

	// calcule le nb de cellules mini pour placer tous les blocs avec un espacement de 1
	var longMini int = len(seqs) - 1 // nb de cellules vides intercalaires
	for _, s := range seqs {
		longMini += s
	}

	result := make(lineList, 0)

	// essaye succesivement le bloc sur toutes les positions possibles
	for startPos := 0; startPos <= (taille - longMini); startPos++ {

		tailleSeqCourante := taille // dernier blocs
		if !lastBloc {
			//  non-derniers blocs
			tailleSeqCourante = startPos + seqs[0] + 1
		}
		seqCourante := make(LigneJeu, tailleSeqCourante)
		for i := 0; i < tailleSeqCourante; i++ {
			// place des cellules pleines entre startPos et la fin du bloc
			if (startPos) <= i && (i < startPos+seqs[0]) {
				seqCourante[i] = cellPlein
				continue
			}
			// place des cellules vides ailleurs (avant et/ou après le bloc)
			seqCourante[i] = cellVide
		}

		// Si c'est le dernier bloc on renvoie juste les séquences courantes
		if lastBloc {
			result = append(result, seqCourante)
		}

		// si ce n'est pas le dernier bloc, appel récursif pour les séquences restantes
		if !lastBloc {
			seqsSuivantes := buildSequences(taille-tailleSeqCourante, seqs[1:])

			//  concatène les séquences suivantes avec la séquence courante
			for i := range seqsSuivantes {
				seq := append(seqCourante, seqsSuivantes[i]...)
				result = append(result, seq)
			}
		}
	}
	return result
}

func (l lineList) String() string {
	var str string = ""
	for _, line := range l {
		str = str + fmt.Sprintf("%v\n", line)
	}
	return str
}

func (pl indexedLineSet) String() string {
	var str string = fmt.Sprintf("Sequences possibles pour la ligne %d\n", pl.num)
	str += pl.lines.String()
	return str
}

func (all lineListSet) String() string {
	var str string
	for _, pl := range all {
		str += pl.String()
	}
	return str
}
