package solver

import (
	"fmt"
	"sync"
	"time"

	TJ "vintz.fr/nonogram/tabjeu"
)

type multiIndex struct {
	idx []int
	max []int
}

func newCompteur(max []int) chan []int {

	taille := len(max)

	// initialise les index
	mi := multiIndex{}
	mi.idx = make([]int, taille)
	mi.max = make([]int, taille)

	// valeur maximale de chaque index elementaire
	copy(mi.max, max)

	// cree le channel pour produire les valeurs successives
	ch := make(chan []int)

	// calcule le nb total d'itérations pour calculer le % de progression
	var nbMax int64 = 1
	for _, m := range mi.max {
		nbMax = nbMax * int64(m)
	}
	fmt.Printf("%v soit %d combinaisons à tester\n", mi.max, nbMax)

	go func() {
		// affichage périodique de l'avancement
		var progress int64        // compteur d'avancement
		var nextNotifVal int64    // prochain point de notification
		const stepNotif = 1000000 // intervalle de notification

		lastNotifTime := time.Now().UnixNano()

		for {
			// renvoie une copie car mi.idx est modifié pendant le traitement
			figé := make([]int, taille)
			copy(figé, mi.idx)
			ch <- figé

			// Affiche la progression
			if progress >= nextNotifVal {
				percentProgress := 100 * float32(progress) / float32(nbMax)
				nextNotifVal += stepNotif

				//cadence
				now := time.Now().UnixNano()
				var tjPerSec int64
				deltaT := now - lastNotifTime
				if deltaT > 0 {
					tjPerSec = 1e9 * stepNotif / (now - lastNotifTime)
				}
				lastNotifTime = now
				fmt.Printf("Progression %f%% %d tests/sec\n", percentProgress, tjPerSec)
			}
			progress++

			// incremente le compteur et termine quand toutes les combinaisons sont épuisées
			overflow := mi.incrémente()
			if overflow {
				close(ch)
				break
			}
		}
	}()
	return ch
}

// incremente fait passer le multiIndex à la valeur suivante
// renvoie true quand il y a overflow
func (mi *multiIndex) incrémente() bool {
	var overflow bool
	for i := range mi.idx {
		mi.idx[i] += 1

		// termine si l'index courant ne génère par de retenue
		if mi.idx[i] < mi.max[i] {
			break
		}
		// sinon on remet l'index courant à 0 et on incrémente le suivant
		mi.idx[i] = 0

		// renvoie l'overflow sur le dernier index
		if i == (len(mi.idx) - 1) {
			overflow = true
		}
	}
	return overflow
}



func makeTache(taille int,
	apl *lineListSet,
	idx []int,
	testFunc func(tj *TJ.TabJeu) bool,
	out chan *TJ.TabJeu) func() {

	return func() {
		tj := make(TJ.TabJeu, taille)
		// remplit tj avec les lignes spécifiées par les index de cptVal
		for i := range idx {
			lignes := (*apl)[i]
			tj[i] = lignes[idx[i]]
		}
		// teste si tj est une solution valide
		if testFunc(&tj) {
			out <- &tj
		}
	}
}

// renvoie une closure qui teste si tj est cohérent avec les longueurs de blocs en colonne
func makeTestSolution(seqColSolution []TJ.BlocCount) func(tj *TJ.TabJeu) bool {
	return func(tj *TJ.TabJeu) bool {
		// compte les longueurs de blocs en colonne
		return tj.CompareBlocsColonnes(seqColSolution)
	}
}

// Trouvesolutions revoie toutes les solutions possibles du nonogramme
// à partir des longueurs de blocs en lignes et en colonnes
func SolveBourrin(prob TJ.Probleme) chan *TJ.TabJeu {
	allSeqs := buildAllSequences(prob.BlocsLignes)

	// closure qui teste si tj est une solution valide
	testFunc := makeTestSolution(prob.BlocsColonnes)

	// calcule le nb de sequences possibles pour chaque ligne du jeu
	// et cree un "multi-index" pour parcourir toutes les sequences possibles
	nbSeq := make([]int, len(allSeqs))
	for i := range allSeqs {
		nbSeq[i] = len(allSeqs[i])
	}
	compteur := newCompteur(nbSeq)

	// recoit la ou les solutions correctes
	out := make(chan *TJ.TabJeu)

	//reçoit les taches à traiter ( closures sans parametres ni valeurs de retour )
	jobsC := make(chan func())

	// cree les taches et envoie dans la queue
	go func() {
		for idx := range compteur {
			task := makeTache(prob.Taille, &allSeqs, idx, testFunc, out)
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
