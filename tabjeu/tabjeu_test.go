package tabjeu

import (
	"testing"
)

func TestTabJeuCreate(t *testing.T) {
	taille := 10
	tj := NewTabJeu(taille, 50, 0)
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
			if cell == Plein {
				nb_plein++
			}
		}
	}
}