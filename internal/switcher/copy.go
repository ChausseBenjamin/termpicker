package switcher

import (
	"github.com/ChausseBenjamin/termpicker/internal/colors"
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
