package preview

import (
	"log/slog"
	"strings"

	"github.com/ChausseBenjamin/termpicker/internal/colors"
	"github.com/ChausseBenjamin/termpicker/internal/util"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	runeBlock     = " "
	defaultHeight = 5
	defaultWidth  = 78
)

type ColorMsg colors.ColorSpace

type Model struct {
	height int
	width  int
	hex    string
	cfg    Config
}

type Config struct {
	PreviewStr string
	PreviewBg  string
	PreviewFg  string
}

func (m *Model) SetColor(hex string) { m.hex = hex }

func (m *Model) SetHeight(size int) { m.height = size }

func (m *Model) SetWidth(size int) { m.width = size }

func New(hex string) *Model {
	return &Model{
		height: defaultHeight,
		width:  defaultWidth,
		hex:    hex,
		cfg: Config{
			PreviewStr: util.DefaultPreviewText,
			PreviewFg:  "#ffffff",
			PreviewBg:  "#000000",
		},
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ColorMsg:
		m.hex = colors.Hex(msg)
	case Config:
		m.cfg = msg
		slog.Info("Updating",
			slog.Group("Model",
				slog.Int("height", m.height),
				slog.Int("width", m.width),
				slog.String("hex", m.hex),
				slog.Group("Config",
					"PreviewString", m.cfg.PreviewStr,
					"Foreground", m.cfg.PreviewFg,
					"Background", m.cfg.PreviewBg,
				),
			),
		)
	}
	return m, nil
}

func (m Model) View() string {
	normStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(m.hex)).
		Foreground(lipgloss.Color(m.cfg.PreviewFg)).
		Align(lipgloss.Center).
		Width(m.width)
	var buffer = 0
	prevRows := ""
	if m.cfg.PreviewStr != "" {
		slog.Info("Detected Custom preview for rendering")
		// When previewing termpicker in relation with other colors,
		// The inverted style will use the target color as a foreground
		// (text) with a predefined comparison color.
		invStyle := normStyle.
			Background(lipgloss.Color(m.cfg.PreviewBg)).
			Foreground(lipgloss.Color(m.hex)).
			Align(lipgloss.Center).
			Width(m.width)
		buffer = 2
		prevRows = invStyle.Render(m.cfg.PreviewStr) + "\n" +
			normStyle.Render(m.cfg.PreviewStr) + "\n"
	}

	oneRow := strings.Repeat(runeBlock, m.width) + "\n"
	block := prevRows + normStyle.Render(strings.Repeat(oneRow, m.height-buffer))
	return block
}
