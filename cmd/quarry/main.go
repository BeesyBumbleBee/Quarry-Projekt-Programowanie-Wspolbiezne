package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

/*
\u263B - postać
\U0001F3ED - fabryka
\U0001F528 - młotek
\u26CF - kilof
\u2692 - narzędzia
*/

type model struct {
	distance int
	pos      int
}

func initialModel() model {
	return model{
		distance: 10,
		pos:      0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "left", "h":
			if m.pos > 0 {
				m.pos--
			}

		case "right", "l":
			if m.pos < m.distance {
				m.pos++
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "\U0001F3ED"
	for i := 0; i < m.distance; i++ {
		if i == m.pos {
			s += fmt.Sprintf("\u263B")
		} else {
			s += fmt.Sprintf(" ")
		}
	}

	if m.pos == m.distance {
		s += "\u263B"
	} else {
		s += "\u2692"
	}

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}
