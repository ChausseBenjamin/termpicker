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
	"github.com/ChausseBenjamin/termpicker/internal/toosmall"
	"github.com/ChausseBenjamin/termpicker/internal/ui"
	"github.com/ChausseBenjamin/termpicker/internal/util"
	"github.com/charmbracelet/bubbles/v2/help"
	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/textinput"
	tea "github.com/charmbracelet/bubbletea/v2"
	lg "github.com/charmbracelet/lipgloss/v2"
)

const (
	IndexRgb int = iota
	IndexHsl
	IndexCmyk
)

type Model struct {
	active   int
	pickers  []picker.Model
	prev     preview.Model
	help     help.Model
	input    textinput.Model
	notice   notices.Model
	fullHelp bool // When false, only show help for the switcher (not children)
}

func New() Model {
	pickers := []picker.Model{ // Order MUST match the Index* constants
		*picker.RGB(),
		*picker.HSL(),
		*picker.CMYK(),
	}

	input := textinput.New()
	input.Prompt = ui.PromptPrefix
	input.Placeholder = ui.PromptPlaceholder

	return Model{
		active:   0,
		pickers:  pickers,
		prev:     *preview.New(colors.Hex(pickers[0].GetColor())),
		help:     help.New(),
		input:    input,
		notice:   notices.New(),
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

func (m *Model) UpdatePreview(cfg preview.Config) {
	newPrev, _ := m.prev.Update(cfg)
	m.prev = newPrev.(preview.Model)
}

func (m *Model) NewNotice(msg string) tea.Cmd {
	return m.notice.New(msg)
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{}

	// Initialize clipboard system with terminal version query
	cmds = append(cmds, util.InitClipboard())

	// The NoticeExpiryMsg is never sent to bubbletea by a tea.Cmd for the initial notices
	// That's why we need to manually reset them here. Otherwise, they would never expire.
	for k := range m.notice.Notices {
		cmds = append(cmds, m.notice.Reset(k))
	}

	for _, p := range m.pickers {
		cmds = append(cmds, p.Init())
	}
	return tea.Batch(cmds...)
}

func (m Model) View() string {
	tabs := make([]string, len(m.pickers))
	for i, p := range m.pickers {
		if i == m.active {
			tabs[i] = ui.Style().TabSel.Render(p.Title())
		} else {
			tabs[i] = ui.Style().TabNorm.Render(p.Title())
		}
	}
	tabStr := strings.Join([]string{
		ui.Style().TabGeom.Render(ui.TabSepLeft),
		strings.Join(tabs, ui.Style().TabGeom.Render(ui.TabSepMid)),
		ui.Style().TabGeom.Render(ui.TabSepRight),
	}, " ")

	pickerStr := m.pickers[m.active].View()
	w := lg.Width(pickerStr)

	m.prev.SetWidth(w)
	previewStr := m.prev.View()

	m.help.Styles.ShortKey.Width(w)

	var helpStr string
	m.help.Width = w
	if m.fullHelp {
		helpStr = m.help.FullHelpView(m.AllKeys())
	} else {
		// This is a hack since the current view has too many keys
		// and the horizontal "ShortHelpView" gets too wide.
		// "FullHelpView" seperates keys by columns (and we only show the first).
		// helpStr = m.help.FullHelpView([][]key.Binding{m.AllKeys()[0]})
		helpStr = m.help.FullHelpView(shortKeys())
	}

	var inputStr string
	if m.input.Focused() {
		m.input.SetWidth(w - lg.Width(ui.PromptPrefix) - 1)
		inputStr = ui.Style().Boxed.Render(m.input.View())
	}

	mainArea := ui.Style().Boxed.Render(strings.Join([]string{
		pickerStr,
		previewStr,
		helpStr,
	}, "\n"))

	return strings.Join(
		[]string{
			tabStr,
			mainArea,
			inputStr,
			m.notice.View(),
		}, "\n")
}

func (m Model) Fits(s tea.WindowSizeMsg) bool {
	return s.Width >= lg.Width(m.View()) && s.Height >= lg.Height(m.View())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keys := newKeybinds()
	cmds := []tea.Cmd{}
	slog.Debug("Received tea.Msg", "tea_msg", msg, "type", fmt.Sprintf("%T", msg))
	switch msg := msg.(type) {
	case notices.NoticeExpiryMsg:
		newNotices, cmd := m.notice.Update(msg)
		m.notice = newNotices.(notices.Model)
		cmds = append(cmds, cmd)

	case tea.WindowSizeMsg:
		if !m.Fits(msg) {
			m.notice.Notices = make(map[string]string)
			smol := toosmall.New(m)
			return smol.Update(msg)
		}

	case preview.Config:
		newPreview, cmd := m.prev.Update(msg)
		cmds = append(cmds, cmd)
		m.prev = newPreview.(preview.Model)

	case tea.KeyMsg:

		if m.input.Focused() && msg.String() != "ctrl+c" {
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
			cmd := m.copyColor(msg.String())
			cmds = append(cmds, cmd)

		case key.Matches(msg, keys.help):
			m.fullHelp = !m.fullHelp

		case key.Matches(msg, keys.insert):
			cmd := m.input.Focus()
			cmds = append(cmds, cmd)

		case key.Matches(msg, keys.suspend):
			return m, tea.Suspend

		case key.Matches(msg, keys.quit):
			return quit.Model{}, tea.Quit

		default: // Update the picker
			newActive, cmd := m.pickers[m.active].Update(msg)
			m.pickers[m.active] = newActive.(picker.Model)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}

	case tea.TerminalVersionMsg:
		// Handle terminal version for clipboard decisions
		util.HandleTerminalVersion(string(msg))

	case util.ClipboardResultMsg:
		// Handle clipboard operation results
		if msg.Success {
			// Show success message
			cmd := m.notice.New(msg.Message)
			cmds = append(cmds, cmd)

			// Show warning if fallback was used
			if msg.UsedFallback {
				warningCmd := m.notice.New(util.CopyWarningMessage(msg.FallbackReason))
				cmds = append(cmds, warningCmd)
			}
		} else {
			// Show error message
			cmd := m.notice.New(msg.Message)
			cmds = append(cmds, cmd)
		}

	default:
	}
	for i, p := range m.pickers {
		newActive, cmd := p.Update(msg)
		m.pickers[i] = newActive.(picker.Model)
		cmds = append(cmds, cmd)
	}
	// Update the preview
	newPreview, cmd := m.prev.Update(m.pickers[m.active].GetColor())
	cmds = append(cmds, cmd)
	m.prev = newPreview.(preview.Model)
	return m, tea.Batch(cmds...)
}
