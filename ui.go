package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

const (
	padding  = 2
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

type tickMsg time.Time

type model struct {
	width       int
	height      int
	quitting    time.Time
	logFileName string
	data        *state
	stopwatch   stopwatch.Model
	spinner     spinner.Model
	progress    progress.Model
}

func newModel(data *state) *model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return &model{
		width:       80, // default width, will be updated on resize
		data:        data,
		logFileName: logFile.Name(),
		stopwatch:   stopwatch.NewWithInterval(time.Millisecond),
		spinner:     s,
		progress:    progress.New(progress.WithDefaultGradient()),
	}
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(tickCmd(), m.spinner.Tick, m.stopwatch.Init())
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if !m.quitting.IsZero() && time.Now().After(m.quitting) {
		log.Debug("quitting...")
		return m, tea.Quit
	}

	quitTime := time.Second / 2
	if m.quitting.IsZero() && m.data.StopDueToFailures() {
		m.quitting = time.Now().Add(quitTime)
		return m, tickCmd()
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.width = min(msg.Width-padding*2-4, maxWidth)
		m.height = msg.Height
		m.progress.Width = m.width
		return m, nil

	case tickMsg:
		if m.quitting.IsZero() && m.progress.Percent() == 1.0 {
			m.quitting = time.Now().Add(quitTime)
		}

		cmd := m.progress.SetPercent(m.data.GetPercentage())
		return m, tea.Batch(tickCmd(), cmd)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:

		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:

		spinner, spinnerCMD := m.spinner.Update(msg)
		m.spinner = spinner

		stopwatch, stopWatchCMD := m.stopwatch.Update(msg)
		m.stopwatch = stopwatch

		return m, tea.Batch(spinnerCMD, stopWatchCMD)
	}
}

func (m *model) View() string {
	apply := lipgloss.NewStyle().MaxHeight(m.height).Render

	pad := strings.Repeat(" ", padding)

	loading := "    "
	if m.data.StopDueToFailures() {
		loading = "🛑  "
	} else if m.data.GetFailures() > 0 {
		loading = "🛜  "
	} else if !m.quitting.IsZero() {
		loading = "🏁  "
	} else if m.data.IsLoading() {
		loading = m.spinner.View() + "  "
	}

	truncate := func(s string) string {
		if len(s) > m.width {
			return s[:m.width-3] + "..."
		}
		return s
	}

	progress := truncate(strings.Join([]string{
		loading,
		m.data.GetText(),
	}, pad))
	elapsed := truncate(helpStyle("Elapsed " + m.stopwatch.View()))

	return apply(strings.Join([]string{
		fmt.Sprintln(),
		pad + progress,
		pad + m.progress.View(),
		fmt.Sprintln(),
		pad + elapsed,
		pad + helpStyle(fmt.Sprintf("Press any key to quit, logs can be found at %s", m.logFileName)),
		fmt.Sprintln(),
	}, fmt.Sprintln()))
}

func tickCmd() tea.Cmd {
	return tea.Tick(40*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
