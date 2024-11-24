package quit

import tea "github.com/charmbracelet/bubbletea"

const byeMsg = "Goodbye!\n"

type Model struct{}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m Model) View() string { return byeMsg }
