package solver

import (
	"fmt"
	_ "fmt"
	"sync"
)

type TaskList []func()

// WorkerPool est un pool de workers qui executent des taches func()
// les taches sont des coroutines sans parametres ni valeurs de retour
type WorkerPool struct {
	ch chan func()     // queue d'entree du worker pool
	wg *sync.WaitGroup // attend la fin de la recursion
}

// NewWorkerPool initialise et démarre un pool de nbWorkers workers
func NewWorkerPool(nbWorkers int) (wp *WorkerPool) {

	wp = new(WorkerPool)
	wp.wg = new(sync.WaitGroup)
	wp.ch = make(chan func())

	// lance les workers
	for i := 0; i < nbWorkers; i++ {
		go func(id int) {
			fmt.Printf("Launched worker #%d\n", id)
			for tache := range wp.ch {
			//	fmt.Printf("Worker #%d  : start task\n", id)
				tache()
			//	fmt.Printf("Worker #%d  : end task\n", id)
				wp.wg.Done()
			}
			fmt.Printf("End of worker #%d\n", id)
		}(i)
	}
	return wp
}

// TryExec execute la tâche avec un worker disponible
// si aucun worker n'est disponible, on renvoie false
func (wp *WorkerPool) TryExec(task func()) (accepted bool) {
	wp.wg.Add(1)
	select {
	case wp.ch <- task: // execution avec un worker disponible
		accepted = true
		// wg.Done() appelé par le worker plus tard
	default:
		accepted = false
		wp.wg.Done()
	}
	return accepted
}

func (wp *WorkerPool) Exec(task func()) {
	wp.wg.Add(1)
	wp.ch <- task
}

// Wait() bloque jusqu'au traitement de la dernière tâche
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}
