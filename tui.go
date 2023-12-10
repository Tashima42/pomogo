package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	baseStyle    = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	helpStyle    = blurredStyle.Copy()
)

type tickMsg time.Time

type model struct {
	table    table.Model
	pomo     pomoModel
	progress progressModel
}

type pomoModel struct {
	focus           time.Duration
	focusBreak      time.Duration
	longFocusBreak  time.Duration
	breaksUntilLong int
}

type progressModel struct {
	percent  float64
	progress progress.Model
}

const (
	padding  = 2
	maxWidth = 80
)

func initialModel() model {
	mdl := model{
		pomo: pomoModel{
			focus:           time.Minute * 25,
			focusBreak:      time.Minute * 5,
			longFocusBreak:  time.Minute * 15,
			breaksUntilLong: 4,
		},
		progress: progressModel{
			progress: progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C")),
		},
	}
	columns := []table.Column{
		{Title: "Name", Width: 17},
		{Title: "Value", Width: 5},
	}
	rows := []table.Row{
		{"Focus", mdl.pomo.focus.String()},
		{"Break", mdl.pomo.focusBreak.String()},
		{"Long Break", mdl.pomo.longFocusBreak.String()},
		{"Breaks Until Long", strconv.Itoa(mdl.pomo.breaksUntilLong)},
	}
	mdl.table = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(4),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	mdl.table.SetStyles(s)
	return mdl
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	case tea.WindowSizeMsg:
		m.progress.progress.Width = msg.Width - padding*2 - 4
		if m.progress.progress.Width > maxWidth {
			m.progress.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		m.progress.percent += 0.25
		if m.progress.percent > 1.0 {
			m.progress.percent = 1.0
			return m, tea.Quit
		}
		return m, tickCmd()

	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var b strings.Builder
	b.WriteString(baseStyle.Render(m.table.View()) + "\n")

	pad := strings.Repeat(" ", padding)
	b.WriteString(pad + m.progress.progress.ViewAs(m.progress.percent) + "\n\n")

	b.WriteString(helpStyle.Render("press 'enter' to edit the value") + "\n")
	return b.String()
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
