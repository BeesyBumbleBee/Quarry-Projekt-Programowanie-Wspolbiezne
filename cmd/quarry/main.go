package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"sync"
	"time"
)

/*
\u263B - postać
\U0001F3ED - fabryka
\U0001F528 - młotek
\u26CF - kilof
\u2692 - narzędzia
*/

type tickMsg time.Time

type model struct {
	workers                  []Worker
	workersPositions         []int
	config                   SimulationConfig
	workerPositionController *WorkerPositionController
}

func initialModel(workerAmount int, config SimulationConfig) model {
	return model{
		workers:                  make([]Worker, workerAmount),
		workersPositions:         make([]int, workerAmount),
		config:                   config,
		workerPositionController: nil,
	}
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	done := true
	go m.workerPositionController.GetEnqueuedMessages(&done)
	workerMsg := WorkerMessage{}
	for done {
		workerMsg = <-m.workerPositionController.DisplayChannel
		if workerMsg.WorkerId == -1 {
			break
		}
		m.workersPositions[workerMsg.WorkerId] = workerMsg.Position
	}
	return m, tickCmd()
}

func (m model) View() string {
	s := ""
	for category := range 3 {
		if category == 1 {
			s += "\n\U0001F3ED"
		} else {
			s += "\n  "
		}
		road := []rune("          ")
		for pos := range 10 {
			id := 0
			for _ = range m.config.WorkersAmount[category] {
				if m.workersPositions[id] == pos {
					road[pos] = rune('\u263B')
				}
				id++
			}
		}
		s += string(road)
		s += "\u2692"
	}

	return s
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*200, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func main() {
	conf := getConfig()
	wg := sync.WaitGroup{}

	workerAmount := 0
	for workerType := range len(conf.WorkersAmount) {
		workerAmount += conf.WorkersAmount[workerType]
	}

	model := initialModel(workerAmount, conf)

	conf.printConfig()

	model.workerPositionController = NewWorkerPositionController(len(model.workers))
	if model.workerPositionController == nil {
		panic("cannot initialize workerPositionController")
	}

	id := 0
	for workerType := range len(conf.WorkersAmount) {
		for _ = range conf.WorkersAmount[workerType] {
			model.workers[id] = initWorker(
				id,
				workerType,
				conf.TimeToTravelEmpty,
				conf.TimeToTravelFull,
				model.workerPositionController.WorkersChannel)
			id++
		}
	}
	go model.workerPositionController.ManageWorkersPositions()
	for i := range len(model.workers) {
		go model.workers[i].Work()
	}

	prog := tea.NewProgram(model)
	if _, err := prog.Run(); err != nil {
		log.Fatal(err)
	}
	wg.Wait()
}
