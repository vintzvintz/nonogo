package partie

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"text/template"
	"time"
	"sync"

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


const (
	JOUE_VIDE = iota
	JOUE_PLEIN
)

const DEFAULT_SIZE = 10
const DEFAULT_SEED = 1005
const DEFAULT_RATIO = 0.45

type Partie struct {
	seed int64
	tj tabjeu.TabJeu
	lock *sync.Mutex
}


func NewPartie(size int, nbPlein int, seed int64) Partie {

	tj_brut := tabjeu.NewTabJeu(size, nbPlein, seed)
	tj := desambigueTabJeu(tj_brut, -1)

	return Partie{
		seed:seed,
		tj:tj,
		lock: new(sync.Mutex) }
}

func NewPartieDefault() Partie {
	var nbPlein int = DEFAULT_RATIO * DEFAULT_SIZE * DEFAULT_SIZE
	return NewPartie(DEFAULT_SIZE, nbPlein, DEFAULT_SEED)
}


func (p *Partie) Clique(action int, ligne, colonne int) {
		p.lock.Lock()
	defer p.lock.Unlock()

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

	// Change l'état de la cellule selon l'action demandée et selon l'état précédent
	switch action {
	case JOUE_PLEIN : c.TogglePlein()
	case JOUE_VIDE : c.ToggleVide()
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
	fmt.Printf("%d solutions et %d cellules ambigües. %d workers en %v\n", nbSol, diff.Count(), nbWorkers, duree)

	valide = tj.Copy()
	valide.ReveleVides(diff)

	return valide
}
