package quit

import (
	"github.com/ChausseBenjamin/termpicker/internal/ui"
	tea "github.com/charmbracelet/bubbletea/v2"
)

const byeMsg = "Goodbye!\n"

type Model struct{}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m Model) View() string { return ui.Style().Quit.Render(byeMsg) }
