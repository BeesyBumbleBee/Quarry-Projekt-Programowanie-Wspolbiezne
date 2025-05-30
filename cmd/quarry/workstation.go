package main

import "sync"

type Workstations struct {
	workersLimit     int
	workersWorking   int
	workstationsFull *sync.Cond
}

func InitWorkstations(workersLimit int) *Workstations {
	return &Workstations{
		workersLimit:     workersLimit,
		workersWorking:   0,
		workstationsFull: sync.NewCond(&sync.Mutex{}),
	}
}

func (w *Workstations) EnterWork() {
	w.workstationsFull.L.Lock()
	for w.workersWorking == w.workersLimit {
		w.workstationsFull.Wait()
	}
	w.workstationsFull.L.Unlock()
	w.workersWorking++
}

func (w *Workstations) ExitWork() {
	w.workersWorking--
	w.workstationsFull.Signal()
}
