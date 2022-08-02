package nonogram


import (
	"sync"
)

type WorkerPool struct {
	ch chan func()     // queue d'entree du worker pool
	wg *sync.WaitGroup // attend la fin de la recursion
}


// NewWorkQueue initialise un pool de workers 
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
func (wp* WorkerPool) AddTasks( tasks []func() ) {
	// envoi des tâches au pool de workers avec une goroutine
	wp.wg.Add(len(tasks))
	go func() {
		for _, t := range tasks {
			wp.ch <- t
		}
	}()

}


func (wp* WorkerPool) Wait() {
	wp.wg.Wait()
}