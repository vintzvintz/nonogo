package solver

import (
	"sync"
)


type TaskList []func()

// WorkerPool est un pool de workers qui executent des taches func() 
// les taches sont des coroutines sans parametres ni valeurs de retour
type WorkerPool struct {
	ch chan func()     // queue d'entree du worker pool
	wg *sync.WaitGroup // attend la fin de la recursion
}


// NewWorkQueue initialise et démarre un pool de nbWorkers workers 
func NewWorkQueue( nbWorkers int) (wp *WorkerPool) {

	if nbWorkers <= 0 {
		return nil
	}

	wp = new(WorkerPool)
	wp.wg = new(sync.WaitGroup)
	// recursion monothread classique si workQueue.ch = nil

	wp.ch = make(chan func())

	// lance les workers
	for i := 0; i < nbWorkers; i++ {
		go func(id int) {
			for tache := range wp.ch {
				tache()
				wp.wg.Done()
			}
		}(i)
	}
	return wp
}

// AddTasks ajoute des taches à traiter par les workers
func (wp* WorkerPool) AddTasks( tasks TaskList ) {
	// envoi des avec une goroutine
	// les taches sont "consommées" à travers le chan par les workers
	wp.wg.Add(len(tasks))
	go func() {
		for _, t := range tasks {
			wp.ch <- t
		}
	}()
}


// Wait() bloque jusqu'au traitement de la dernière tâche
func (wp* WorkerPool) Wait() {
	wp.wg.Wait()
}