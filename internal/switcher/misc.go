package switcher

import (
	"log/slog"

	"github.com/ChausseBenjamin/termpicker/internal/colors"
	"github.com/ChausseBenjamin/termpicker/internal/parse"
	"github.com/ChausseBenjamin/termpicker/internal/util"
	tea "github.com/charmbracelet/bubbletea/v2"
)

func (m Model) copyColor(format string) tea.Cmd {
	pc := m.pickers[m.active].GetColor().ToPrecise()
	var colorStr string

	switch format {
	case cpHex:
		colorStr = colors.Hex(m.pickers[m.active].GetColor())
	case cpRGB:
		rgb := colors.RGB{}.FromPrecise(pc).(colors.RGB)
		colorStr = rgb.String()
	case cpHSL:
		hsl := colors.HSL{}.FromPrecise(pc).(colors.HSL)
		colorStr = hsl.String()
	case cpCMYK:
		cmyk := colors.CMYK{}.FromPrecise(pc).(colors.CMYK)
		colorStr = cmyk.String()
	case cpOKLCH:
		oklch := colors.OKLCH{}.FromPrecise(pc).(colors.OKLCH)
		colorStr = oklch.String()
	case cpEscFG:
		colorStr = colors.EscapedSeq(m.pickers[m.active].GetColor(), true)
	case cpEscBG:
		colorStr = colors.EscapedSeq(m.pickers[m.active].GetColor(), false)
	default:
		return func() tea.Msg {
			return util.ClipboardResultMsg{
				Success: false,
				Message: "Copy format not supported",
			}
		}
	}

	return util.SmartCopyToClipboard(colorStr)
}

func (m *Model) SetColorFromText(colorStr string) string {
	color, err := parse.Color(colorStr)
	if err != nil {
		slog.Error("Failed to parse color", util.ErrKey, err)
		return err.Error()
	} else {
		pc := color.ToPrecise()
		switch color.(type) {
		case colors.RGB:
			m.UpdatePicker(IndexRgb, pc)
			m.SetActive(IndexRgb)
		case colors.CMYK:
			m.UpdatePicker(IndexCmyk, pc)
			m.SetActive(IndexCmyk)
		case colors.HSL:
			m.UpdatePicker(IndexHsl, pc)
			m.SetActive(IndexHsl)
		case colors.OKLCH:
			m.UpdatePicker(IndexOklch, pc)
			m.SetActive(IndexOklch)
		}
		return "Color set to " + colorStr
	}
}
