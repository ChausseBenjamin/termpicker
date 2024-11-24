package switcher

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbletea-app-template/internal/colors"
	"github.com/charmbracelet/bubbletea-app-template/internal/picker"
	"github.com/charmbracelet/bubbletea-app-template/internal/preview"
	"github.com/charmbracelet/bubbletea-app-template/internal/quit"
)

type Model struct {
	active  int
	pickers []picker.Model
	preview preview.Model
}

func New(pickers []picker.Model) Model {
	return Model{
		active:  0,
		pickers: pickers,
		preview: *preview.New(colors.Hex(pickers[0].GetColor())),
	}
}

func (m Model) fixSel(val int) int {
	size := len(m.pickers)
	return (val%size + size) % size
}

func (m *Model) Next() int {
	m.active = m.fixSel(m.active + 1)
	return m.active
}

func (m *Model) Prev() int {
	m.active = m.fixSel(m.active - 1)
	return m.active
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{}
	for _, p := range m.pickers {
		cmds = append(cmds, p.Init())
	}
	return tea.Batch(cmds...)
}

func (m Model) View() string {
	v := "|"
	for i, p := range m.pickers {
		if i == m.active {
			v += fmt.Sprintf(">%s<|", p.Title())
		} else {
			v += fmt.Sprintf(" %s |", p.Title())
		}
	}
	return fmt.Sprintf("%s\n%s\n%s", v, m.pickers[m.active].View(), m.preview.View())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keys := newKeybinds()
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.next):
			cs := m.pickers[m.active].GetColor()
			m.Next()
			m.pickers[m.active].SetColor(cs)
		case key.Matches(msg, keys.prev):
			cs := m.pickers[m.active].GetColor()
			m.Prev()
			m.pickers[m.active].SetColor(cs)
		case key.Matches(msg, keys.quit):
			return quit.Model{}, tea.Quit
			// return m, tea.Quit
		default:
			// Update the picker
			newActive, cmd := m.pickers[m.active].Update(msg)
			m.pickers[m.active] = newActive.(picker.Model)
			cmds = append(cmds, cmd)
			// Update the preview
			newPreview := preview.New(colors.Hex(m.pickers[m.active].GetColor()))
			m.preview = *newPreview

			return m, tea.Batch(cmds...)
		}
	}
	for i, p := range m.pickers {
		newActive, cmd := p.Update(msg)
		m.pickers[i] = newActive.(picker.Model)
		cmds = append(cmds, cmd)
	}
	// Update the preview
	newPreview := preview.New(colors.Hex(m.pickers[m.active].GetColor()))
	m.preview = *newPreview
	return m, tea.Batch(cmds...)
}
