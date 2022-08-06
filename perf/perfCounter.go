package perf

import (
	"fmt"
	"sync/atomic"
	"time"
)

type PerfCounter struct {
	count int64
	stop  chan struct{}
	out   chan string
}

// NewPerfCounter cree et démarre un compteur d'evenements
func NewPerfCounter(interval time.Duration, out chan string) *PerfCounter {
	pc := new(PerfCounter)
	pc.stop = make(chan struct{})
	pc.out = out
	pc.start(interval)
	return pc
}

// Stop() arrete le compteur
func (pc *PerfCounter) Stop() {
	pc.stop <- struct{}{}
}

// Int() incrémente de compteur de n unités. Int() est thread-safe
func (pc *PerfCounter) Inc(n int64) {
	atomic.AddInt64( &pc.count, n)

}

// Get() renvoie la valeur courante du compteur
func (pc *PerfCounter) Get() int64 {
	return pc.count
}

func (pc *PerfCounter) sendMsg(msg string) {
	if pc.out != nil {
		pc.out <- msg
	}
}

func (pc *PerfCounter) start(interval time.Duration) {
	go func() {
		startTime := time.Now()
		lastTime := startTime
		lastCount := pc.Get()
	loop:
		for {
			tick := time.After(interval)
			select {
			case now := <-tick:
				count := pc.Get()
				rate := (count - lastCount) * 1e6 / now.Sub(lastTime).Nanoseconds()
				lastTime, lastCount = now, count
				pc.sendMsg(fmt.Sprintf("Vitesse %d k/s      \r", rate))

			case <-pc.stop:
				count := pc.Get()
				avg := count * 1e6 / time.Since(startTime).Nanoseconds()
				pc.sendMsg(fmt.Sprintf("Vitesse moyenne %d k/s Total %d\n", avg, count))
				break loop
			}
		}
		if pc.out != nil {
			close(pc.out)
		}
	}()
}
