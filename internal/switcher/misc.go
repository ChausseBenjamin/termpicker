package switcher

import (
	"log/slog"

	"github.com/ChausseBenjamin/termpicker/internal/colors"
	"github.com/ChausseBenjamin/termpicker/internal/parse"
	"github.com/ChausseBenjamin/termpicker/internal/util"
)

const (
	okCpMsg = "Copied %s to clipboard as %s"
)

func (m Model) copyColor(format string) string {
	pc := m.pickers[m.active].GetColor().ToPrecise()
	switch format {
	case cpHex:
		return util.Copy(colors.Hex(m.pickers[m.active].GetColor()))
	case cpRGB:
		rgb := colors.RGB{}.FromPrecise(pc).(colors.RGB)
		return util.Copy(rgb.String())
	case cpHSL:
		hsl := colors.HSL{}.FromPrecise(pc).(colors.HSL)
		return util.Copy(hsl.String())
	case cpCMYK:
		cmyk := colors.CMYK{}.FromPrecise(pc).(colors.CMYK)
		return util.Copy(cmyk.String())
	default:
		return "Copy format not supported"
	}
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
		}
		return "Color set to " + colorStr
	}
}
