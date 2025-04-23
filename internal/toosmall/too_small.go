package toosmall

import (
	"fmt"
	"strings"

	"github.com/ChausseBenjamin/termpicker/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
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
	return lg.Width(m.parent.View())
}

func (m Model) parentHeight() int {
	return lg.Height(m.parent.View())
}
