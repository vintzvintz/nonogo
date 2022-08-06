package perf

import (
	"sync"
	"testing"
	"time"
	"os"
)

func TestPerfCounter(t *testing.T) {

	const nbEvt int64 = 1000000
	const nbGoroutines = 200
	wgCount := new(sync.WaitGroup)
	pc := NewPerfCounter(time.Second/10, os.Stdout)


	// lance le comptage
	wgCount.Add(nbGoroutines)
	for r := 0; r < nbGoroutines; r++ {
		go func() {
			var i int64
			for i = 0; i < nbEvt; i++ {
				pc.Inc(1)
			}
			wgCount.Done()
		}()
	}

	wgCount.Wait()
	pc.Stop()

	got := pc.Get()
	want := nbEvt * nbGoroutines

	if got != want {
		t.Errorf("Nombre d'evements comptés: %d  (attendu %d)", got, want)
	}

}
