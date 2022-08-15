package tabjeu

import (
	"testing"
)

type tjDef struct {
	taille int
	ratio float32
}

func TestTabJeuCreate(t *testing.T) {
	tbl := []tjDef{
		{ taille: 1,   ratio: 0.00},
		{ taille: 1,   ratio: 1.00},
		{ taille: 10,  ratio: 0.00},
		{ taille: 10,  ratio: 0.10},
		{ taille: 10,  ratio: 1.00},
		{ taille: 15,  ratio: 0.50},
		{ taille: 100, ratio: 0.90},
	}
	for _, def := range tbl {
		testTj(t, def.taille, def.ratio)
	}
}

func testTj(t *testing.T, taille int, ratio float32) {
	tj := NewTabJeu(taille, ratio, 0)

	if len(tj) != taille {
		t.Errorf("Nb lignes = %d, attendu %d", len(tj), taille)
	}
	var gotPlein int
	for _, ligne := range tj {
		if len(ligne) != taille {
			t.Errorf("Nb colonnes = %d, attendu %d", len(tj), taille)
		}
		for _, cell := range ligne {
			if cell.EstPlein() {
				gotPlein++
			}
		}
	}
	wantPlein := int(float32(taille*taille)*ratio)
	if gotPlein != wantPlein {
		t.Errorf("Nb cellules remplies = %d, attendu %d", gotPlein, wantPlein )
	}
}
