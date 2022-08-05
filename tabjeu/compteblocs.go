package tabjeu

import (
	"fmt"
)

type SeqCount []int

// Direction de comptage des séquences
const (
	LIGNE   = 1
	COLONNE = 2
)

type Probleme struct {
	Taille      int
	SeqLignes   []SeqCount
	SeqColonnes []SeqCount
}

func (tj *TabJeu) MakeProbleme() Probleme {
	return Probleme{
		Taille:      len(*tj),
		SeqLignes:   tj.CompteBlocs(LIGNE),
		SeqColonnes: tj.CompteBlocs(COLONNE),
	}
}

// CompteBlocs compte les blocs dans la direction indiquée
// renvoie une liste  de longueurs des blocs cellules pleines consécutives
func (tj TabJeu) CompteBlocs(direction int) []SeqCount {

	taille := len(tj)
	resultat := make([]SeqCount, taille)
	// index i : ligne en mode 'ligne', ou colonne en mode 'colonne'
	// index j : colonne en mode 'ligne' ou colonne en mode 'ligne'
	for i := 0; i < taille; i++ {

		blocs := make(SeqCount, taille) // le nb max de blocs en réalité est (taille/2 + 1)
		var long int                    // longueur du bloc courant
		var rang int                    // rang du bloc courant
		for j := 0; j < taille; j++ {

			var cell CelluleBase
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
			if (!cell.EstPlein() && long > 0) || j == taille-1 {
				if long > 0 {
					blocs[rang] = long
					rang++
					long = 0
				}
			}
		}
		// tronque la liste des blocs à la bonne longueur
		blocs = blocs[:rang]
		resultat[i] = blocs
	}
	return resultat
}

// CompteBlocsCompare compte les blocs dans la direction indiquée
// renvoie vrai si les blocs sont identiques à la référence
// version optimisée de CompteBlocs sans allocations
func (tj TabJeu) CompteBlocsCompare(seqRef []SeqCount) bool {

	taille := len(tj)
	for i := 0; i < taille; i++ {

		var long int // longueur du bloc courant
		var rang int // rang du bloc courant
		for j := 0; j < taille; j++ {

			var cell CelluleBase = tj[j][i]

			// debut ou continuation d'un bloc
			if cell.EstPlein() {
				long++
			}

			// fin d'un bloc ou fin de la ligne
			if (!cell.EstPlein() && long > 0) || j == taille-1 {
				// compare avec la référence
				if (rang < len(seqRef[i])) && (long != seqRef[i][rang]) {
					return false
				}
				// prepare pour compter  le bloc suivant
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
// ceci est utile surtout pour l'affichage
func transposeSeqColonnes(seqC []SeqCount) []SeqCount {
	// Transposition des séquences en colonnes
	seqTranspo := make([]SeqCount, 0)

	// 	on ne connait pas à l'avance le nb max de séquences
	for rang := 0; ; rang++ {
		var empty bool = true // sort de la boucle quand toutes les séquences en colonnes sont épuisées

		// transpo contient les séquences en colonnes de même rang
		ligneTranspo := make(SeqCount, len(seqC))

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
