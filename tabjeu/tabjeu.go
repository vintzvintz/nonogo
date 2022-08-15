package tabjeu

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const DEFAULT_SIZE = 15
const DEFAULT_SEED = 1005
const DEFAULT_RATIO = 0.44

type Cellule uint8

const (
	PLEIN     Cellule = 1 << iota // Etat de base de la cellule : vide ou plein
	REVELE                        // L'état est révélé dès le début du jeu
	JOUEVIDE                      // La cellule est jouée vide
	JOUEPLEIN                     // La cellule est jouée pleine
)

// TabJeu contient les cellules du jeu sous forme de slice 2D
type LigneJeu []Cellule
type TabJeu []LigneJeu
type BlocCount []int
type BlocCountList []BlocCount
type TransposedBlocCount []int

const SEP = ""

func (c Cellule) EstPlein() bool {
	return c&PLEIN == PLEIN
}

func (c *Cellule) Remplit() {
	*c = *c | PLEIN
}

func (c Cellule) EstRévélé() bool {
	return c&REVELE == REVELE
}

func (c *Cellule) Révèle() {
	*c = *c | REVELE
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
	nbPlein := int(float32(taille*taille) * ratioRemplissage)
	if nbPlein > taille*taille {
		panic("Ratio de remplissage trop élevé")
	}

	// remplit le nb de cellule requis
	for n := 0; n < nbPlein; n++ {
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
	str := make([]string, len(tj))
	for l, ligne := range tj {
		str[l] = ligne.String()
	}
	return str
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
	return strings.Join(tj.StringsSlice(), "\n") + "\n"
}

func (lj LigneJeu) String() string {
	var str = make([]string, len(lj))
	for i, cell := range lj {
		str[i] = cell.String()
	}
	return "[" + strings.Join(str, SEP) + "]"
}

func (c Cellule) String() string {
	if c.EstPlein() {
		return "\u2588\u2588"
	}
	return "  " // deux espaces pour une cellule vide
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
