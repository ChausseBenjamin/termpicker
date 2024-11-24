package picker

import (
	"fmt"

	"github.com/ChausseBenjamin/termpicker/internal/colors"
	"github.com/ChausseBenjamin/termpicker/internal/slider"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	activeRune = '>'
)

type Model struct {
	title   string
	active  int
	sliders []slider.Model
}

func (m *Model) Next() int {
	m.active = m.fixSel(m.active + 1)
	return m.active
}

func (m *Model) Prev() int {
	m.active = m.fixSel(m.active - 1)
	return m.active
}

func (m *Model) Sel(i int) int {
	m.active = m.fixSel(i)
	return m.active
}

func (m Model) fixSel(val int) int {
	size := len(m.sliders)
	return (val%size + size) % size
}

func New(sliders []slider.Model, title string) *Model {
	return &Model{
		title:   title,
		active:  0,
		sliders: sliders,
	}
}

func (m Model) Title() string {
	return m.title
}

func (m Model) GetColor() colors.ColorSpace {
	switch m.title {
	case "RGB":
		return colors.RGB{
			R: m.sliders[0].Val(),
			G: m.sliders[1].Val(),
			B: m.sliders[2].Val(),
		}
	case "CMYK":
		return colors.CMYK{
			C: m.sliders[0].Val(),
			M: m.sliders[1].Val(),
			Y: m.sliders[2].Val(),
			K: m.sliders[3].Val(),
		}
	case "HSL":
		return colors.HSL{
			H: m.sliders[0].Val(),
			S: m.sliders[1].Val(),
			L: m.sliders[2].Val(),
		}
	default: // Default to white if we don't know the color space
		return colors.RGB{
			R: 255,
			G: 255,
			B: 255,
		}
	}
}

func (m Model) SetColor(c colors.ColorSpace) {
	p := c.ToPrecise()
	switch m.title {
	case "RGB":
		rgb := colors.RGB{}.FromPrecise(p).(colors.RGB)
		m.sliders[0].Set(rgb.R)
		m.sliders[1].Set(rgb.G)
		m.sliders[2].Set(rgb.B)
	case "CMYK":
		cmyk := colors.CMYK{}.FromPrecise(p).(colors.CMYK)
		m.sliders[0].Set(cmyk.C)
		m.sliders[1].Set(cmyk.M)
		m.sliders[2].Set(cmyk.Y)
		m.sliders[3].Set(cmyk.K)
	case "HSL":
		hsl := colors.HSL{}.FromPrecise(p).(colors.HSL)
		m.sliders[0].Set(hsl.H)
		m.sliders[1].Set(hsl.S)
		m.sliders[2].Set(hsl.L)
	}
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{}
	for _, s := range m.sliders {
		cmds = append(cmds, s.Init())
	}
	return tea.Batch(cmds...)
}

func (m Model) View() string {
	var s string

	carriageReturn := ""
	for i, slider := range m.sliders {
		if i > 0 {
			carriageReturn = "\n"
		}
		if i == m.active {
			s += fmt.Sprintf("%v%c %s", carriageReturn, activeRune, slider.View())
		} else {
			s += fmt.Sprintf("%v  %s", carriageReturn, slider.View())
		}
	}
	return s
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keys := newKeybinds()
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.next):
			m.Next()
		case key.Matches(msg, keys.prev):
			m.Prev()
		default:
			newActive, cmd := m.sliders[m.active].Update(msg)
			m.sliders[m.active] = newActive.(slider.Model)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}
	}
	// Keys are only sent to the active sliders
	// However, other messages (ex: tick, resize) must be sent to all
	for i, s := range m.sliders {
		newSlider, cmd := s.Update(msg)
		m.sliders[i] = newSlider.(slider.Model)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}
