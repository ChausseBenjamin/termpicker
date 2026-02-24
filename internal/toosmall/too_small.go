package toosmall

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	lg "charm.land/lipgloss/v2"
	"github.com/ChausseBenjamin/termpicker/internal/ui"
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

func (m Model) View() tea.View {
	return tea.NewView(ui.Style().Notice.Align(lg.Center).
		Faint(false).Foreground(lg.Color("#e06060")).Bold(true).
		Render(
			fmt.Sprintf("%s\nTerminal too small!\nPlease resize...",
				strings.Repeat(" ", m.w),
			),
		))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
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
	// In v2, we get the view directly from the model
	view := m.parent.View()
	return lg.Width(view.Content)
}

func (m Model) parentHeight() int {
	// In v2, we get the view directly from the model
	view := m.parent.View()
	return lg.Height(view.Content)
}
