package tabjeu

import (
	"fmt"
)

type BlocCount []int
type TransposedBlocCount []int

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
func (tj TabJeu) CompteBlocs(direction int) []BlocCount {

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
			if cell == Plein {
				long++
			}

			// fin d'un bloc ou fin de la ligne
			if ( cell == Vide && long > 0) || j == taille-1 {
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
func (tj TabJeu) CompareBlocsColonnes(blocsRef []BlocCount) bool {

	taille := len(tj)
	for i := 0; i < taille; i++ {

		var long int // longueur du bloc courant
		var rang int // rang du bloc courant
		for j := 0; j < taille; j++ {

			// inversion i/j pour parcourir dans l'axes des colonnes
			var cell Cellule = tj[j][i]

			// debut ou continuation d'un bloc
			if cell == Plein {
				long++
			}
			// fin d'un bloc ou fin de la ligne
			if ( cell != Plein && long > 0) || j == taille-1 {
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

func (tj TabJeu) AfficheAvecComptes() {

	seqL := tj.CompteBlocs(LIGNE)
	seqC := tj.CompteBlocs(COLONNE)
	seqCTranspo := transposeSeqColonnes(seqC)

	for _, l := range seqCTranspo {
		fmt.Printf("%v\n", l)
	}
	for i := range tj {
		fmt.Printf("%v %v\n", tj[i], seqL[i])
	}

}

// transposeSeqColonnes convertit en lignes les comptes de séquences en colonne
// ceci est utile seulement pour l'affichage
func transposeSeqColonnes(seqC []BlocCount) []TransposedBlocCount {

	seqTranspo := make( []TransposedBlocCount, 0)

	// 	on ne connait pas à l'avance le nb max de blocs
	for rang := 0; ; rang++ {
		var empty bool = true // sort de la boucle quand toutes les séquences en colonnes sont épuisées

		// transpo contient les séquences en colonnes de même rang
		ligneTranspo := make(TransposedBlocCount, len(seqC))

		for i := range seqC {
			// valeur par défaut 0 si plus de séquences dans la colonne
			if rang < len(seqC[i]) {
				empty = false
				ligneTranspo[i] = seqC[i][rang]
			}
		}
		if empty {
			break
		}
		seqTranspo = append(seqTranspo, ligneTranspo)
	}
	return seqTranspo
}
