package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
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

	viewport viewport.Model
	ready    bool
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
		toType: generateRandomSentence(top1000),
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

var correct = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#24273a")).
	Background(lipgloss.Color("#a6da95"))

var incorrect = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#24273a")).
	Background(lipgloss.Color("#ed8796"))

var textBody = lipgloss.NewStyle().
	Padding(2).
	Align(lipgloss.Center)

func (m model) content() string {
	var res strings.Builder
	res.WriteString(strings.Repeat("\n", m.viewport.Height/2-1))
	for _, v := range m.typed {
		if v.correct {
			res.WriteString(correct.Render(v.key))
		} else {
			res.WriteString(incorrect.Render(v.key))
		}
	}
	for _, s := range m.toType {
		res.WriteString(s)
	}

	return textBody.
		Width(m.viewport.Width).
		Render(res.String())
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
			m.deleteTypedWhile(func(c string) bool {
				return true
			})

		default:
			for _, c := range keyTyped {
				correct := string(c) == m.toType[0]
				if correct {
					m.toType = m.toType[1:]
				}
				m.typed = append(m.typed, typedAttempt{key: string(c), correct: correct})
			}

			if len(m.toType) == 0 {
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height)
			m.viewport.YPosition = 10
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height
		}
	}

	m.viewport.SetContent(m.content())

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return m.viewport.View()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
