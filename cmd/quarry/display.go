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
			fmt.Sprintf("Currently placing: %2s", m.workerInStorage),
			""})

	return storage
}

func (m *model) printWorkstations() string {
	workstation := ""
	i := m.cfg.QuarryWorkplaces - 1/3
	for station := 0; station < m.cfg.QuarryWorkplaces; station++ {
		workerAtStation := " "
		if station < m.workersAtWork {
			workerAtStation = "\u263B"
		}
		workstation += fmt.Sprintf(" %s \u2692  ", workerAtStation)
		if i == 0 {
			workstation += "\n"
			i = m.cfg.QuarryWorkplaces - 1/3
		}
		i--
	}
	return workstation
}

func (m *model) View() string {
	screen := ""

	maxWidth := 0
	terminalWidth = m.winWidth - 4
	terminalHeight = m.winHeight

	// Storage
	storage := m.printStorageCells(m.storage.cells, m.storage.verticalBars, m.storage.horizontalBars)

	// Road
	road := ""
	roadCells := make([]rune, 150)
	for i := 0; i < 150; i++ {
		roadCells[i] = '.'
	}

	for _, pos := range m.workerPos {
		roadCells[pos] = '\u263B'
	}

	for i := 0; i < 150; i++ {
		if (i%50 == 0) && (i != 0) {
			road += "\n\n\n"
		}
		road += string(roadCells[i])
	}

	// Workstation
	workstation := m.printWorkstations()

	// MessageBox
	messageBox := m.message

	// StorageQueue
	storageQ := "Waiting at storage: "
	for worker, inQueue := range m.storageQueue {
		if inQueue {
			storageQ += m.workers[worker].id + " "
		}
	}

	// WorkstationQueue
	workstationQ := "Waiting at workstations: "
	for worker, inQueue := range m.workstationQueue {
		if inQueue {
			workstationQ += m.workers[worker].id + " "
		}
	}

	// Config
	config := m.cfg.printConfig()
	// TODO: Add dynamic max log amount based on terminal height

	// Logs
	logs := "LOGS:\n"
	for log := range m.logs {
		logs += m.logs[log] + "\n"
	}

	// Creating views
	storageWidth := (terminalWidth * 2 / 5) - lipgloss.Width(storage)
	storage = windowStyle.PaddingLeft(storageWidth / 2).PaddingRight(storageWidth / 2).Render(storage)
	storageHeight := lipgloss.Height(storage) - 2

	roadWidth := (terminalWidth * 2 / 5) - lipgloss.Width(road)
	road = lipgloss.Place(roadWidth, 9, 0, 1, road)
	road = windowStyle.PaddingRight(roadWidth / 2).PaddingLeft(roadWidth / 2).Height(storageHeight).Render(road)

	workstationWidth := (terminalWidth * 1 / 5) - lipgloss.Width(workstation)
	workstation = windowStyle.Height(storageHeight).PaddingRight(workstationWidth / 2).PaddingLeft(workstationWidth / 2).Render(workstation)

	maxWidth += lipgloss.Width(storage) + lipgloss.Width(workstation) + lipgloss.Width(road) - 3

	messageBox = lipgloss.Place(maxWidth, 2, 0.5, 0, messageBox)
	messageBox = windowStyle.Width(maxWidth + 1).Render(messageBox)
	storageQ = windowStyle.Width(maxWidth / 2).Render(storageQ)
	workstationQ = windowStyle.Width(maxWidth/2 - 1).Render(workstationQ)
	config = windowStyle.Width(maxWidth * 4 / 7).Height(terminalHeight - lipgloss.Height(messageBox) - lipgloss.Height(workstationQ) - lipgloss.Height(workstation) - 2).Render(config)
	logs = windowStyle.Width(maxWidth * 3 / 7).Height(lipgloss.Height(config) - 2).Render(logs)

	// Gluing views together
	screen1 := messageBox
	screen2 := lipgloss.JoinHorizontal(0, storageQ, workstationQ)
	screen3 := lipgloss.JoinHorizontal(0, storage, road, workstation)
	screen4 := lipgloss.JoinHorizontal(0, config, logs)

	screen = lipgloss.JoinVertical(0, screen1, screen2, screen3, screen4)

	return screen
}
