package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"os"
)

func main() {
	conf, err := startMainMenu()
	if err != nil {
		os.Exit(0)
	}
	err = startSimulation(conf)
	if err != nil {
		log.Fatal(err)
	}
}

func startSimulation(cfg SimulationConfig) error {
	var err error = nil

	workerAmount := 0
	for workerType := range len(cfg.WorkersAmount) {
		workerAmount += cfg.WorkersAmount[workerType]
	}

	model := initialModel(cfg, workerAmount, 40)
	prog := tea.NewProgram(&model, tea.WithAltScreen(), tea.WithoutSignalHandler()) // implement ^C as interrupt or other way to close the program
	model.storage = InitStorage(cfg, prog)
	model.workstations = InitWorkstations(cfg.QuarryWorkplaces)

	i := 0
	for workerType := range len(cfg.WorkersAmount) {
		for workerNo := range cfg.WorkersAmount[workerType] {
			var workerId string
			workerId, err = makeWorkerId(workerType, workerNo)

			if err != nil {
				return err
			}

			model.workerPos[workerId] = 0
			model.workstationQueue = append(model.workstationQueue, false)
			model.storageQueue = append(model.storageQueue, false)
			model.workerIds[workerId] = i
			model.workers[i] = initWorker(workerId, workerType, model.storage, model.workstations, prog, cfg)
			i++
		}
	}

	if _, err = prog.Run(); err != nil {
		return err
	}
	// Clean up
	defer func() {
		for w := range model.workers {
			model.workers[w].done <- true
		}
		model.storage.hasPalletCleared.Broadcast() // make sure no worker is waiting for storage conditional variable
	}()

	return nil
}
