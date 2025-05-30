package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

func main() {
	conf := getConfig()

	workerAmount := 0
	for workerType := range len(conf.WorkersAmount) {
		workerAmount += conf.WorkersAmount[workerType]
	}

	model := initialModel(conf, workerAmount, 15)
	prog := tea.NewProgram(&model, tea.WithAltScreen(), tea.WithoutSignalHandler()) // implement ^C as interrupt or other way to close the program
	model.storage = InitStorage(conf, prog)
	model.workstations = InitWorkstations(conf.QuarryWorkplaces)

	i := 0
	for workerType := range len(conf.WorkersAmount) {
		for workerNo := range conf.WorkersAmount[workerType] {
			workerId, err := makeWorkerId(workerType, workerNo)

			if err != nil {
				log.Fatal(err)
			}
			model.workerPos[workerId] = 0
			model.workstationQueue = append(model.workstationQueue, false)
			model.storageQueue = append(model.storageQueue, false)
			model.workerIds[workerId] = i
			model.workers[i] = initWorker(workerId, workerType, model.storage, model.workstations, prog, conf)
			i++
		}
	}

	if _, err := prog.Run(); err != nil {
		log.Fatal(err)
	}

	// Clean up
	for w := range model.workers {
		model.workers[w].done <- true
	}
}
