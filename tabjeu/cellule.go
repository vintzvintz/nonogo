
package tabjeu

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
type TransposedBlocCount []int    // pour l'affichage en mode console

const SEP = ""

func (c Cellule) EstPlein() bool { return c&PLEIN == PLEIN }
func (c *Cellule) Remplit() { *c = *c | PLEIN }

func (c Cellule) EstRévélé() bool { return c&REVELE == REVELE }
func (c *Cellule) Révèle() { *c = *c | REVELE }

func (c Cellule) EstJouéPlein() bool { return c & JOUEPLEIN == JOUEPLEIN}
func (c Cellule) EstJouéVide() bool { return c & JOUEVIDE == JOUEVIDE}
func (c *Cellule) JoueAucun() { *c = *c &^ ( JOUEVIDE | JOUEPLEIN) }
func (c *Cellule) JoueVide() { *c = *c | JOUEVIDE ;	*c = *c &^ JOUEPLEIN }
func (c *Cellule) JouePlein() {	*c = *c | JOUEPLEIN; *c = *c &^ JOUEVIDE }


const (
	IMG_AUCUN int  = iota
	IMG_VIDE
	IMG_PLEIN
)

// Image détermine l'affichage de la cellule à présenter au joueur
func (c Cellule) Image() (img int) {
	if c.EstJouéPlein() || (c.EstRévélé() && c.EstPlein()) {
		return IMG_PLEIN
	}
	if c.EstJouéVide() || (c.EstRévélé() && !c.EstPlein()) {
		return IMG_VIDE
	}
	return IMG_AUCUN
}

func (c Cellule) String() string {
	if c.EstRévélé() {
		return ".."
	}
	if c.EstPlein() {
		return "\u2588\u2588"
	}
	return "  " // deux espaces pour une cellule vide
}
