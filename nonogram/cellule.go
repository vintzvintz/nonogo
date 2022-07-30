package nonogram

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type étatBase int
type étatJoué int

const (
	vide  étatBase = 0
	plein étatBase = 1
)

const (
	blanc   étatJoué = 0
	coché   étatJoué = 1
	colorié étatJoué = 2
)

type cellule struct {
	base étatBase
	joué étatJoué
}

type CelluleBase interface {
	estPlein() bool
}

// TabJeu contient les cellules du jeu sous forme de slice 2D
type LigneJeu []*cellule
type TabJeu []LigneJeu

func (c cellule) estPlein() bool {
	return c.base == plein
}

// NewTabJeu crée un nouveau tableau de jeu
func NewTabJeu(taille int, ratioRemplissage int) TabJeu {
	rand.Seed(time.Now().Unix())

	// tableau de jeu représenté par un slice à 2 dimensions
	tj := make(TabJeu, taille)

	for l := range tj {

		ligne := make([]*cellule, taille)
		for col := range ligne {

			// tj[i], cellules = cellules[:taille], cellules[taille:]
			// alloue les cellules une par une et colorie certaines aleatoirement
			// TODO allouer les cellules en une seule fois
			cell := new(cellule)
			if rand.Intn(100) < ratioRemplissage {
				cell.base = plein
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

func (sc seqCount) String() string {
	var elems []string
	for _, ints := range sc {
		elems = append(elems, fmt.Sprintf("%2d", ints))
	}
	return strings.Join(elems, " ")
}

func (tj TabJeu) String() string {
	return strings.Join(tj.StringsSlice(), "\n") + "\n"
}

func (c cellule) String() string {
	tBase := [...]rune{' ', '\u2588'}

	s := fmt.Sprintf("%c%c", tBase[c.base], tBase[c.base])
	return s
}
