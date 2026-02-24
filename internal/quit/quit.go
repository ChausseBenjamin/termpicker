package quit

import (
	tea "charm.land/bubbletea/v2"
	"github.com/ChausseBenjamin/termpicker/internal/ui"
)

const byeMsg = "Goodbye!\n"

type Model struct{}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m Model) View() tea.View { return tea.NewView(ui.Style().Quit.Render(byeMsg)) }
