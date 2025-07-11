package toosmall

import (
	"fmt"
	"strings"

	"github.com/ChausseBenjamin/termpicker/internal/ui"
	tea "github.com/charmbracelet/bubbletea/v2"
	lg "github.com/charmbracelet/lipgloss/v2"
)

type Model struct {
	parent tea.Model
	w      int
}

func New(parent tea.Model) Model {
	return Model{parent: parent}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	return ui.Style().Notice.Align(lg.Center).
		Faint(false).Foreground(lg.Color("#e06060")).Bold(true).
		Render(
			fmt.Sprintf("%s\nTerminal too small!\nPlease resize...",
				strings.Repeat(" ", m.w),
			),
		)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+z":
			return m, tea.Suspend
		}
	case tea.WindowSizeMsg:
		if msg.Height >= m.parentHeight() && msg.Width >= m.parentWidth() {
			return m.parent.Update(msg)
		}
		m.w = msg.Width
		return m, nil
	}
	return m, nil
}

func (m Model) parentWidth() int {
	// In v2, we need to handle the ViewModel interface
	if vm, ok := m.parent.(tea.ViewModel); ok {
		return lg.Width(vm.View())
	}
	return 0
}

func (m Model) parentHeight() int {
	// In v2, we need to handle the ViewModel interface
	if vm, ok := m.parent.(tea.ViewModel); ok {
		return lg.Height(vm.View())
	}
	return 0
}
