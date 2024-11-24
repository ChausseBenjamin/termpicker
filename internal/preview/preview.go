package preview

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const runeBlock = "█"

type Model struct {
	size int // height of the square in rows
	hex  string
}

func (m *Model) Color(hex string) { m.hex = hex }

func New(hex string) *Model {
	return &Model{
		size: 5,
		hex:  hex,
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m Model) View() string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(m.hex))
	// size is doubled since terminal cells are 2:1 (h:w)
	oneRow := strings.Repeat(runeBlock, m.size*2)
	block := strings.Repeat(oneRow+"\n", m.size)
	return style.Render(block)
}
