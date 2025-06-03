package main

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strconv"
)

const (
	workers_amountA = iota
	workers_amountB
	workers_amountC
	stones_massesA
	stones_massesB
	stones_massesC
	stone_masses_limits1
	stone_masses_limits2
	stone_masses_limits3
	stones_extraction_timeAMIN
	stones_extraction_timeAMAX
	stones_extraction_timeBMIN
	stones_extraction_timeBMAX
	stones_extraction_timeCMIN
	stones_extraction_timeCMAX
	time_to_travel_emptyMIN
	time_to_travel_emptyMAX
	time_to_travel_fullMIN
	time_to_travel_fullMAX
	time_to_place_stoneMIN
	time_to_place_stoneMAX
	time_to_place_insulationMIN
	time_to_place_insulationMAX
	time_to_change_palletMIN
	time_to_change_palletMAX
	quarry_workplaces
)

const (
	selector = '>'
)

func numberValidator(s string) error {
	num, err := strconv.Atoi(s)
	if num < 0 {
		err = errors.New("number has to be a positive integer")
		return err
	}
	return err
}

type (
	errMsg error
)

type MainMenu struct {
	cfg                 SimulationConfig
	inputs              []textinput.Model
	configInputFocused  int
	menuOptionFocused   int
	err                 error
	inMenu              bool
	inConfig            bool
	menuOptions         []rune
	exitApp             bool
	winWidth, winHeight int
}

func initialMainMenu() MainMenu {
	var inputs = make([]textinput.Model, 26)
	for input := range inputs {
		inputs[input] = textinput.New()
		inputs[input].Validate = numberValidator
		inputs[input].CharLimit = 6
		inputs[input].Width = 8
		inputs[input].Placeholder = "0"
		inputs[input].Prompt = " "
	}
	inputs[0].Focus()
	menuOptions := make([]rune, 2)
	for i := range menuOptions {
		menuOptions[i] = ' '
	}
	menuOptions[0] = selector

	return MainMenu{
		inputs:             inputs,
		winWidth:           0,
		winHeight:          0,
		configInputFocused: 0,
		menuOptionFocused:  0,
		err:                nil,
		cfg:                SimulationConfig{},
		inConfig:           false,
		exitApp:            false,
		menuOptions:        menuOptions,
	}
}

func startMainMenu() (SimulationConfig, error) {
	var err error = nil
	mainMenu := initialMainMenu()
	mainMenu.cfg, err = getConfig()
	if err != nil {
		mainMenu.cfg = defaultConfig()
	}
	menu := tea.NewProgram(&mainMenu, tea.WithAltScreen(), tea.WithoutSignalHandler()) // implement ^C as interrupt or other way to close the program

	_, err = menu.Run()
	cfg, err := configFromInputs(mainMenu.inputs)
	if err != nil {
		return defaultConfig(), nil
	}

	if mainMenu.exitApp {
		err = errors.New("exit app")
		return SimulationConfig{}, err
	}
	return cfg, err
}

func (m *MainMenu) Init() tea.Cmd {
	m.inputs[workers_amountA].SetValue(strconv.Itoa(m.cfg.WorkersAmount[0]))
	m.inputs[workers_amountB].SetValue(strconv.Itoa(m.cfg.WorkersAmount[1]))
	m.inputs[workers_amountC].SetValue(strconv.Itoa(m.cfg.WorkersAmount[2]))
	m.inputs[stones_extraction_timeAMIN].SetValue(strconv.Itoa(m.cfg.StonesExtractionTime[0][0]))
	m.inputs[stones_extraction_timeAMAX].SetValue(strconv.Itoa(m.cfg.StonesExtractionTime[0][1]))
	m.inputs[stones_extraction_timeBMIN].SetValue(strconv.Itoa(m.cfg.StonesExtractionTime[1][0]))
	m.inputs[stones_extraction_timeBMAX].SetValue(strconv.Itoa(m.cfg.StonesExtractionTime[1][1]))
	m.inputs[stones_extraction_timeCMIN].SetValue(strconv.Itoa(m.cfg.StonesExtractionTime[2][0]))
	m.inputs[stones_extraction_timeCMAX].SetValue(strconv.Itoa(m.cfg.StonesExtractionTime[2][1]))
	m.inputs[stones_massesA].SetValue(strconv.Itoa(m.cfg.StonesMasses[0]))
	m.inputs[stones_massesB].SetValue(strconv.Itoa(m.cfg.StonesMasses[1]))
	m.inputs[stones_massesC].SetValue(strconv.Itoa(m.cfg.StonesMasses[2]))
	m.inputs[stone_masses_limits1].SetValue(strconv.Itoa(m.cfg.StoneMassesLimits[0]))
	m.inputs[stone_masses_limits2].SetValue(strconv.Itoa(m.cfg.StoneMassesLimits[1]))
	m.inputs[stone_masses_limits3].SetValue(strconv.Itoa(m.cfg.StoneMassesLimits[2]))
	m.inputs[time_to_travel_emptyMIN].SetValue(strconv.Itoa(m.cfg.TimeToTravelEmpty[0]))
	m.inputs[time_to_travel_emptyMAX].SetValue(strconv.Itoa(m.cfg.TimeToTravelEmpty[1]))
	m.inputs[time_to_travel_fullMIN].SetValue(strconv.Itoa(m.cfg.TimeToTravelFull[0]))
	m.inputs[time_to_travel_fullMAX].SetValue(strconv.Itoa(m.cfg.TimeToTravelFull[1]))
	m.inputs[time_to_place_stoneMIN].SetValue(strconv.Itoa(m.cfg.TimeToPlaceStone[0]))
	m.inputs[time_to_place_stoneMAX].SetValue(strconv.Itoa(m.cfg.TimeToPlaceStone[1]))
	m.inputs[time_to_place_insulationMIN].SetValue(strconv.Itoa(m.cfg.TimeToPlaceInsulation[0]))
	m.inputs[time_to_place_insulationMAX].SetValue(strconv.Itoa(m.cfg.TimeToPlaceInsulation[1]))
	m.inputs[time_to_change_palletMIN].SetValue(strconv.Itoa(m.cfg.TimeToChangePallet[0]))
	m.inputs[time_to_change_palletMAX].SetValue(strconv.Itoa(m.cfg.TimeToChangePallet[1]))
	m.inputs[quarry_workplaces].SetValue(strconv.Itoa(m.cfg.QuarryWorkplaces))
	return nil
}

func (m *MainMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.winWidth = msg.Width
		m.winHeight = msg.Height
		break
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.inConfig {
				if m.configInputFocused == len(m.inputs)-1 {
					return m, tea.Quit
				}
				m.nextInput()
			} else {
				switch m.menuOptionFocused {
				case 0:
					return m, tea.Quit
				case 1:
					m.inConfig = true
				}
			}
		case tea.KeyCtrlQ:
			m.exitApp = true
			return m, tea.Quit
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			m.inConfig = false
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		case tea.KeyUp:
			m.menuOptionFocused = (m.menuOptionFocused - 1) % len(m.menuOptions)
			for i := range m.menuOptions {
				if i == m.menuOptionFocused {
					m.menuOptions[i] = selector
				} else {
					m.menuOptions[i] = ' '
				}
			}
		case tea.KeyDown:
			m.menuOptionFocused = (m.menuOptionFocused + 1) % len(m.menuOptions)
			for i := range m.menuOptions {
				if i == m.menuOptionFocused {
					m.menuOptions[i] = selector
				} else {
					m.menuOptions[i] = ' '
				}
			}
		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.configInputFocused].Focus()

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)

}

func (m *MainMenu) View() string {
	var view string
	if m.inConfig {
		view = fmt.Sprintf(
			`%s
%s

Workers:                 A = %s B = %s C = %s 
Stone Masses:            A = %s B = %s C = %s

Stone Masses Limits:     1.Layer = %s 2.Layer = %s 3.Layer = %s

Stones Extraction Times: A (min = %s max = %s) B (min = %s max = %s) C (min = %s max = %s)

Time To Travel Empty:       min = %s max = %s
Time To Travel Full:        min = %s max = %s
Time To Place Stone:        min = %s max = %s
Time To Place Insulation:   min = %s max = %s
Time To Change Pallet:      min = %s max = %s

Quarry Workplaces = %s

All times are given in ms.
Time to travel empty and time to travel full are times it takes for worker to take 1 step,
there are 50 steps between storage and workstations.

Press TAB \ Shift-TAB to cycle between input fields.
Press ESC to go back to main menu.
Press CTRL+C to start simulation.
`,
			padString("Simulation Configuration", " ", m.winWidth),
			padString("", "=", m.winWidth*2),
			m.inputs[workers_amountA].View(), m.inputs[workers_amountB].View(), m.inputs[workers_amountC].View(),
			m.inputs[stones_massesA].View(), m.inputs[stones_massesB].View(), m.inputs[stones_massesC].View(),
			m.inputs[stone_masses_limits1].View(), m.inputs[stone_masses_limits2].View(), m.inputs[stone_masses_limits3].View(),
			m.inputs[stones_extraction_timeAMIN].View(), m.inputs[stones_extraction_timeAMAX].View(), m.inputs[stones_extraction_timeBMIN].View(), m.inputs[stones_extraction_timeBMAX].View(), m.inputs[stones_extraction_timeCMIN].View(), m.inputs[stones_extraction_timeCMAX].View(),
			m.inputs[time_to_travel_emptyMIN].View(), m.inputs[time_to_travel_emptyMAX].View(),
			m.inputs[time_to_travel_fullMIN].View(), m.inputs[time_to_travel_fullMAX].View(),
			m.inputs[time_to_place_stoneMIN].View(), m.inputs[time_to_place_stoneMAX].View(),
			m.inputs[time_to_place_insulationMIN].View(), m.inputs[time_to_place_insulationMAX].View(),
			m.inputs[time_to_change_palletMIN].View(), m.inputs[time_to_change_palletMAX].View(),
			m.inputs[quarry_workplaces].View())
	} else {
		view = fmt.Sprintf(
			`%s
%s

%s
%s


Cycle between options using up and down arrow keys.
Press Enter to select option.
Press Ctrl-Q to quit.
`,
			padString("The Quarry - Main Menu", " ", m.winWidth),
			padString("", "=", m.winWidth*2),
			padString(fmt.Sprintf("%cStart Simulation    ", m.menuOptions[0]), " ", m.winWidth),
			padString(fmt.Sprintf("%cChange Config       ", m.menuOptions[1]), " ", m.winWidth))
	}
	return view
}

// nextInput focuses the next input field
func (m *MainMenu) nextInput() {
	m.configInputFocused = (m.configInputFocused + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
func (m *MainMenu) prevInput() {
	m.configInputFocused--
	// Wrap around
	if m.configInputFocused < 0 {
		m.configInputFocused = len(m.inputs) - 1
	}
}

func configFromInputs(inputs []textinput.Model) (SimulationConfig, error) {
	var cfg SimulationConfig
	var err error

	cfg.WorkersAmount[0], _ = strconv.Atoi(inputs[workers_amountA].Value())
	cfg.WorkersAmount[1], _ = strconv.Atoi(inputs[workers_amountB].Value())
	cfg.WorkersAmount[2], _ = strconv.Atoi(inputs[workers_amountC].Value())
	cfg.StonesExtractionTime[0][0], _ = strconv.Atoi(inputs[stones_extraction_timeAMIN].Value())
	cfg.StonesExtractionTime[0][1], _ = strconv.Atoi(inputs[stones_extraction_timeAMAX].Value())
	cfg.StonesExtractionTime[1][0], _ = strconv.Atoi(inputs[stones_extraction_timeBMIN].Value())
	cfg.StonesExtractionTime[1][1], _ = strconv.Atoi(inputs[stones_extraction_timeBMAX].Value())
	cfg.StonesExtractionTime[2][0], _ = strconv.Atoi(inputs[stones_extraction_timeCMIN].Value())
	cfg.StonesExtractionTime[2][1], _ = strconv.Atoi(inputs[stones_extraction_timeCMAX].Value())
	cfg.StonesMasses[0], _ = strconv.Atoi(inputs[stones_massesA].Value())
	cfg.StonesMasses[1], _ = strconv.Atoi(inputs[stones_massesB].Value())
	cfg.StonesMasses[2], _ = strconv.Atoi(inputs[stones_massesC].Value())
	cfg.StoneMassesLimits[0], _ = strconv.Atoi(inputs[stone_masses_limits1].Value())
	cfg.StoneMassesLimits[1], _ = strconv.Atoi(inputs[stone_masses_limits2].Value())
	cfg.StoneMassesLimits[2], _ = strconv.Atoi(inputs[stone_masses_limits3].Value())
	cfg.TimeToTravelEmpty[0], _ = strconv.Atoi(inputs[time_to_travel_emptyMIN].Value())
	cfg.TimeToTravelEmpty[1], _ = strconv.Atoi(inputs[time_to_travel_emptyMAX].Value())
	cfg.TimeToTravelFull[0], _ = strconv.Atoi(inputs[time_to_travel_fullMIN].Value())
	cfg.TimeToTravelFull[1], _ = strconv.Atoi(inputs[time_to_travel_fullMAX].Value())
	cfg.TimeToPlaceStone[0], _ = strconv.Atoi(inputs[time_to_place_stoneMIN].Value())
	cfg.TimeToPlaceStone[1], _ = strconv.Atoi(inputs[time_to_place_stoneMAX].Value())
	cfg.TimeToPlaceInsulation[0], _ = strconv.Atoi(inputs[time_to_place_insulationMIN].Value())
	cfg.TimeToPlaceInsulation[1], _ = strconv.Atoi(inputs[time_to_place_insulationMAX].Value())
	cfg.TimeToChangePallet[0], _ = strconv.Atoi(inputs[time_to_change_palletMIN].Value())
	cfg.TimeToChangePallet[1], _ = strconv.Atoi(inputs[time_to_change_palletMAX].Value())
	cfg.QuarryWorkplaces, _ = strconv.Atoi(inputs[quarry_workplaces].Value())

	return cfg, err
}

func padString(str string, padChar string, winWidth int) string {
	padding := ""
	padded := str
	for i := 0; i < (winWidth/2)-len(str); i++ {
		padding += padChar
	}
	padded = padding + padded
	return padded
}
