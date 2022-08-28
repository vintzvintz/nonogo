package tabjeu

import (
	"testing"
)

type tjDef struct {
	taille int
	nbPlein int
}

func TestTabJeuCreate(t *testing.T) {
	tbl := []tjDef{
		{ taille: 1,   nbPlein: 0},
		{ taille: 1,   nbPlein: 1},
		{ taille: 10,  nbPlein: 0},
		{ taille: 10,  nbPlein: 10},
		{ taille: 10,  nbPlein: 100},
		{ taille: 15,  nbPlein: 130},
		{ taille: 100, nbPlein: 9000},
	}
	for _, def := range tbl {
		testTj(t, def.taille, def.nbPlein)
	}
}

func testTj(t *testing.T, taille int, nbPlein int ) {
	tj := NewTabJeu(taille, nbPlein, 0)

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
	if gotPlein != nbPlein {
		t.Errorf("Nb cellules remplies = %d, attendu %d", gotPlein, nbPlein )
	}
}
