package solver

import (
	"testing"
	TJ "vintz.fr/nonogram/tabjeu"

)

func TestTabJeuCreate(t *testing.T) {
	taille := 10
	tj := TJ.NewTabJeu(taille, 50, 0)
	if tj == nil {
		t.Errorf("Wesh error")
	}

	var nb_plein int
	if len(tj) != taille {
		t.Errorf("Nb lignes = %d, attendu %d", len(tj), taille)
	}
	for _, ligne := range tj {
		if len(ligne) != taille {
			t.Errorf("Nb colonnes = %d, attendu %d", len(tj), taille)
		}
		for _, cell := range ligne {
			if cell.EstPlein() {
				nb_plein++
			}
		}
	}
}

func TestMain(t *testing.T) {
	tj := TJ.NewTabJeu(15, 40, 1003)
	tj.AfficheAvecComptes()
	prob := tj.MakeProbleme()
	Bench(prob, false)
}
