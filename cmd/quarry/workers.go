package main

import (
	"errors"
	tea "github.com/charmbracelet/bubbletea"
	"math/rand"
	"strconv"
	"time"
)

type Worker struct {
	workerType     int
	id             string
	position       int
	goingToStorage bool
	done           chan bool
	rand           *rand.Rand
	storage        *Storage
	workstations   *Workstations
	program        *tea.Program
	cfg            SimulationConfig
}

type workerFinishedWork struct {
	workerId string
}
type workerAtStorage struct {
	workerId string
}

type workerWorking struct {
	workerId string
}

type workerAtWork struct {
	workerId string
}
type workerMoveMsg struct {
	workerId   string
	workerType int
	position   int
}

func initWorker(id string, workerType int, storage *Storage, workstations *Workstations, program *tea.Program, cfg SimulationConfig) Worker {
	return Worker{
		id:             id,
		workerType:     workerType,
		position:       48,
		goingToStorage: true,
		done:           make(chan bool, 1),
		storage:        storage,
		workstations:   workstations,
		program:        program,
		cfg:            cfg,
		rand:           rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func makeWorkerId(workerType int, workerNo int) (string, error) {
	types := make([]rune, 3)
	types[0] = 'A'
	types[1] = 'B'
	types[2] = 'C'

	if workerType >= 3 {
		err := errors.New("worker type out of range")
		return "", err
	}

	return string(types[workerType]) + strconv.Itoa(workerNo+1), nil
}

func (w *Worker) Work() {
	for {
		select {
		case <-w.done:
			return
		default:
		}

		if w.position == 49 && w.goingToStorage {
			w.goingToStorage = false
			w.tryWork()
		} else if w.position == 0 && !w.goingToStorage {
			w.goingToStorage = true
			if !w.tryPlace() {
				return
			}
		} else {
			if w.goingToStorage {
				w.move(1)
			} else {
				w.move(-1)
			}
		}

	}
}

func (w *Worker) move(direction int) {
	if direction >= 1 {
		time.Sleep(time.Duration(w.cfg.TimeToTravelEmpty[0]+
			w.rand.Intn(w.cfg.TimeToTravelEmpty[1]-w.cfg.TimeToTravelEmpty[0])) * time.Millisecond)
	} else {
		time.Sleep(time.Duration(w.cfg.TimeToTravelFull[0]+
			w.rand.Intn(w.cfg.TimeToTravelFull[1]-w.cfg.TimeToTravelFull[0])) * time.Millisecond)
	}
	w.position += direction
	w.program.Send(workerMoveMsg{w.id, w.workerType, w.position})
}

func (w *Worker) delayPlace() {
	time.Sleep(time.Duration(w.cfg.TimeToPlaceStone[0]+
		w.rand.Intn(w.cfg.TimeToPlaceStone[1]-w.cfg.TimeToPlaceStone[0])) * time.Millisecond)
}

func (w *Worker) tryPlace() bool {
	w.program.Send(workerAtStorage{w.id})
	w.storage.hasPalletCleared.L.Lock()
	for !w.storage.Place(w.workerType, w.id, w.delayPlace) {
		select { // check if program wasn't termintated while waiting
		case <-w.done:
			return false
		default:
		}
		w.storage.hasPalletCleared.Wait()
	}
	w.storage.hasPalletCleared.L.Unlock()
	return true
}

func (w *Worker) tryWork() {
	w.program.Send(workerAtWork{w.id})
	w.workstations.EnterWork()
	w.program.Send(workerWorking{w.id})
	time.Sleep(time.Duration(w.cfg.StonesExtractionTime[w.workerType][0]+
		w.rand.Intn(w.cfg.StonesExtractionTime[w.workerType][1]-w.cfg.StonesExtractionTime[w.workerType][0])) * time.Millisecond)
	w.workstations.ExitWork()
	w.program.Send(workerFinishedWork{w.id})
}
