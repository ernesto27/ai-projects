package tui

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/crypto/ssh/terminal"
)

func GetTerminalWidth() int {
	width, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil || width == 0 {
		return 80
	}
	return width
}

var (
	// Compact header
	HeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.AdaptiveColor{Light: "#6366F1", Dark: "#818CF8"}).
		Padding(0, 2).
		Align(lipgloss.Center)

	// Compact info text
	InfoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#9CA3AF"}).
		Italic(true)

	// Compact prompt section
	PromptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#10B981")).
		Bold(true).
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color("#10B981")).
		PaddingLeft(1)

	// Elegant response section
	ResponseStyle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#1F2937", Dark: "#F9FAFB"}).
		Background(lipgloss.AdaptiveColor{Light: "#F8FAFC", Dark: "#111827"}).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#E5E7EB", Dark: "#374151"}).
		Padding(2, 3).
		MarginBottom(1).
		Width(GetTerminalWidth() - 4)

	// Error styling
	ErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#EF4444")).
		Bold(true).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#DC2626")).
		Padding(1, 2).
		MarginBottom(1).
		Width(GetTerminalWidth() - 4)

	// Loading style
	LoadingStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6366F1")).
		Bold(true).
		MarginTop(1).
		MarginBottom(1)

	// Divider style
	DividerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#E5E7EB", Dark: "#374151"}).
		MarginTop(1).
		MarginBottom(1)
)