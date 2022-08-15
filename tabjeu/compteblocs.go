package tabjeu

import (
	"fmt"
)

type Diff [][]bool

// Direction de comptage des séquences
const (
	LIGNE   = 1
	COLONNE = 2
)

type Probleme struct {
	Taille      int
	BlocsLignes   []BlocCount
	BlocsColonnes []BlocCount
}

func (tj *TabJeu) MakeProbleme() Probleme {
	return Probleme{
		Taille:      len(*tj),
		BlocsLignes:   tj.CompteBlocs(LIGNE),
		BlocsColonnes: tj.CompteBlocs(COLONNE),
	}
}

// CompteBlocs compte les blocs dans la direction indiquée
// renvoie une liste  de longueurs des blocs cellules pleines consécutives
func (tj TabJeu) CompteBlocs(direction int) BlocCountList {

	taille := len(tj)
	resultat := make([]BlocCount, taille)
	// index i : ligne en mode 'ligne', ou colonne en mode 'colonne'
	// index j : colonne en mode 'ligne' ou colonne en mode 'ligne'
	for i := 0; i < taille; i++ {

		blocs := make(BlocCount, 0, taille/2 +1 ) // le nb max de blocs est (taille/2 + 1) car il faut  des separateurs
		var long int                    // longueur du bloc courant
		var cell Cellule
		for j := 0; j < taille; j++ {

			switch direction {
			case LIGNE:
				cell = tj[i][j]
			case COLONNE:
				cell = tj[j][i]
			default:
				panic(fmt.Sprintf("Direction %d inconnue", direction))
			}
			if cell.EstPlein() {
				long++
			}

			// fin d'un bloc ou fin de la ligne
			if ( !cell.EstPlein() && long > 0) || j == taille-1 {
				if long > 0 {
					blocs = append(blocs, long)
					long = 0
				}
			}
		}
		resultat[i] = blocs
	}
	return resultat
}

// CompteBlocsCompare compte les blocs dans la direction indiquée
// renvoie vrai si les blocs sont identiques à la référence
// version optimisée de CompteBlocs sans allocations
func (tj TabJeu) CompareBlocsColonnes(blocsRef BlocCountList) bool {

	taille := len(tj)
	for i := 0; i < taille; i++ {

		var long int // longueur du bloc courant
		var rang int // rang du bloc courant
		for j := 0; j < taille; j++ {

			// inversion i/j pour parcourir dans l'axes des colonnes
			var cell Cellule = tj[j][i]

			// debut ou continuation d'un bloc
			if cell.EstPlein() {
				long++
			}
			// fin d'un bloc ou fin de la ligne
			if ( !cell.EstPlein() && long > 0) || j == taille-1 {
				// compare avec la longueur du bloc de référence
				if (rang < len(blocsRef[i])) && (long != blocsRef[i][rang]) {
					return false
				}
				// prepare le comptage du bloc suivant
				rang++
				long = 0
			}
		}
	}
	return true
}

// NewDiff() alloue un tableau de booléens pour recevoir la comparaison de deux tabjeu.TabJeu
func NewDiff(taille int) (diff Diff) {
	diff = make(Diff, taille)
	for i := range diff {
		diff[i] = make([]bool, taille)
	}
	return diff
}

// Count() compte le nombre d'elements 'true'
func (diff Diff) Count() (nb int) {
	for _, ligne := range diff{
		for _, cell := range ligne {
			if cell {
				nb++
			}
		}
	}
	return nb
}

// Compare() compare tj avec ref et renvoie un tableau avec 'true' pour chaque cellule différente
func (tj TabJeu) Compare(ref TabJeu, diff Diff) {

	taille := len(tj)

	if len(ref) != taille  || len(diff)!=taille {
		panic("Nombre de lignes incohérent")
	}

	//compare chaque cellule et note les différences dans *diff
	for i, ligne := range tj {

		if len(ligne)!=taille || len(ref[i])!=taille || len(diff[i])!=taille {
			panic("Nombre de colonnes incohérent")
		}
	
		for j, cell := range ligne {
			if cell.EstPlein() != ref[i][j].EstPlein() {
				diff[i][j] = true
			}
		}
	}
}


func (tj TabJeu) AfficheAvecComptes() {

	blocsLignes := tj.CompteBlocs(LIGNE)
	blocsCols := tj.CompteBlocs(COLONNE)
	blocsColsT := blocsCols.transpose()

	for _, l := range blocsColsT {
		fmt.Printf("%v\n", l)
	}
	for i := range tj {
		fmt.Printf("%v %v\n", tj[i], blocsLignes[i])
	}

}
