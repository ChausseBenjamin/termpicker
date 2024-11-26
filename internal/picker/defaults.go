package picker

import (
	"github.com/ChausseBenjamin/termpicker/internal/progress"
	"github.com/ChausseBenjamin/termpicker/internal/slider"
)

func RGB() *Model {
	r := slider.New('R', 255, progress.WithGradient("#660000", "#ff0000"))
	g := slider.New('G', 255, progress.WithGradient("#006600", "#00ff00"))
	b := slider.New('B', 255, progress.WithGradient("#000066", "#0000ff"))
	rgb := New([]slider.Model{r, g, b}, "RGB")
	return rgb
}

func CMYK() *Model {
	c := slider.New('C', 100, progress.WithGradient("#006666", "#00ffff"))
	m := slider.New('M', 100, progress.WithGradient("#660066", "#ff00ff"))
	y := slider.New('Y', 100, progress.WithGradient("#666600", "#ffff00"))
	k := slider.New('K', 100, progress.WithSolidFill("#000000"))
	cmyk := New([]slider.Model{c, m, y, k}, "CMYK")
	return cmyk
}

func HSL() *Model {
	h := slider.New('H', 360, progress.WithDefaultGradient())
	s := slider.New('S', 100, progress.WithGradient("#a68e59", "#ffae00"))
	l := slider.New('L', 100, progress.WithGradient("#222222", "#ffffff"))
	hsl := New([]slider.Model{h, s, l}, "HSL")
	return hsl
}
