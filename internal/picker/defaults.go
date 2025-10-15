package picker

import (
	"github.com/ChausseBenjamin/termpicker/internal/slider"
	"github.com/ChausseBenjamin/termpicker/internal/ui"
)

func RGB() *Model {
	return New(
		[]slider.Model{
			slider.New('R', 255, ui.Style().Sliders.R...),
			slider.New('G', 255, ui.Style().Sliders.G...),
			slider.New('B', 255, ui.Style().Sliders.B...),
		}, "RGB")
}

func CMYK() *Model {
	return New(
		[]slider.Model{
			slider.New('C', 100, ui.Style().Sliders.C...),
			slider.New('M', 100, ui.Style().Sliders.M...),
			slider.New('Y', 100, ui.Style().Sliders.Y...),
			slider.New('K', 100, ui.Style().Sliders.K...),
		}, "CMYK")
}

func HSL() *Model {
	return New(
		[]slider.Model{
			slider.New('H', 360, ui.Style().Sliders.H...),
			slider.New('S', 100, ui.Style().Sliders.S...),
			slider.New('L', 100, ui.Style().Sliders.L...),
		}, "HSL")
}

func OKLCH() *Model {
	return New(
		[]slider.Model{
			slider.New('H', 360, ui.Style().Sliders.OH...),  // 0-360 as-is
			slider.New('C', 500, ui.Style().Sliders.OC...),  // 0-0.5 scaled to 0-500
			slider.New('L', 1000, ui.Style().Sliders.OL...), // 0-1 scaled to 0-1000
		}, "OKLCH")
}
