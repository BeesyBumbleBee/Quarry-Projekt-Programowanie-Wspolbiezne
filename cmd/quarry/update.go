package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

type tickMsg time.Time

type model struct {
	winWidth, winHeight int
	palletsFilled       int
	cfg                 SimulationConfig
	workers             []Worker
	workerInStorage     string
	workerPos           map[string]int
	storage             *Storage
	workstations        *Workstations
	message             string
	logs                []string
	logsSettings        [2]int
	workerIds           map[string]int
	storageQueue        []bool
	workstationQueue    []bool
	workersAtWork       int
}

const (
	title       = "The Quarry\n"
	mainMessage = "Press q to quit\t Press TAB to modify config"
)

// Init
func initialModel(cfg SimulationConfig, workerAmout int, maxLogs int) model {
	return model{
		palletsFilled:    0,
		winWidth:         0,
		winHeight:        0,
		cfg:              cfg,
		workerPos:        make(map[string]int),
		workers:          make([]Worker, workerAmout),
		workerInStorage:  "",
		storage:          nil,
		workstations:     nil,
		message:          title + mainMessage,
		logsSettings:     [2]int{0, maxLogs},
		logs:             make([]string, maxLogs),
		workerIds:        make(map[string]int),
		storageQueue:     []bool{},
		workstationQueue: []bool{},
		workersAtWork:    0,
	}
}

func (m *model) log(msg string, colorFg [3]int, colorBg ...[3]int) {
	clearColor := "\033[0m"
	colorCodeFg := fmt.Sprintf("\033[38;2;%d;%d;%dm", colorFg[0], colorFg[1], colorFg[2])
	colorCodeBg := ""
	if len(colorBg) > 0 {
		colorCodeBg = fmt.Sprintf("\033[48;2;%d;%d;%dm", colorBg[0][0], colorBg[0][1], colorBg[0][2])
	}
	if m.logsSettings[0] != m.logsSettings[1] {
		m.logs[m.logsSettings[0]] = colorCodeFg + colorCodeBg + timeStamp() + " : " + msg + clearColor
		m.logsSettings[0]++
		return
	}
	for i := range m.logsSettings[1] - 1 {
		m.logs[i] = m.logs[i+1]
	}
	m.logs[m.logsSettings[1]-1] = colorCodeFg + colorCodeBg + timeStamp() + " : " + msg + clearColor

}

func timeStamp() string {
	return time.Now().Format("15:04:05")
}

func (m *model) Init() tea.Cmd {
	for worker := range m.workers {
		go m.workers[worker].Work()
	}
	return nil
}

// Update model
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil
	switch msg := msg.(type) {
	case tickMsg:
		m.message = title + mainMessage
		break
	case tea.WindowSizeMsg:
		m.winWidth = msg.Width
		m.winHeight = msg.Height
		break
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
		break
	case workerMoveMsg:
		m.workerPos[msg.workerId] = msg.position + (50 * msg.workerType)
		break
	case workerAtStorage:
		m.storageQueue[m.workerIds[msg.workerId]] = true
		m.log(msg.workerId+" is at storage.", [3]int{224, 224, 224})
		break
	case workerAtWork:
		m.workstationQueue[m.workerIds[msg.workerId]] = true
		m.workersAtWork++
		m.log(msg.workerId+" is at work.", [3]int{160, 160, 160})
		break
	case workerWorking:
		m.workstationQueue[m.workerIds[msg.workerId]] = false
		m.log(msg.workerId+" started working.", [3]int{180, 180, 180})
	case workerFinishedWork:
		m.workstationQueue[m.workerIds[msg.workerId]] = false
		m.workersAtWork--
		m.log(msg.workerId+" finished work.", [3]int{160, 160, 160})
		break
	case storageWorkerPlacing:
		m.storageQueue[m.workerIds[msg.workerId]] = false
		m.log(msg.workerId+" is placing stone block.", [3]int{255, 255, 204})
		m.workerInStorage = msg.workerId
		break
	case storageWorkerFinishedPlacing:
		m.log(msg.workerId+" is finished placing stone block.", [3]int{204, 255, 204})
		m.workerInStorage = ""
		break
	case storageWorkerCantPlace:
		m.storageQueue[m.workerIds[msg.workerId]] = true
		m.log(msg.workerId+" can't place stone block.", [3]int{255, 0, 0})
		m.workerInStorage = ""
		break
	case palletFullMsg:
		m.palletsFilled++
		m.log("Pallet Full. Replacing...", [3]int{0, 0, 0}, [3]int{0, 153, 0})
		m.message = title + fmt.Sprintf("\033[38;2;252;189;0mPallet #%-2d was filled in %3d seconds!\033[0m",
			m.palletsFilled, int(msg.timeTook.Seconds()))
		cmd = tickCmd()
		break
	case placeInsulationMsg:
		m.log("Layer is full. Placing Insulation...", [3]int{0, 0, 0}, [3]int{0, 153, 0})
		break
	}

	return m, cmd
}

// Helper Functions
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
