package solver

import (
	"fmt"
	"sync"
	TJ "vintz.fr/nonogram/tabjeu"
)

// lineList est une liste de lignes
// contrairement à TabJeu la largeur n'est pas forcément égale à la hauteur
type lineList []TJ.LigneJeu

// allPossibleLines est un slice de la même longueur que TabJeu
// chaque élement est l'ensemble des valeurs possibles pour la ligne correspondante
type lineListSet []lineList

// structure intermédiaire pour construire l'ensemble des lignes possibles
type indexedLineSet struct {
	num   int
	lines lineList
}

type allPossibleBlocs struct {
	rows lineListSet
	cols lineListSet
}

func buildPossibleBlocs(tj TJ.TabJeu) allPossibleBlocs {

	// compte les blocs de cellules consecutuves dans tj
	// construit toutes les lignes et colonnes possibles avec les longueurs de blocs mesurées
	// conserve uniquement celles compatibles avec les cellules révéléés dans tj

	bcLigne := tj.CompteBlocs(TJ.LIGNE)
	allRows := buildAllSequences(bcLigne)
	possibleRows := filtreSequences(allRows, tj, TJ.LIGNE)

	bcCol := tj.CompteBlocs(TJ.COLONNE)
	allCols := buildAllSequences(bcCol)
	possibleCols := filtreSequences(allCols, tj, TJ.COLONNE)

	return allPossibleBlocs{
		rows: possibleRows,
		cols: possibleCols,
	}
}

// filtreSequences renvoie les elements de all compatibles avec les cellules révélées dans tj
func filtreSequences(all lineListSet, tj TJ.TabJeu, direction TJ.Direction) (filtered lineListSet) {

	taille := len(tj)
	if len(all) != taille {
		panic("Incohérence des dimensions")
	}

	filteredSet := make(lineListSet, taille)
	for i, lignes := range all { // de 0  à taille-1
		filteredLines := make(lineList, 0)
	iterlines:
		for _, ligne := range lignes { // itere sur toutes les possibilités pour la ligne ou colonne i
			for j, cell := range ligne { // itere sur les cellules dans l'autre dimension ( colonne ou ligne )

				// récupere la cellule correspondante dans tj
				var refCell TJ.Cellule
				switch direction {
				case TJ.LIGNE:
					refCell = tj[i][j]
				case TJ.COLONNE:
					refCell = tj[j][i]
				}

				// si la cellule dans tj est révélée, elle doit être identique
				// si elle ne l'est pas, on ignore la ligne courante
				if refCell.EstRévélé() && (refCell.EstPlein() != cell.EstPlein()) {
					continue iterlines
				}
			}
			filteredLines = append(filteredLines, ligne)
		}
		filteredSet[i] = filteredLines
	}
	return filteredSet
}

// buildAllSequences construit la liste des ensembles de lignes (ou colonnes)
// à partir d'une liste de listes de longueurs de blocs
func buildAllSequences(blocs []TJ.BlocCount) lineListSet {
	taille := len(blocs)
	result := make(lineListSet, taille)

	ch := make(chan indexedLineSet)
	wg := new(sync.WaitGroup)

	// construit les ensembles de lignes possibles en parallèle
	for i := range blocs {
		wg.Add(1)
		go func(taille int, sc TJ.BlocCount, idx int) {
			lines := buildSequences(taille, sc)
			ch <- indexedLineSet{num: idx, lines: lines}
		}(taille, blocs[i], i)
	}

	// reçoit les resultats pour chaque ligne
	go func() {
		for lines := range ch {
			result[lines.num] = lines.lines
			wg.Done()
		}
	}()

	wg.Wait()
	close(ch)
	return result
}

// buildSequences construit recursivement toutes les combinaisons possibles
// correspondant à une liste de longueurs de blocs pour une ligne ou une colonne
func buildSequences(taille int, blocs TJ.BlocCount) lineList {
	// cas particulier des lignes complètement vides
	if len(blocs) == 0 {
		ligneVide := make(TJ.LigneJeu, taille) // zero-value = Vide
		return lineList{ligneVide}
	}

	// cas particulier pour le dernier bloc (pas de séparateur à la fin)
	lastBloc := len(blocs) == 1

	// calcule le nb de cellules mini pour placer tous les blocs avec un espacement de 1
	var longMini int = len(blocs) - 1 // nb de cellules vides intercalaires
	for _, s := range blocs {
		longMini += s
	}

	result := make(lineList, 0)

	// place succesivement le bloc sur toutes les positions possibles
	for startPos := 0; startPos <= (taille - longMini); startPos++ {

		tailleSeqCourante := taille // dernier blocs
		if !lastBloc {
			//  non-derniers blocs
			tailleSeqCourante = startPos + blocs[0] + 1
		}
		seqCourante := make(TJ.LigneJeu, tailleSeqCourante)
		for i := 0; i < tailleSeqCourante; i++ {
			// par défaut, toutes les cellules sont initialisées vides, inutile de les remplir explicitement
			//seqCourante[i] = cellVide

			// remplit la ligne à partir de startPos jusqu'à la fin du bloc
			if (startPos) <= i && (i < startPos+blocs[0]) {
				seqCourante[i].Remplit()
				continue
			}
		}

		// Si c'est le dernier bloc on renvoie juste les séquences courantes
		if lastBloc {
			result = append(result, seqCourante)
		}

		// si ce n'est pas le dernier bloc, appel récursif pour les séquences restantes
		if !lastBloc {
			seqsSuivantes := buildSequences(taille-tailleSeqCourante, blocs[1:])

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

func (all lineListSet) String() string {
	var str string
	for _, pl := range all {
		str += pl.String()
	}
	return str
}

func (pl indexedLineSet) String() string {
	var str string = fmt.Sprintf("Sequences possibles pour la ligne %d\n", pl.num)
	str += pl.lines.String()
	return str
}
