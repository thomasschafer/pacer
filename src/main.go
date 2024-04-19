package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type typedAttempt struct {
	key     string
	correct bool
}

type page int

const (
	typingPage page = iota
	helpPage
	resultsPage
)

type model struct {
	toType         []string
	testWordLength int
	typed          []typedAttempt

	startTime   time.Time
	finishTime  time.Time
	testStarted bool

	currentPage page

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
	m := model{}
	m.resetState()
	return m
}

func (m *model) resetState() {
	m.toType = generateRandomSentence(top1000)
	m.testWordLength = len(strings.Fields(strings.Join(m.toType, "")))
	m.typed = []typedAttempt{}
	m.startTime = time.Now()
	m.finishTime = time.Now() // If only Go had a Maybe type
	m.currentPage = typingPage
	m.testStarted = false
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

var header = lipgloss.NewStyle().
	BorderBottom(true).
	BorderStyle(lipgloss.NormalBorder()).
	Padding(0, 2).
	Align(lipgloss.Center)

var helpTextStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	PaddingLeft(2).
	PaddingRight(4).
	Align(lipgloss.Left)

var textBody = lipgloss.NewStyle().
	Padding(2).
	Align(lipgloss.Center)

var tableRow = lipgloss.NewStyle().
	Width(20).
	Padding(0, 1)

func buildTable(rows [][]string) *table.Table {
	return table.New().
		StyleFunc(func(row, col int) lipgloss.Style { return tableRow }).
		Rows(rows...)
}

var allPagesHelp = buildTable([][]string{
	{"?", "help"},
	{"ctrl+c", "quit"},
})

var typingTestHelp = buildTable([][]string{
	{"tab", "new test"},
})

var otherPagesHelp = buildTable([][]string{
	{"enter", "new test"},
	{"q", "quit"},
})

func (m model) content() string {
	var res strings.Builder
	res.WriteString(header.Render("PACER\nPress ? for help"))

	switch m.currentPage {
	case typingPage:
		res.WriteString(strings.Repeat("\n", max(0, m.viewport.Height/2-3))) // TODO: deduplicate this and remove hard-coded nums
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
	case helpPage:
		res.WriteString(strings.Repeat("\n", max(0, m.viewport.Height/2-12)))
		res.WriteString(
			"\nAll pages\n" +
				allPagesHelp.String() +
				"\n\nTyping test\n" +
				typingTestHelp.String() +
				"\n\nOther pages\n" +
				otherPagesHelp.String(),
		)
	case resultsPage:
		res.WriteString(strings.Repeat("\n", max(0, m.viewport.Height/2-5)))
		timeTakenSecs := float32(m.finishTime.Sub(m.startTime)) / 1e9
		res.WriteString(fmt.Sprintf("Time taken: %.2fs\n", timeTakenSecs))
		res.WriteString(fmt.Sprintf("Words per minute: %.2f\n\n", float32(m.testWordLength*60)/timeTakenSecs))
		res.WriteString("Press enter to start a new test")
	}

	return textBody.
		Width(m.viewport.Width).
		Render(res.String())
}

func (m *model) endTest() {
	m.currentPage = resultsPage
	m.finishTime = time.Now()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		keyTyped := msg.String()

		switch m.currentPage {
		case typingPage:
			if !m.testStarted {
				m.testStarted = true
				m.startTime = time.Now()
			}
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

			case "tab":
				m.resetState()

			case "?":
				m.currentPage = helpPage

			default:
				for _, c := range keyTyped {
					if len(m.toType) == 0 {
						m.endTest()
					}
					correct := string(c) == m.toType[0]
					if correct {
						m.toType = m.toType[1:]
					}
					m.typed = append(m.typed, typedAttempt{key: string(c), correct: correct})
				}
				if len(m.toType) == 0 {
					m.endTest()
				}
			}
		case helpPage, resultsPage:
			switch keyTyped {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "enter", "esc", "tab":
				m.resetState()
			case "?":
				m.currentPage = helpPage
			}
		}

	case tea.WindowSizeMsg:
		if !m.ready {
			// Wait until we've received the window dimensions before
			// we can initialize the viewport
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
