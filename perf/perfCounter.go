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

func NewPerfCounter(interval time.Duration, out chan string) *PerfCounter {
	pc := new(PerfCounter)
	pc.stop = make(chan struct{})
	pc.start(interval)
	if out != nil {
		pc.out = out
	}
	return pc
}

func (pc *PerfCounter) Stop() {
	fmt.Println("Stop requested")
	pc.stop <- struct{}{}
	if pc.out != nil{
		close(pc.out)
	}
}


func (pc *PerfCounter) BlockingStop() {
	fmt.Println("Blocking stop requested")
	pc.stop <- struct{}{}

	// attend la fermeture complete du compteur
	<-pc.stop

	if pc.out != nil {
		close(pc.out)
	}
}


func (pc *PerfCounter) Inc() {
	atomic.AddInt64( &pc.count, 1)

}

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
				fmt.Println("Stopped !")
				count := pc.Get()
				avg := count * 1e6 / time.Since(startTime).Nanoseconds()
				pc.sendMsg(fmt.Sprintf("Vitesse moyenne %d k/s Total %d\n", avg, count))
				break loop
			}
		}
		// ecrit dans le chan pour signaler la fin du compteur
		close(pc.stop)
	}()
}
