package slider

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	label    byte
	progress progress.Model
	max      int
	current  int
	mappings keybinds
}

func New(label byte, maxVal int, opts ...progress.Option) Model {
	slider := Model{
		label: label,
		progress: progress.New(
			progress.WithoutPercentage(),
		),
		max:      maxVal,
		current:  maxVal / 2,
		mappings: newKeybinds(),
	}
	for _, opt := range opts {
		opt(&slider.progress)
	}
	return slider
}

func (m Model) Title() string { return fmt.Sprintf("%c", m.label) }

func (m Model) Init() tea.Cmd {
	// Triggering a frame message Update here will force the progress bar to
	// render immediately. This is necessary because progress bars only render
	// when their state changes. Without this, there is a disrepancy between
	// the initial state of the progress bar and the initial state of the slider.

	// There's no sugar-coating it: This is a hack. But it works...
	_, cmd := m.Update(progress.FrameMsg{})
	return cmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	keys := newKeybinds()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.incRegular):
			m.IncPcnt(0.05)
		case key.Matches(msg, keys.decRegular):
			m.DecPcnt(0.05)
		case key.Matches(msg, keys.incPrecise):
			m.Inc(1)
		case key.Matches(msg, keys.decPrecise):
			m.Dec(1)
		}
		return m, m.progress.SetPercent(m.Pcnt())
	case progress.FrameMsg:
		if m.progress.Percent() != m.Pcnt() {
			cmds = append(cmds, m.progress.SetPercent(m.Pcnt()))
		}
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		cmds = append(cmds, cmd)
	default:
	}
	return m, tea.Batch(cmds...)
}

func (m Model) ViewValue(current int) string {
	return fmt.Sprintf("(%3d/%d)", current, m.max)
}

func (m Model) View() string {
	return fmt.Sprintf("%v: %v %v", m.Title(), m.progress.View(), m.ViewValue(m.current))
}
