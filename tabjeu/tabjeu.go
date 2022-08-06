package tabjeu

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type étatBase int
type étatJoué int

const (
	Vide  étatBase = 0
	Plein étatBase = 1
)

const (
	Blanc   étatJoué = 0
	coché   étatJoué = 1
	colorié étatJoué = 2
)

type Cellule struct {
	Base étatBase
	Joué étatJoué
}

type CelluleBase interface {
	EstPlein() bool
}

// TabJeu contient les cellules du jeu sous forme de slice 2D
type LigneJeu []*Cellule
type TabJeu []LigneJeu

func (c Cellule) EstPlein() bool {
	return c.Base == Plein
}

// NewTabJeu crée un nouveau tableau de jeu
func NewTabJeu(taille int, ratioRemplissage int, seed int64) TabJeu {
	if seed == 0 {
		seed = time.Now().Unix()
	}
	rand.Seed(seed)

	// tableau de jeu représenté par un slice à 2 dimensions
	tj := make(TabJeu, taille)

	for l := range tj {

		ligne := make([]*Cellule, taille)
		for col := range ligne {

			// tj[i], cellules = cellules[:taille], cellules[taille:]
			// alloue les cellules une par une et colorie certaines aleatoirement
			// TODO allouer les cellules en une seule fois
			cell := new(Cellule)
			if rand.Intn(100) < ratioRemplissage {
				cell.Base = Plein
			}
			ligne[col] = cell
		}
		tj[l] = ligne
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
	tBase := [...]rune{' ', '\u2588'}

	s := fmt.Sprintf("%c%c", tBase[c.Base], tBase[c.Base])
	return s
}
