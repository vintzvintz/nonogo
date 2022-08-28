package tabjeu

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const DEFAULT_SIZE = 10
const DEFAULT_SEED = 1006
const DEFAULT_RATIO = 0.45

// Direction de comptage des séquences
type Direction int
const (
	LIGNE   = 1
	COLONNE = 2
)

// TabJeu contient les cellules du jeu sous forme de slice 2D
type LigneJeu []Cellule
type TabJeu []LigneJeu
type BlocCount []int
type BlocCountList []BlocCount
type TransposedBlocCount []int    // pour l'affichage en mode console
type Diff [][]bool           // comparaison de deux tabjeu

const SEP = ""

// NewTabJeu crée un nouveau tableau de jeu
func NewTabJeu(taille int, nbPlein int, seed int64) TabJeu {

	// alloue toutes les cellules en une seule fois
	cellules := make(LigneJeu, taille*taille)

	// initialise le pseudo-random
	if seed == 0 {
		seed = time.Now().Unix()
	}
	rand.Seed(seed)

	if nbPlein > taille*taille {
		panic("Ratio de remplissage trop élevé")
	}

	// remplit le nb de cellule requis
	for n := 0; n < nbPlein; n++ {
		// recommence tant qu'on ne tombe pas sur une cellule vide
		for {
			i := rand.Intn(len(cellules))
			if !cellules[i].EstPlein() {
				cellules[i].Remplit()
				break
			}
		}
	}

	// tabjeu est un slice à 2 dimensions construit en découpant le bloc de cellules alouées
	tj := make(TabJeu, taille)
	for l := range tj {
		tj[l], cellules = cellules[:taille], cellules[taille:]
	}
	return tj
}

func (tj TabJeu) ReveleVides(ambigu Diff) {
	for i := range tj {
		for j := range tj[i] {
			if ambigu[i][j] && !tj[i][j].EstPlein() {
				tj[i][j].Révèle()
				fmt.Printf("Révèle cellule (%d,%d)\n", i, j)
			}
		}
	}
}

func (tj TabJeu) ComptePlein() (joué, total int) {
	for _, ligne := range tj {
		for _, cell := range ligne {
			if cell.EstPlein() {
				total++
				if cell.EstRévélé() || cell.EstJouéPlein() {
					joué++
				}
			}
		}
	}
	return joué, total
}

func (src TabJeu) Copy() (dst TabJeu) {
	dst = NewTabJeu(len(src), 0, 0)
	for i := range src {
		copy(dst[i], src[i])
	}
	return dst
}

// CompteBlocs compte les blocs dans la direction indiquée
// renvoie une liste  de longueurs des blocs cellules pleines consécutives
func (tj TabJeu) CompteBlocs(direction Direction) BlocCountList {

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

// Compare() compare tj avec ref et met à jour diff avec 'true' pour chaque cellule différente
func (tj TabJeu) Compare(ref TabJeu, diff Diff) (nbDiff int) {

	taille := len(tj)

	if len(ref) != taille  || (diff!=nil && len(diff)!=taille) {
		panic("Nombre de lignes incohérent")
	}

	//compare chaque cellule et note les différences dans diff
	for i, ligne := range tj {

		if len(ligne)!=taille || len(ref[i])!=taille || (diff!=nil && len(diff[i])!=taille) {
			panic("Nombre de colonnes incohérent")
		}
	
		for j, cell := range ligne {
			if cell.EstPlein() != ref[i][j].EstPlein() {
				if diff!=nil {
					diff[i][j] = true
				}
				nbDiff++
			}
		}
	}
	return nbDiff
}


func (sc BlocCount) String() string {
	var elems []string
	for _, count := range sc {
		if count == 0 {
			elems = append(elems, "  ")
			continue
		}
		elems = append(elems, fmt.Sprintf("%2d", count))
	}
	return strings.Join(elems, SEP)
}

func (tbc TransposedBlocCount) String() string {
	bc := BlocCount(tbc)
	return bc.String()
}

func (tj TabJeu) String() string {

	blocsLignes := tj.CompteBlocs(LIGNE)
	blocsCols := tj.CompteBlocs(COLONNE)
	blocsColsT := blocsCols.transpose()

	var str string
	for _, l := range blocsColsT {
		str += fmt.Sprintf("%v\n", l)
	}
	for i := range tj {
		str += fmt.Sprintf("%v %v\n", tj[i], blocsLignes[i])
	}
	return str
}

func (lj LigneJeu) String() string {
	var str = make([]string, len(lj))
	for i, cell := range lj {
		str[i] = cell.String()
	}
	return "[" + strings.Join(str, SEP) + "]"
}

// transposeSeqColonnes convertit en lignes les comptes de séquences en colonne
func (blocs BlocCountList) transpose() []TransposedBlocCount {

	seqTranspo := make([]TransposedBlocCount, 0)

	for rang := 0; ; rang++ {
		var empty bool = true // sort de la boucle quand toutes les séquences en colonnes sont épuisées

		// transpo contient les séquences en colonnes de même rang
		ligneTranspo := make(TransposedBlocCount, len(blocs))

		for i := range blocs {
			// commence par le dernier bloc (bas du tableau)
			rang_inverse := len(blocs[i]) - 1 - rang
			if rang_inverse >= 0 {
				empty = false
				ligneTranspo[i] = blocs[i][rang_inverse]
			}
		}
		if empty {
			break
		}
		// 	on ne connait pas à l'avance le nb max de blocs
		// donc on alloue/copie pour chaque rangée avec inversion de l'ordre
		prev := seqTranspo
		seqTranspo = []TransposedBlocCount{ligneTranspo}
		seqTranspo = append(seqTranspo, prev...)
	}
	return seqTranspo
}
