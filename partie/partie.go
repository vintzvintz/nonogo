package partie

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"text/template"
	"time"

	"vintz.fr/nonogram/solver"
	"vintz.fr/nonogram/tabjeu"
)

type Etat int

const (
	PAS_COMMENCE Etat = iota
	EN_COURS
	GAGNE
	PERDU
)


type Partie struct {
	seed int64
	tj tabjeu.TabJeu
}


func NewPartie(size int, nbPlein int, seed int64) Partie {

	tj_brut := tabjeu.NewTabJeu(size, nbPlein, seed)
	tj := desambigueTabJeu(tj_brut, -1)

	return Partie{seed, tj}
}

func NewPartieDefault() Partie {
	var nbPlein int = tabjeu.DEFAULT_RATIO * tabjeu.DEFAULT_SIZE * tabjeu.DEFAULT_SIZE
	return NewPartie(tabjeu.DEFAULT_SIZE, nbPlein, tabjeu.DEFAULT_SEED)
}


func (p *Partie) Clique(ligne, colonne int) {

	// ignore les coordonnées invalides
	if( ligne<0 || ligne >= len(p.tj) ){
		log.Printf("Ignore Clique() sur ligne %d invalide\n", ligne)
		return
	}
	if( colonne<0 || colonne >= len(p.tj) ){
		log.Printf("Ignore Clique() sur colonne %d invalide\n", colonne)
		return
	}

	c := &p.tj[ligne][colonne]

	// aucune action sur les cellules révélées
	if c.EstRévélé() {
		log.Printf("Ignore Clique() sur cellule (%d,%d) révélée\n", ligne, colonne)
		return
	}
	// cycle aucun -> plein -> vide
	if c.EstJouéPlein() {
		c.JoueVide()
	} else if c.EstJouéVide() {
		c.JoueAucun()
	} else {
		c.JouePlein()
	}
}


func (p Partie) Html() (*bytes.Buffer, error) {

	buf := bytes.NewBuffer(nil)

	tmpl, err := template.ParseFiles("templates/common.tmpl")
	if err != nil {
		return buf, err
	}

	err = tmpl.ExecuteTemplate(buf, "page", RenderPartie(p))
	if err != nil {
		return nil, err
	}

	return buf, nil
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
