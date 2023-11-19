package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

const (
	gameStateMain = 0
	gameStatePlay = 1
	gameStateOver = 2

	gameModeAgility    = "Agility"
	gameModeMemory     = "Memory"
	gameModePerception = "Perception"
	gameModeLogic      = "Logic"

	arrowUpChar    = "↑"
	arrowDownChar  = "↓"
	arrowLeftChar  = "←"
	arrowRightChar = "→"
)

var (
	styleCenter = lipgloss.NewStyle().
		Bold(true).
		Width(60).
		Padding(1).
		Margin(1).
		Align(lipgloss.Center).
		BorderStyle(lipgloss.ThickBorder())
)

type model struct {
	GameState int
	GameMode  string

	Score               int
	AgilityCurrentArrow string

	Ticks    int
	Frames   int
	Progress float64
	Loaded   bool
	Quitting bool
}

func main() {
	initialModel := model{0, "", 0, "", 10, 0, 0, false, false}
	p := tea.NewProgram(initialModel, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}

func (m model) Init() tea.Cmd {
	return tick()
}

// Main update, which calls the appropriate sub-update
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if (k == "q" || k == "esc" || k == "ctrl+c") && m.GameState != gameStateOver {
			m.Quitting = true
			return m, tea.Quit
		}
	}

	// Hand off the message and model to the appropriate update function
	// for the appropriate view based on the current state
	switch m.GameState {
	case gameStateMain:
		return updateMainMenu(msg, m)
	case gameStateOver:
		return updateGameOver(msg, m)
	case gameStatePlay:
		switch m.GameMode {
		case gameModeAgility:
			return updateAgility(msg, m)
		}
	}

	return nil, nil
}

// Main view, which calls the appropriate sub-view
func (m model) View() string {
	result := ""

	switch m.GameState {
	case gameStateMain:
		result = viewMainMenu(m)
	case gameStateOver:
		result = viewGameOver(m)
	case gameStatePlay:
		result = viewAgility(m)
	}

	return styleCenter.Render(result)
}

// Sub-update functions

func updateMainMenu(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "w", "up":
			m.GameMode = gameModeAgility
			m.GameState = gameStatePlay

		case "s", "down":
			m.GameMode = gameModeMemory
			m.GameState = gameStatePlay
		}
	}

	return m, frame()
}

func updateGameOver(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		default:
			m.GameState = gameStateMain
			m.GameMode = ""
		}
	}

	return m, frame()
}

func updateAgility(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "w", "up":
			if m.AgilityCurrentArrow == arrowUpChar {
				m.AgilityCurrentArrow = ""
				m.Score += 1
			} else {
				m.GameState = gameStateOver
				return m, nil
			}

		case "s", "down":
			if m.AgilityCurrentArrow == arrowDownChar {
				m.AgilityCurrentArrow = ""
				m.Score += 1
			} else {
				m.GameState = gameStateOver
				return m, nil
			}

		case "a", "left":
			if m.AgilityCurrentArrow == arrowLeftChar {
				m.AgilityCurrentArrow = ""
				m.Score += 1
			} else {
				m.GameState = gameStateOver
				return m, nil
			}

		case "d", "right":
			if m.AgilityCurrentArrow == arrowRightChar {
				m.AgilityCurrentArrow = ""
				m.Score += 1
			} else {
				m.GameState = gameStateOver
				return m, nil
			}
		}
	}

	if m.AgilityCurrentArrow == "" {
		num, _ := rand.Int(rand.Reader, big.NewInt(4))

		switch num.Int64() {
		case 0:
			m.AgilityCurrentArrow = arrowUpChar
		case 1:
			m.AgilityCurrentArrow = arrowDownChar
		case 2:
			m.AgilityCurrentArrow = arrowLeftChar
		case 3:
			m.AgilityCurrentArrow = arrowRightChar
		}
	}

	return m, nil
}

// Sub-views

func viewMainMenu(m model) string {
	tpl := "%s\n\n"
	tpl += "Press q, esc, or ctrl+c to quit"

	choices := fmt.Sprintf(
		" %s%s\n %s%s\n %s%s\n %s%s",
		arrowUpChar,
		" Agility: Press the arrows keys you see",
		arrowDownChar,
		" Memory: Remember the arrow key order",
		arrowLeftChar,
		" Perception: Choose the arrow key that's most numerous",
		arrowRightChar,
		" Logic: Deduce the arrow key next in the pattern",
	)

	return fmt.Sprintf(tpl, choices)
}

func viewGameOver(m model) string {
	tpl := "Good job!\n\n"
	tpl += "Game: %s\n"
	tpl += "Score: %d\n\n"
	tpl += "Press any key to continue"

	return fmt.Sprintf(tpl, m.GameMode, m.Score)
}

func viewAgility(m model) string {
	tpl := "  %s\n %s%s%s\n  %s\n\n\n"
	tpl += "Score: %d"

	return fmt.Sprintf(tpl, m.AgilityCurrentArrow, m.AgilityCurrentArrow, m.AgilityCurrentArrow, m.AgilityCurrentArrow, m.AgilityCurrentArrow, m.Score)
}

// Utils

type (
	tickMsg  struct{}
	frameMsg struct{}
)

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func frame() tea.Cmd {
	return tea.Tick(time.Second/60, func(time.Time) tea.Msg {
		return frameMsg{}
	})
}
