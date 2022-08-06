package solver

import (
	"testing"
	TJ "vintz.fr/nonogram/tabjeu"

)


func ttttest(t *testing.T) {
	tj := TJ.NewTabJeu(15, 45, 1003)
	tj.AfficheAvecComptes()
	prob := tj.MakeProbleme()
	Bench(prob, false)
}
