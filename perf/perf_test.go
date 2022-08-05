package perf

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestPerfCounter(t *testing.T) {

	const nbEvt int64 = 100
	const nbGoroutines = 200000
	wgCount := new(sync.WaitGroup)
	wgPrint := new(sync.WaitGroup)

	// affiche l'avancement
	output := make(chan string)
	wgPrint.Add(1)
	go func() {
		for s := range output {
			fmt.Print(s)
		}
		wgPrint.Done()
	}()

	pc := NewPerfCounter(time.Second/10, output)

	// lance le comptage
	wgCount.Add(nbGoroutines)
	for r := 0; r < nbGoroutines; r++ {
		go func() {
			var i int64
			for i = 0; i < nbEvt; i++ {
				pc.Inc()
			}
			wgCount.Done()
		}()
	}

	wgCount.Wait()
	pc.BlockingStop()
	wgPrint.Wait()

	got := pc.Get()
	want := nbEvt * nbGoroutines

	if got != want {
		t.Errorf("Nombre d'evements comptés: %d  (attendu %d)", got, want)
	}

}
