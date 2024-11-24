package preview

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	runeBlock     = "█"
	defaultHeight = 5
	defaultWidth  = 78
)

type Model struct {
	height int
	width  int
	hex    string
}

func (m *Model) SetColor(hex string) { m.hex = hex }

func (m *Model) SetHeight(size int) { m.height = size }

func (m *Model) SetWidth(size int) { m.width = size }

func New(hex string) *Model {
	return &Model{
		height: defaultHeight,
		width:  defaultWidth,
		hex:    hex,
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m Model) View() string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(m.hex))
	oneRow := strings.Repeat(runeBlock, m.width)
	block := strings.Repeat(oneRow+"\n", m.height)
	return style.Render(block)
}
