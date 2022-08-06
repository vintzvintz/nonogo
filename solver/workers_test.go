package solver

import(
	"testing"
	"math/rand"
	"runtime"
	"time"
)


// makeTasks crée une liste de nb tâches 
// chaque tâche renvoie un entier aléatoire lorsqu'elle est executée
// somme est la somme des entiers qui seront ecrits par les tâches
func makeTestTasks(t *testing.T, nb int, outC chan int ) (tl TaskList, somme int)  {

	tl = make( TaskList, nb)
	//var tmpVal int
	dbg := make( []int, nb)   // affiche les valeurs attendues
	for i:=0;i<nb;i++ {

		val := rand.Intn(1000)
		task := func() {
			outC <- val
		}
		tl[i] = task
		somme += val
		dbg[i] = val
	}
	return tl, somme
}

func testWorkers(t *testing.T, nbWorkers, nbListes, nbTasks int) {

	valuesC := make(chan int)      // reçoit les valeurs produites par les tâches
	testOkC := make(chan bool)    // pour attendre la fin du test avant de renvoyer le resultat

	var attendu int

	wp := NewWorkQueue( nbWorkers )

	// envoie des listes de taches au workerpool
	for nList:=0; nList<nbListes; nList++ {
		tl, somme := makeTestTasks(t, nbTasks, valuesC)
		wp.AddTasks(tl)
		attendu += somme

	}

	// additionne les valeurs reçues et compare avec le résutat attendu 
	go func() {
		var resultat int
		for v := range valuesC {
			//t.Log( "Reçoit", v)
			resultat += v
		}
		ok := resultat==attendu
		if !ok {
			t.Errorf("WorkerPool ne renvoie pas le résulat attendu. (got %d, want %d)\n",resultat, attendu)
		}
		testOkC <- ok
	}()

	wp.Wait()
	close(valuesC)

	// attend la fin du test pour terminer
	<-testOkC
}


func TestWorkerPool(t *testing.T) {
	rand.Seed(time.Hour.Nanoseconds())
	nbWorkers := []int{1, 2, runtime.NumCPU(), 100 }
	//nbWorkers = []int{1}
	for _,n := range nbWorkers {
		t.Logf("WorkerPool test %d workers", n)
		testWorkers(t, n, 500, 500)
	}
}