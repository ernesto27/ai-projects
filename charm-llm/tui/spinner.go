package tui

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SpinnerModel struct {
	spinner spinner.Model
	message string
	ctx     context.Context
	cancel  context.CancelFunc
}

type spinnerFinishedMsg struct{}

func NewSpinnerModel(message string) *SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	ctx, cancel := context.WithCancel(context.Background())

	return &SpinnerModel{
		spinner: s,
		message: message,
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (m *SpinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			m.cancel()
			return m, tea.Quit
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case spinnerFinishedMsg:
		return m, tea.Quit
	}
	return m, nil
}

func (m *SpinnerModel) View() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginLeft(2)

	return style.Render(m.spinner.View() + " " + m.message)
}

func (m *SpinnerModel) Stop() {
	m.cancel()
}

func (m *SpinnerModel) Context() context.Context {
	return m.ctx
}

// ShowSpinnerWhile shows a spinner for the duration of the provided function
func ShowSpinnerWhile(message string, fn func(ctx context.Context) error) error {
	model := NewSpinnerModel(message)

	program := tea.NewProgram(model)

	// Start the spinner in a goroutine
	go func() {
		if _, err := program.Run(); err != nil {
			fmt.Println("Error running spinner:", err)
		}
	}()

	// Execute the function
	err := fn(model.Context())

	// Stop the spinner
	program.Send(spinnerFinishedMsg{})
	time.Sleep(100 * time.Millisecond) // Small delay to ensure clean exit
	program.Quit()

	return err
}
