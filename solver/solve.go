package solver

import (
	"fmt"
	TJ "vintz.fr/nonogram/tabjeu"
)

// lineList est une liste de lignes
// contrairement à TabJeu la largeur n'est pas forcément égale à la hauteur
type lineList []TJ.LigneJeu

// allPossibleLines est un slice de la même longueur que TabJeu
// chaque élement est l'ensemble des valeurs possibles pour la ligne correspondante
type lineListSet []lineList

// structure intermédiaire pour construire l'ensemble des lignes possibles
type indexedLineSet struct {
	num   int
	lines lineList
}

// cellVide et cellPlein sont les éléments de base pour construire les sequences possibles
// on aurait aussi bien pu utiliser des booleéns à la place des pointeurs
var cellVide = &TJ.Cellule{Base: TJ.Vide, Joué: TJ.Blanc}
var cellPlein = &TJ.Cellule{Base: TJ.Plein, Joué: TJ.Blanc}

type allPossibleBlocs struct {
	rows lineListSet
	cols lineListSet
}

type SolverFunc func(TJ.Probleme) chan *TJ.TabJeu


func (l lineList) String() string {
	var str string = ""
	for _, line := range l {
		str = str + fmt.Sprintf("%v\n", line)
	}
	return str
}

func (all lineListSet) String() string {
	var str string
	for _, pl := range all {
		str += pl.String()
	}
	return str
}

func (pl indexedLineSet) String() string {
	var str string = fmt.Sprintf("Sequences possibles pour la ligne %d\n", pl.num)
	str += pl.lines.String()
	return str
}


