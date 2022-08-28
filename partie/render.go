package partie

import (
	"vintz.fr/nonogram/tabjeu"
)

type RenderPartie Partie

type RenderLigne struct {
	Blocs    tabjeu.BlocCount
	Cellules []RenderCellule
}

type RenderCellule struct {
	Colonne, Ligne int
	Img            string
	Lien           string
}

func (p RenderPartie) Taille() int {
	return len(p.tj)
}

func (p RenderPartie) DimensionCellules() float32 {
	return 50 / float32(p.Taille())
}

func (p RenderPartie) Seed() int64 {
	return p.seed
}

func (p RenderPartie) PremiereLigne() tabjeu.BlocCountList {
	return p.tj.CompteBlocs(tabjeu.COLONNE)
}

func (p RenderPartie) Lignes() []RenderLigne {

	blocsLignes := p.tj.CompteBlocs(tabjeu.LIGNE)

	var rLignes []RenderLigne
	for l, ligne := range p.tj {
		var rCells []RenderCellule

		for c, cell := range ligne {
			var img string
			switch cell.Image() {
			case tabjeu.IMG_AUCUN:
				img = SVG_AUCUN
			case tabjeu.IMG_VIDE:
				img = SVG_VIDE
			case tabjeu.IMG_PLEIN:
				img = SVG_PLEIN
			default:
				panic("tabjeu.Image() renvoie une valeur inattendue")
			}

			var lien string
			if !cell.EstRévélé() {
				lien = ACTION_CLIC
			}

			rCell := RenderCellule{
				Ligne:   l,
				Colonne: c,
				Img:     img,
				Lien:    lien,
			}
			rCells = append(rCells, rCell)
		}

		rLigne := RenderLigne{
			Blocs:    blocsLignes[l],
			Cellules: rCells,
		}
		rLignes = append(rLignes, rLigne)
	}
	return rLignes
}

func (p RenderPartie) NbJouéPlein() int {
	joué, _ := p.tj.ComptePlein()
	return joué
}

func (p RenderPartie) NbPlein() int {
	_, total := p.tj.ComptePlein()
	return total
}
