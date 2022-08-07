package tabjeu

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

/*
type étatBase int
type étatJoué int
*/


const TAILLE_MAX = 50   // longueur fixe pour tenter des allocations sur la stack au lieu du heap

type Cellule byte

const (
	Vide  Cellule = 0
	Plein Cellule = 1
)

/*
const (
	Blanc   étatJoué = 0
	coché   étatJoué = 1
	colorié étatJoué = 2
)
*/

/*
type CelluleBase interface {
	EstPlein() bool
}
*/

// TabJeu contient les cellules du jeu sous forme de slice 2D
type LigneJeu []Cellule
type TabJeu []LigneJeu

/*
func (c Cellule) EstPlein() bool {
	return c == Plein
}
*/

// NewTabJeu crée un nouveau tableau de jeu
func NewTabJeu(taille int, ratioRemplissage int, seed int64) TabJeu {
	if seed == 0 {
		seed = time.Now().Unix()
	}
	rand.Seed(seed)

	// tableau de jeu représenté par un slice à 2 dimensions
	tj := make(TabJeu, taille)

	// alloue toutes les cellules en une seule fois
	cellules := make(LigneJeu, taille*taille)

	// colorie certaines cellules
	for i:= range cellules {
		if rand.Intn(100) < ratioRemplissage {
			cellules[i] = Plein
		}
	}
	// construit tabjeu en découpant en lignes les cellules alouées
	for l := range tj {
		tj[l], cellules = cellules[:taille], cellules[taille:]
	}
	return tj
}

func (tj TabJeu) StringsSlice() []string {
	tjStrings := make([]string, len(tj))
	for l, ligneCells := range tj {
		ligneStrings := make([]string, len(tj))
		for col, cell := range ligneCells {
			ligneStrings[col] = cell.String()
		}
		tjStrings[l] = "[" + strings.Join(ligneStrings, " ") + "]"
	}
	return tjStrings
}

func (sc BlocCount) String() string {
	var elems []string
	for _, ints := range sc {
		elems = append(elems, fmt.Sprintf("%2d", ints))
	}
	return strings.Join(elems, " ")
}

func (tbc TransposedBlocCount) String() string {
	bc := BlocCount(tbc)
	return bc.String()
}

func (tj TabJeu) String() string {
	return strings.Join(tj.StringsSlice(), "\n") + "\n"
}

func (c Cellule) String() string {
	s := "  "  // deux espaces pour une cellule vide
	if c == Plein {
		s = "\u2588\u2588"
	}
	return s
}
