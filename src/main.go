package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type typedAttempt struct {
	key     string
	correct bool
}

type model struct {
	toType []string
	typed  []typedAttempt
}

func (m model) keysTyped() string {
	res := []string{}
	for _, v := range m.typed {
		res = append(res, v.key)
	}
	return strings.Join(res, "")
}

func initialModel() model {
	return model{
		toType: generateRandomSentence(),
		typed:  []typedAttempt{},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) handleBackspace() {
	if len(m.typed) > 0 {
		last := m.typed[len(m.typed)-1]
		if last.correct {
			m.toType = append([]string{last.key}, m.toType...)
		}
		m.typed = m.typed[:(len(m.typed) - 1)]
	}
}

func (m *model) deleteTypedWhile(pred func(string) bool) {
	for len(m.typed) > 0 && pred(m.typed[len(m.typed)-1].key) {
		m.handleBackspace()
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		keyTyped := msg.String()

		switch keyTyped {
		case "ctrl+c":
			return m, tea.Quit

		case "backspace":
			m.handleBackspace()

		case "alt+backspace", "ctrl+w":
			m.deleteTypedWhile(func(c string) bool {
				return c == " "
			})
			m.deleteTypedWhile(func(c string) bool {
				return c != " "
			})

		case "ctrl+u", "cmd+backspace":
			m.typed = []typedAttempt{}

		default:
			correct := false
			if keyTyped == string(m.toType[0]) {
				m.toType = m.toType[1:]
				correct = true
			}
			m.typed = append(m.typed, typedAttempt{key: keyTyped, correct: correct})
			if len(m.toType) == 0 {
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

var correct = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#24273a")).
	Background(lipgloss.Color("#a6da95"))

var incorrect = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#24273a")).
	Background(lipgloss.Color("#ed8796"))

func (m model) View() string {
	res := ""
	for _, v := range m.typed {
		if v.correct {
			res += correct.Render(v.key)
		} else {
			res += incorrect.Render(v.key)
		}
	}
	res += strings.Join(m.toType, "")
	return res
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
