package perf

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

type PerfCounter struct {
	count  int64
	stop   chan struct{}
	out    io.Writer
	wg     *sync.WaitGroup
}

// NewPerfCounter cree et démarre un compteur d'evenements
func NewPerfCounter(interval time.Duration, out io.Writer) *PerfCounter {
	pc := new(PerfCounter)
	pc.stop = make(chan struct{})
	pc.out = out
	pc.wg = new(sync.WaitGroup)

	// le WaitGroup sert à bloquer sur Stop() pour garantir que le dernier message sera affiché
	pc.wg.Add(1)
	go pc.run(interval, out)
	return pc
}

// Stop() arrete le compteur
func (pc *PerfCounter) Stop() {
	pc.stop <- struct{}{}
	pc.wg.Wait() // attend l'affichage des derniers messages avant de quitter
}

// Int() incrémente de compteur de n unités. Int() est thread-safe
func (pc *PerfCounter) Inc(n int64) {
	atomic.AddInt64(&pc.count, n)
}

// Int() incrémente de compteur de n unités. Int() est thread-safe
func (pc *PerfCounter) Inc1() {
	atomic.AddInt64(&pc.count, 1)
}

// Get() renvoie la valeur courante du compteur
func (pc *PerfCounter) Get() int64 {
	return atomic.LoadInt64(&pc.count)
}

func (pc *PerfCounter) run(interval time.Duration, out io.Writer) {

	// fonction utilitaire pour afficher la vitesse et la progression
	sendMsg := func(msg string) {
		if out != nil {
			pc.out.Write([]byte(msg))
		}
	}

	startTime := time.Now()
	lastTime := startTime
	lastCount := pc.Get() //should be 0....
	tick := time.NewTicker(interval)

loop:
	for {
		select {
		case now := <-tick.C:
			count := pc.Get()
			rate := (count - lastCount) * 1e6 / now.Sub(lastTime).Nanoseconds()
			sendMsg(fmt.Sprintf("Vitesse %d k/s   \r", rate))
			lastTime, lastCount = now, count
		case <-pc.stop:
			count := pc.Get()
			avg := count * 1e6 / time.Since(startTime).Nanoseconds()
			sendMsg(fmt.Sprintf("Vitesse moyenne %d k/s Total %d\n", avg, count))
			break loop
		}
	}
	tick.Stop()
	pc.wg.Done()
}
