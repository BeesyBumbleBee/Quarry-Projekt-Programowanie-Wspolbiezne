package main

import (
	"math/rand/v2"
	"time"
)

type Worker struct {
	Id                int
	Category          int
	Position          int // 0 - 100 gdzie 0 oznacza przy składowisku a 100 przy stanowisku pracy
	TimeToTravelEmpty int
	TimeToTravelFull  int
	MovingToStorage   bool
	PositionChannel   chan WorkerMessage
}

func initWorker(id int, category int, timeToTravelEmpty int, timeToTravelFull int, positionChannel chan WorkerMessage) Worker {
	return Worker{
		Id:                id,
		Category:          category,
		Position:          0,
		TimeToTravelEmpty: timeToTravelEmpty,
		TimeToTravelFull:  timeToTravelFull,
		MovingToStorage:   false,
		PositionChannel:   positionChannel,
	}
}

func (w Worker) Work() {
	for {
		if w.MovingToStorage {
			time.Sleep(time.Duration(rand.IntN(w.TimeToTravelFull)+w.TimeToTravelFull) * time.Millisecond)
			w.Position -= 1
			if w.Position < 0 {
				w.MovingToStorage = false
			}
		} else {
			time.Sleep(time.Duration(rand.IntN(w.TimeToTravelEmpty)+w.TimeToTravelEmpty) * time.Millisecond)
			w.Position += 1
			if w.Position > 10 {
				w.MovingToStorage = true
			}
		}
		w.PositionChannel <- WorkerMessage{WorkerId: w.Id, Position: w.Position}
	}
}
