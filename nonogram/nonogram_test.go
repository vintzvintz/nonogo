package nonogram

import (
	"testing"
	//"vintz.fr/gotest/nonogram"
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
			if cell.estPlein() {
				nb_plein++
			}
		}
	}
}
