package partie

import (
	"fmt"
	"runtime"
	"time"
	"bytes"
	"text/template"

	"vintz.fr/nonogram/solver"
	"vintz.fr/nonogram/tabjeu"
)


type Partie struct {
	Seed int
	TauxRemplissage float32
	Tj tabjeu.TabJeu
}


type RenderPartie struct {
	Taille int
	DimensionCellules int
	PremiereLigne tabjeu.BlocCountList
	Lignes []RenderLigne
}

type RenderLigne struct {
	Blocs tabjeu.BlocCount
	Cellules []RenderCellule
}

type RenderCellule struct {
	Col,Ligne int
	Img string
	Lien string
}


func newLevel(size int, ratio float32, seed int64) tabjeu.TabJeu {
	tj := tabjeu.NewTabJeu(size, ratio, seed)
	return desambigueTabJeu(tj, -1)
}

func NewDefault() Partie {
	return Partie{
		Seed:  tabjeu.DEFAULT_SEED,
		TauxRemplissage: tabjeu.DEFAULT_RATIO,
		Tj: newLevel(tabjeu.DEFAULT_SIZE, tabjeu.DEFAULT_RATIO, tabjeu.DEFAULT_SEED),
	}
}

func desambigueTabJeu(tj tabjeu.TabJeu, nbWorkers int) (valide tabjeu.TabJeu) {

	taille := len(tj)
	startTime := time.Now()
	if nbWorkers < 0 {
		nbWorkers = runtime.NumCPU()
	}

	var nbSol int
	var prev tabjeu.TabJeu
	var diff tabjeu.Diff = tabjeu.NewDiff(taille)
	for solution := range solver.SolveConcurrent(tj, nbWorkers, true) {
		if prev != nil {
			solution.Compare(prev, diff)
		}
		prev = solution
		nbSol++
	}

	duree := time.Since(startTime)
	fmt.Printf(" %d solutions et %d cellules ambigües. %d workers en %v\n", nbSol, diff.Count(), nbWorkers, duree)
	
	valide = tj.Copy()
	valide.ReveleVides(diff)

	return valide
}


func (p Partie) prepareRenderData() (r *RenderPartie) {

	blocsLignes := p.Tj.CompteBlocs( tabjeu.LIGNE )
	
	var rLignes []RenderLigne
	for l := range p.Tj {
		var rCells []RenderCellule

		for c, cell := range p.Tj[l] {
			var img string
			switch cell.Image() {
				case tabjeu.IMG_AUCUN: img = "vide.svg"
				case tabjeu.IMG_VIDE: img = "croix.svg"
				case tabjeu.IMG_PLEIN: img = "carré.svg"
			default: panic( "tabjeu.Image() renvoie une valeur inattendue")
			}
			rCell := RenderCellule{ 
				Ligne: l,
				Col: c,
				Img: img,
				Lien: "clic",
			}
			rCells = append(rCells, rCell)
		}

		rLigne := RenderLigne{ 
			Blocs: blocsLignes[l],
			Cellules: rCells,
		}
		rLignes = append(rLignes, rLigne)
	}

	return &RenderPartie{
		Taille:len(p.Tj),
		DimensionCellules: 6,
		PremiereLigne: p.Tj.CompteBlocs( tabjeu.COLONNE ),
		Lignes: rLignes,
		}
}

func (p Partie) Html() (*bytes.Buffer, error) {

	buf := bytes.NewBuffer(nil)

	tmpl, err := template.ParseFiles("templates/common.tmpl")
	if err!=nil {
		return buf, err
	}

	data  := p.prepareRenderData()

	err = tmpl.ExecuteTemplate(buf, "page", data)
	if err!=nil {
		return nil, err
	}

	return buf, nil
}
