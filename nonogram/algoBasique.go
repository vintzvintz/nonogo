package nonogram

import (

	//"fmt"
	"sync"
)

func makeTache(taille int,
	apl *lineListSet,
	idx []int,
	//seqColonnes []seqCount,
	testFunc func(tj *TabJeu) bool,
	out chan *TabJeu) func() {

	return func() {
		tj := make(TabJeu, taille)
		// remplit tj avec les lignes spécifiées par les index de cptVal
		for i := range idx {
			lignes := (*apl)[i]
			tj[i] = lignes[idx[i]]
		}
		// teste si tj est une solution valide
		if testFunc(&tj) {
			//if tj.CompteBlocsCompare(seqColonnes) {
			out <- &tj
		}
	}
}

// renvoie une closure qui teste si tj est cohérent avec les longueurs de blocs en colonne
func makeTestSolution(seqColSolution []seqCount) func(tj *TabJeu) bool {
	return func(tj *TabJeu) bool {
		// compte les longueurs de blocs en colonne
		return tj.CompteBlocsCompare(seqColSolution)
	}
}

// Trouvesolutions revoie toutes les solutions possibles du nonogramme
// à partir des longueurs de blocs en lignes et en colonnes
func SolveBourrin(prob Probleme) chan *TabJeu {
	allSeqs := buildAllSequences(prob.taille, prob.seqLignes)

	// closure qui teste si tj est une solution valide
	testFunc := makeTestSolution(prob.seqColonnes)

	// calcule le nb de sequences possibles pour chaque ligne du jeu
	// et cree un "multi-index" pour parcourir toutes les sequences possibles
	nbSeq := make([]int, len(allSeqs))
	for i := range allSeqs {
		nbSeq[i] = len(allSeqs[i])
	}
	compteur := NewCompteur(nbSeq)

	// recoit la ou les solutions correctes
	out := make(chan *TabJeu)

	//reçoit les taches à traiter ( closures sans parametres ni valeurs de retour )
	jobsC := make(chan func())

	// cree les taches et envoie dans la queue
	go func() {
		for idx := range compteur {
			task := makeTache(prob.taille, &allSeqs, idx, testFunc, out)
			jobsC <- task
		}
		// ferme le channel lorsque la derniere tâche est traitée
		close(jobsC)
	}()

	// cree un pool de workers qui traitent les taches placées dans la queue
	wg := new(sync.WaitGroup)
	for workerId := 0; workerId < 8; workerId++ {
		wg.Add(1)
		go func(id int, queue chan func()) {
			var nb int
			for task := range queue { // termine lorsque la queue est fermée
				task()
				nb++
			}
			//fmt.Printf("Worker %d terminé après %d taches\n", id, nb)
			wg.Done()
		}(workerId, jobsC)
	}

	// ferme le channel de sortie quand tous les workers se terminent
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
