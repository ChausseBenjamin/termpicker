package switcher

import (
	"fmt"
	"log/slog"

	"github.com/ChausseBenjamin/termpicker/internal/colors"
	"github.com/ChausseBenjamin/termpicker/internal/picker"
	"github.com/ChausseBenjamin/termpicker/internal/preview"
	"github.com/ChausseBenjamin/termpicker/internal/quit"
	"github.com/ChausseBenjamin/termpicker/internal/util"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	active   int
	pickers  []picker.Model
	preview  preview.Model
	help     help.Model
	fullHelp bool // When false, only show help for the switcher (not children)
}

func New(pickers []picker.Model) Model {
	return Model{
		active:   0,
		pickers:  pickers,
		preview:  *preview.New(colors.Hex(pickers[0].GetColor())),
		help:     help.New(),
		fullHelp: false,
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
	tabs := "|"
	for i, p := range m.pickers {
		if i == m.active {
			tabs += fmt.Sprintf(">%s<|", p.Title())
		} else {
			tabs += fmt.Sprintf(" %s |", p.Title())
		}
	}

	pickerView := m.pickers[m.active].View()
	w := lipgloss.Width(pickerView)

	m.help.Styles.ShortKey.Width(w)
	var helpstr string
	if m.fullHelp {
		helpstr = m.help.FullHelpView(m.AllKeys())
	} else {
		// This is a hack since the current view has too many keys
		// and the horizontal "ShortHelpView" gets too wide.
		// "FullHelpView" seperates keys by columns (and we only show the first).
		// helpstr = m.help.FullHelpView([][]key.Binding{m.AllKeys()[0]})
		helpstr = m.help.FullHelpView(shortKeys())
	}

	return fmt.Sprintf("%s\n%s\n%s\n%v",
		tabs,
		pickerView,
		m.preview.View(),
		helpstr,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keys := newKeybinds()
	cmds := []tea.Cmd{}
	slog.Info("Received tea.Msg", "tea_msg", msg, "type", fmt.Sprintf("%T", msg))
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

		case key.Matches(msg, keys.cpHex):
			util.Copy(colors.Hex(m.pickers[m.active].GetColor()))

		case key.Matches(msg, keys.cpRgb):
			pc := m.pickers[m.active].GetColor().ToPrecise()
			rgb := colors.RGB{}.FromPrecise(pc).(colors.RGB)
			util.Copy(rgb.String())

		case key.Matches(msg, keys.cpHsl):
			pc := m.pickers[m.active].GetColor().ToPrecise()
			hsl := colors.HSL{}.FromPrecise(pc).(colors.HSL)
			util.Copy(hsl.String())

		case key.Matches(msg, keys.cpCmyk):
			pc := m.pickers[m.active].GetColor().ToPrecise()
			cmyk := colors.CMYK{}.FromPrecise(pc).(colors.CMYK)
			util.Copy(cmyk.String())

		case key.Matches(msg, keys.help):
			m.fullHelp = !m.fullHelp

		case key.Matches(msg, keys.quit):
			return quit.Model{}, tea.Quit

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
	default:
		// fmt.Printf("\nmsg: %T\n", msg)
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
