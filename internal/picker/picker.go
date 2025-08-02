package picker

import (
	"fmt"
	"strings"

	"github.com/ChausseBenjamin/termpicker/internal/colors"
	"github.com/ChausseBenjamin/termpicker/internal/slider"
	"github.com/ChausseBenjamin/termpicker/internal/ui"
	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
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
	case "OKLCH":
		return colors.OKLCH{
			L: float64(m.sliders[0].Val()) / 1000.0, // Scale back from 0-1000 to 0-1
			C: float64(m.sliders[1].Val()) / 1000.0, // Scale back from 0-500 to 0-0.5
			H: float64(m.sliders[2].Val()),          // Use as-is 0-360
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
	case "OKLCH":
		oklch := colors.OKLCH{}.FromPrecise(p).(colors.OKLCH)
		m.sliders[0].Set(int(oklch.L * 1000.0)) // Scale 0-1 to 0-1000
		m.sliders[1].Set(int(oklch.C * 1000.0)) // Scale 0-0.5 to 0-500
		m.sliders[2].Set(int(oklch.H))          // Use as-is 0-360
	}
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{}
	for _, s := range m.sliders {
		cmds = append(cmds, s.Init())
	}
	return tea.Batch(cmds...)
}

func ViewSlider(active bool, s slider.Model) string {
	if active {
		return fmt.Sprintf("%s %s",
			ui.Style().PickerCursor.Render(ui.PickerSelRune),
			s.View(),
		)
	}
	return fmt.Sprintf("  %s",
		s.View(),
	)
}

func (m Model) View() string {
	sliderList := make([]string, len(m.sliders))
	for i, slider := range m.sliders {
		sliderList[i] = ViewSlider(i == m.active, slider)
	}
	return strings.Join(sliderList, "\n")
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
