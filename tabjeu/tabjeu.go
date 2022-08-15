package tabjeu

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

/*
struct remplacée par un uint8 et des bitmasks => gain de 20% à 30% sur la vitesse du solveur...
type Cellule struct {
	plein bool     // Etat de base de la cellule : vide ou plein
	revelé bool    // L'état est révélé dès le début du jeu
	jouéVide bool  // La cellule est jouée vide
	jouéPlein bool // La cellule est jouée pleine
}
*/


type Cellule uint8

const (
	plein Cellule = 1 << iota+1   // Etat de base de la cellule : vide ou plein
	revelé               // L'état est révélé dès le début du jeu
	jouéVide             // La cellule est jouée vide
	jouéPlein            // La cellule est jouée pleine
)

// TabJeu contient les cellules du jeu sous forme de slice 2D
type LigneJeu []Cellule
type TabJeu []LigneJeu


func (c Cellule) EstPlein() bool {
	return c & plein != 0
}

func (c *Cellule) Remplit() {
	*c = *c | plein 
}


// NewTabJeu crée un nouveau tableau de jeu
func NewTabJeu(taille int, ratioRemplissage float32, seed int64) TabJeu {
	if seed == 0 {
		seed = time.Now().Unix()
	}
	rand.Seed(seed)

	// alloue toutes les cellules en une seule fois
	cellules := make(LigneJeu, taille*taille)

	// calcule le nb de cellules à remplir
	nbPlein := int( float32(taille * taille) * ratioRemplissage )
	if nbPlein > taille*taille {
		panic( "Ratio de remplissage trop élevé")
	}

	// remplit le nb de cellule requis
	for n:=0; n<nbPlein; n++ {
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
	if c.EstPlein() {
		s = "\u2588\u2588"
	}
	return s
}
