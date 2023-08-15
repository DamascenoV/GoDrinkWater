package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	title                     = "GoDrinkWater - "
	windowWidth, windowHeight int
	quitKeys                  = key.NewBinding(key.WithKeys("esc", "q"))
)

type model struct {
	duration   time.Duration
	timePassed time.Duration
	start      time.Time
	timer      timer.Model
	progress   progress.Model
	quit       bool
}

func (m model) Init() tea.Cmd {
	return m.timer.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case timer.TickMsg:
		var commands []tea.Cmd
		var cmd tea.Cmd
		m.timePassed += m.timer.Interval
		countdown := m.timePassed.Milliseconds() * 100 / m.duration.Milliseconds()
		commands = append(commands, m.progress.SetPercent(float64(countdown)/100))
		m.timer, cmd = m.timer.Update(msg)
		commands = append(commands, cmd)

		return m, tea.Batch(commands...)

	case progress.FrameMsg:
		modelProgress, cmd := m.progress.Update(msg)
		m.progress = modelProgress.(progress.Model)

		return m, cmd

	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)

		return m, cmd

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - 2*2 - 4
		windowHeight = msg.Height
		windowWidth = msg.Width
		m.progress.Width = 80

		return m, nil

	case timer.TimeoutMsg:
		m.quit = true

		return m, tea.Quit

	case tea.KeyMsg:
		if key.Matches(msg, quitKeys) {
			m.quit = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.quit {
		return "\n"
	}

	view := lipgloss.NewStyle().Bold(true).Render(title)
	view += lipgloss.NewStyle().Bold(true).Render(m.timer.View() + "\n" + m.progress.View())

	viewWidth, viewHeight := lipgloss.Size(view)

	width := (windowWidth - viewWidth) / 2
	height := (windowHeight - viewHeight) / 2

	return lipgloss.NewStyle().Margin(height, width).Render(view)
}

func Run(duration time.Duration, interval time.Duration, opts []tea.ProgramOption) (tea.Model, error) {
	p, err := tea.NewProgram(model{
		duration: duration,
		timer:    timer.NewWithInterval(duration, interval),
		progress: progress.New(progress.WithDefaultGradient()),
		start:    time.Now(),
	}, opts...).Run()
	if err != nil {
		return nil, err
	}

	return p, nil
}
