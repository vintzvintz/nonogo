package nonogram

import (
	"fmt"
	"time"
)

type multiIndex struct {
	idx []int
	max []int
}

func NewCompteur(max []int) chan []int {

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
