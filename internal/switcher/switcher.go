package switcher

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/ChausseBenjamin/termpicker/internal/colors"
	"github.com/ChausseBenjamin/termpicker/internal/notices"
	"github.com/ChausseBenjamin/termpicker/internal/picker"
	"github.com/ChausseBenjamin/termpicker/internal/preview"
	"github.com/ChausseBenjamin/termpicker/internal/quit"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	IndexRgb int = iota
	IndexCmyk
	IndexHsl
)

type Model struct {
	active   int
	pickers  []picker.Model
	preview  preview.Model
	help     help.Model
	input    textinput.Model
	notices  notices.Model
	fullHelp bool // When false, only show help for the switcher (not children)
}

func New() Model {
	pickers := []picker.Model{ // Order MUST match the Index* constants
		*picker.RGB(),
		*picker.CMYK(),
		*picker.HSL(),
	}
	return Model{
		active:   0,
		pickers:  pickers,
		preview:  *preview.New(colors.Hex(pickers[0].GetColor())),
		help:     help.New(),
		input:    textinput.New(),
		notices:  notices.New(),
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

func (m *Model) SetActive(i int) {
	m.active = m.fixSel(i)
}

func (m *Model) UpdatePicker(i int, c colors.ColorSpace) {
	m.pickers[i].SetColor(c)
}

func (m *Model) NewNotice(msg string) tea.Cmd {
	return m.notices.New(msg)
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{}

	// The NoticeExpiryMsg is never sent to bubbletea by a tea.Cmd for the initial notices
	// That's why we need to manually reset them here. Otherwise, they would never expire.
	for k := range m.notices.Notices {
		cmds = append(cmds, m.notices.Reset(k))
	}

	for _, p := range m.pickers {
		cmds = append(cmds, p.Init())
	}
	return tea.Batch(cmds...)
}

func (m Model) View() string {
	norm := lipgloss.NewStyle().Faint(true)
	bright := lipgloss.NewStyle().Faint(false)

	delims := [3]string{"[ ", " | ", "]"}
	for i, d := range delims {
		delims[i] = bright.Render(d)
	}

	var sections []string
	for i, p := range m.pickers {
		if i == m.active {
			sections = append(
				sections,
				bright.
					Underline(true).
					Bold(true).
					Render(p.Title()),
			)
		} else {
			sections = append(sections, norm.Render(p.Title()))
		}
	}
	tabs := "[ " + strings.Join(sections, " | ") + " ]"

	pickerView := m.pickers[m.active].View()
	boxStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true, true, false, true)
	w := lipgloss.Width(pickerView)
	pickerView = boxStyle.Render(pickerView)

	m.preview.SetWidth(w)
	boxStyle = boxStyle.Border(lipgloss.RoundedBorder(), false, true, false, true)
	previewStr := boxStyle.Render(m.preview.View())

	m.help.Styles.ShortKey.Width(w)
	boxStyle = boxStyle.Border(lipgloss.RoundedBorder(), false, true, true, true).Width(w)

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
	helpstr = boxStyle.Render(helpstr)

	inputStr := ""
	if m.input.Focused() {
		boxStyle = boxStyle.Border(lipgloss.RoundedBorder(), true, true, true, true).Width(w)
		inputStr = boxStyle.Render(m.input.View())
	}

	return fmt.Sprintf("%s\n%s\n%s\n%v\n%v\n%v",
		tabs,
		pickerView,
		previewStr,
		helpstr,
		inputStr,
		m.notices.View(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keys := newKeybinds()
	cmds := []tea.Cmd{}
	slog.Info("Received tea.Msg", "tea_msg", msg, "type", fmt.Sprintf("%T", msg))
	switch msg := msg.(type) {
	case notices.NoticeExpiryMsg:
		newNotices, cmd := m.notices.Update(msg)
		m.notices = newNotices.(notices.Model)
		cmds = append(cmds, cmd)

	case tea.KeyMsg:

		if m.input.Focused() {
			keys.esc.SetEnabled(true)
			keys.confirm.SetEnabled(true)
			if key.Matches(msg, keys.esc) {
				m.input.Blur()
			} else if key.Matches(msg, keys.confirm) {
				m.input.Blur()
				cmds = append(
					cmds,
					m.NewNotice(m.SetColorFromText(m.input.Value())),
					m.Init(), // Will force a slider update/animation
				)
			}
			newInput, cmd := m.input.Update(msg)
			m.input = newInput
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}

		switch {
		case key.Matches(msg, keys.next):
			cs := m.pickers[m.active].GetColor()
			m.Next()
			m.pickers[m.active].SetColor(cs)

		case key.Matches(msg, keys.prev):
			cs := m.pickers[m.active].GetColor()
			m.Prev()
			m.pickers[m.active].SetColor(cs)

		case key.Matches(msg, keys.copy):
			cmd := m.notices.New(m.copyColor(msg.String()))
			cmds = append(cmds, cmd)

		case key.Matches(msg, keys.help):
			m.fullHelp = !m.fullHelp

		case key.Matches(msg, keys.insert):
			cmd := m.input.Focus()
			cmds = append(cmds, cmd)

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

			newNotices, cmd := m.notices.Update(msg)
			m.notices = newNotices.(notices.Model)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}
	default:
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
