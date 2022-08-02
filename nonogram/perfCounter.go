package nonogram

import (
	"fmt"
	"sync"
	"time"
)

type PerfCounter struct {
	lock  sync.Mutex
	count int64
	stop  chan struct{}
}

func NewPerfCounter(interval time.Duration) *PerfCounter {
	pc := new(PerfCounter)
	pc.stop = make(chan struct{})
	pc.start(interval)
	return pc
}

func (pc *PerfCounter) Stop() {
	pc.stop <- struct{}{}
}

func (pc *PerfCounter) Inc() {
	pc.lock.Lock()
	defer pc.lock.Unlock()
	pc.count++
}

func (pc *PerfCounter) get() int64 {
	pc.lock.Lock()
	defer pc.lock.Unlock()
	return pc.count
}

func (pc *PerfCounter) start(interval time.Duration) {
	go func() {
		startTime := time.Now()
		lastTime := startTime
		lastCount := pc.get()
	loop:
		for {
			tick := time.After(interval)
			select {
			case now := <-tick:
				count := pc.get()
				rate := (count - lastCount) * 1e6 / now.Sub(lastTime).Nanoseconds()
				lastTime, lastCount = now, count
				fmt.Printf("Vitesse %d k/s      \r", rate)
			case <-pc.stop:
				count := pc.get()
				avg := count * 1e6 / time.Since(startTime).Nanoseconds()
				fmt.Printf("Vitesse moyenne %d k/s Total %d\n", avg, count)
				break loop
			}
		}
	}()
}
