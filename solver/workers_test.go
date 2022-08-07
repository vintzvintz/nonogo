package solver

import(
	"testing"
	"math/rand"
	"runtime"
	"time"
)


func TestExec(t *testing.T) {
	rand.Seed(time.Hour.Nanoseconds())
	nbWorkers := []int{ 1, 2, runtime.NumCPU(), 100 }
	// attention, deadlock avec 0 workers
	for _,n := range nbWorkers {
		testExec(t, n, 1000, 420 )
	}
}


func TestTryExec(t *testing.T) {
	rand.Seed(time.Hour.Nanoseconds())
	nbWorkers := []int{ 0, 1, 2, runtime.NumCPU(), 100 }
	// attention, execution synchrone avec 0 workers 
	for _,n := range nbWorkers {
		testTryExec(t, n, 1000, 421 )
	}
}


func testTryExec(t *testing.T, nbWorkers, nbTasks, sizeTask int) {

	valuesC := make(chan int)      // reçoit les valeurs produites par les tâches
	//testOkC := make(chan bool)    // pour attendre la fin du test avant de renvoyer le resultat

	wp := NewWorkerPool( nbWorkers )
	tl, wantSum := prepareTasks(t, nbTasks, sizeTask, valuesC )

	// lance la vérification des résultats (goroutine)
	okC := checkInts(t, valuesC, wantSum, sizeTask*nbTasks)

	// envoie des listes de taches au workerpool
	var nbFail int
	for _, task := range tl {
		accepted := wp.TryExec(task)    // non bloquant
		if( !accepted ) {
			// en cas d'echec de TryExec, on execute la tache de façon synchrone
			nbFail++
			task()
		}
	}

	t.Logf("TryExec() avec %d workers / %d taches. (%d refus du worker pool)", nbWorkers, nbTasks, nbFail)
	wp.Wait()
	close(valuesC)   // provoque la fin de checkInts()

	// attend le signal de fin de checkInts()
	<-okC 
}


func testExec(t *testing.T, nbWorkers, nbTasks, sizeTask int) {

	valuesC := make(chan int)      // reçoit les valeurs produites par les tâches
	//testOkC := make(chan bool)    // pour attendre la fin du test avant de renvoyer le resultat

	wp := NewWorkerPool( nbWorkers )
	tl, wantSum := prepareTasks(t, nbTasks, sizeTask, valuesC )

	// lance la vérification des résultats (goroutine)
	okC := checkInts(t, valuesC, wantSum, sizeTask*nbTasks)

	// envoie des listes de taches au workerpool
	for _, task := range tl {
		wp.Exec(task)    // bloque tant qu'il n'y a pas de worker dispo
	}
	t.Logf("TryExec() avec %d workers / %d taches", nbWorkers, nbTasks)
	wp.Wait()
	close(valuesC)   // provoque la fin de checkInts()

	// attend le signal de fin de checkInts()
	<-okC 
}


func checkInts ( t *testing.T, in chan int, wantSum, wantNb int ) (okC chan bool) {
	okC = make(chan bool)
	var gotSum, gotNb int
	go func() {
		for v := range in {
			gotSum += v
			gotNb += 1
		}
		ok := (gotSum==wantSum) && (gotNb==wantNb)
		if !ok {
			t.Errorf("WorkerPool ne renvoie pas le résulat attendu. (gotSum=%d/gotNb=%d, wantSum=%d/wantNb=%d)\n",gotSum, gotNb, wantSum, wantNb)
		}
		okC <- ok
		close(okC)
	}()
	return okC
}

func prepareTasks(t *testing.T, nbTask, sizeTask int, valuesC chan int) (tl []func(), attendu int) {
	
	// prepare une liste de taches
	tl = make ( []func(), nbTask )	
	for i := range tl {
		task, somme := makeTestTask( t, sizeTask, valuesC)
		tl[i] = task
		attendu += somme
	}
	return tl, attendu
}


// makeTasks crée une liste de nb tâches 
// chaque tâche renvoie des entiers aléatoires dans outC lorsqu'elle est executée
// somme est la somme des entiers qui seront ecrits par les tâches
func makeTestTask( t *testing.T, nb int, outC chan int ) (task func(), somme int)  {

	// prepare une liste de nombres
	vals := make( []int, nb)
	for i := range vals {
		vals[i] = rand.Intn(3000)
		somme += vals[i]
	}

	task = func() {
		for _, v := range vals {
			//if rand.Intn(10000)!=0 {   // pour déclencher des erreurs aleatoires
				outC <- v 
			//}
		}
	}

	return task, somme
}
