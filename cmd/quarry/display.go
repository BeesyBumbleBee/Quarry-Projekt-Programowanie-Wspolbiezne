package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strconv"
)

/*
\u263B - postać
\U0001F3ED - fabryka
\u2692 - narzędzia
*/

var (
	windowStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("#FFF")).
			Padding(1, 2, 1, 2)
	terminalWidth  = 0
	terminalHeight = 0
	grey           = "\033[48;2;120;120;120m\033[38;2;0;0;0m"
	darkGrey       = "\033[48;2;119;119;100m\033[38;2;0;0;0m"
	brown          = "\033[48;2;153;76;0m\033[38;2;102;51;0m"
	normal         = "\033[0m"
)

func printPalletRow(cells [3]int, verticalBars [2]bool, msg [3]string) string {
	A := [6]string{}
	B := [8]string{}

	for i, cell := range cells {
		if cell == 0 {
			A[i*2] = normal
			A[i*2+1] = " "
		} else {
			A[i*2] = grey
			A[i*2+1] = strconv.Itoa(cell)
		}
	}
	for i, bar := range verticalBars {
		if bar {
			B[i*4] = ""
			B[i*4+1] = darkGrey
			B[i*4+2] = " "
			B[i*4+3] = grey
		} else {
			B[i*4] = normal
			B[i*4+1] = brown
			B[i*4+2] = "|"
			B[i*4+3] = normal
		}
	}

	res := fmt.Sprintf("%s   %s%s%s%s%s   %s%s%s%s%s   %s     || %s\n",
		A[0], B[0], B[1], B[2], B[3], A[2], B[4], B[5], B[6], B[7], A[4], normal, msg[0])
	res += fmt.Sprintf("%s %s %s%s%s%s%s %s %s%s%s%s%s %s %s     || %s\n",
		A[0], A[1], B[0], B[1], B[2], B[3], A[2], A[3], B[4], B[5], B[6], B[7], A[4], A[5], normal, msg[1])
	res += fmt.Sprintf("%s   %s%s%s%s%s   %s%s%s%s%s   %s     || %s\n",
		A[0], B[0], B[1], B[2], B[3], A[2], B[4], B[5], B[6], B[7], A[4], normal, msg[2])
	return res
}

func printPalletHorizontalBars(horizontalBars [3]bool, msg string) string {
	B := [9]string{}
	for i, bar := range horizontalBars {
		if bar {
			B[i*3] = darkGrey
			B[i*3+1] = " "
			B[i*3+2] = brown
		} else {
			B[i*3] = ""
			B[i*3+1] = "-"
			B[i*3+2] = ""
		}
	}
	res := fmt.Sprintf("%s%s%s%s%s%s+%s%s%s%s%s+%s%s%s%s%s%s     || %s\n",
		brown, B[0], B[1], B[1], B[1], B[2], B[3], B[4], B[4], B[4], B[5], B[6], B[7], B[7], B[7], B[8], normal, msg)
	return res
}

func (m *model) printRoad() string {
	var road string

	roadCells := make([]rune, 150)
	for i := 0; i < 150; i++ {
		roadCells[i] = '.'
	}

	for _, pos := range m.workerPos {
		if pos != -1 {
			roadCells[pos] = '\u263B'
		}
	}

	for i := 0; i < 150; i++ {
		if (i%50 == 0) && (i != 0) {
			road += "\n\n\n"
		}
		road += string(roadCells[i])
	}
	return road
}

func (m *model) printStorageCells(cells [3][3]int, verticalBars [3][2]bool, horizontalBars [2][3]bool) string {

	storage := printPalletRow(cells[0], verticalBars[0],
		[3]string{
			fmt.Sprintf("Pallet #%-2d", m.palletsFilled+1),
			fmt.Sprintf("Layer: %1d", m.storage.level),
			fmt.Sprintf("Mass limit: %2d", m.cfg.StoneMassesLimits[m.storage.level])})
	storage += printPalletHorizontalBars(horizontalBars[0],
		fmt.Sprintf("Current mass: %2d", m.storage.totalMass))
	storage += printPalletRow(cells[1], verticalBars[1],
		[3]string{
			"",
			"",
			""})
	storage += printPalletHorizontalBars(horizontalBars[1], "")
	storage += printPalletRow(cells[2], verticalBars[2],
		[3]string{
			"",
			fmt.Sprintf("Currently placing: %4s", m.workerInStorage),
			""})

	return storage
}

func (m *model) printWorkstations(width int) string {
	workstation := ""
	bytes := 0
	for station := 0; station < m.cfg.QuarryWorkplaces; station++ {
		workerAtStation := " "
		if station < m.workersAtWork {
			workerAtStation = "\u263B"
		}
		workstation += fmt.Sprintf(" %s \u2692  ", workerAtStation)
		bytes += 6
		if bytes%width == 0 {
			workstation += "\n"
		}
	}
	return workstation
}

func (m *model) View() string {
	screen := ""

	terminalWidth = m.winWidth - 8
	terminalHeight = m.winHeight - 4

	// Storage
	storageWidth := terminalWidth / 4
	storage := m.printStorageCells(m.storage.cells, m.storage.verticalBars, m.storage.horizontalBars)
	storage = windowStyle.PaddingLeft(5).Width(storageWidth).Render(storage)
	storageHeight := lipgloss.Height(storage) - 2

	// Road
	roadWidth := terminalWidth / 4
	road := m.printRoad()
	road = lipgloss.Place(roadWidth, 9, 0, 1, road)
	road = windowStyle.PaddingLeft(3).Width(roadWidth).Height(storageHeight).Render(road)

	// StorageQueue
	storageQWidth := terminalWidth * 3 / 8
	storageQ := "Waiting at storage: "
	for worker, inQueue := range m.storageQueue {
		if inQueue {
			storageQ += m.workers[worker].id + " "
		}
	}
	storageQ = windowStyle.Width(storageQWidth).Render(storageQ)

	// WorkstationQueue
	workstationQWidth := terminalWidth*2/8 + 1

	workstationQ := "Waiting at workstations: "
	for worker, inQueue := range m.workstationQueue {
		if inQueue {
			workstationQ += m.workers[worker].id + " "
		}
	}
	workstationQ = windowStyle.Width(workstationQWidth).Render(workstationQ)

	// Workstation
	workstationWidth := terminalWidth / 8
	workstation := m.printWorkstations(workstationWidth - 4)
	workstation = windowStyle.Width(workstationWidth).Height(storageHeight).Render(workstation)

	// MessageBox
	messageBox := m.message
	messageBox = lipgloss.Place(terminalWidth, 2, 0.5, 0, messageBox)
	messageBox = windowStyle.Width(terminalWidth + 4).Render(messageBox)

	// Logs
	logsWidth := terminalWidth * 3 / 8
	logsHeight := terminalHeight - lipgloss.Height(messageBox) + 2
	logs := "LOGS:\n"
	for i := 0; i < m.logsSize-1; i++ {
		logs += m.logs[i]
		if i != m.logsSize-1 {
			logs += "\n"
		}
	}
	logs = windowStyle.Width(logsWidth).Height(logsHeight).Render(logs)
	m.logsDesiredCapacity = lipgloss.Height(logs) - 4

	// Config
	configWidth := terminalWidth*5/8 + 3
	configHeight := terminalHeight - storageHeight - lipgloss.Height(messageBox) - lipgloss.Height(storageQ)

	config := m.cfg.printConfig()
	config = windowStyle.Width(configWidth).Height(configHeight).Render(config)

	// Gluing views together
	screen1 := messageBox
	screen2 := lipgloss.JoinHorizontal(0, storageQ, workstationQ)
	screen3 := lipgloss.JoinHorizontal(0, storage, road, workstation)
	screen4 := lipgloss.JoinHorizontal(0, config)

	screen5 := lipgloss.JoinVertical(0, screen2, screen3, screen4)
	screen6 := lipgloss.JoinHorizontal(0, screen5, logs)
	screen = lipgloss.JoinVertical(1, screen1, screen6)

	return screen
}
