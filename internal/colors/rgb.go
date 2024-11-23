package colors

import "math"

type RGB struct {
	R int // 0-255
	G int // 0-255
	B int // 0-255
}

func (c RGB) ToPrecise() PreciseColor {
	return PreciseColor{
		R: float64(c.R) / 255,
		G: float64(c.G) / 255,
		B: float64(c.B) / 255,
	}
}

func (c RGB) FromPrecise(p PreciseColor) ColorSpace {
	return RGB{
		R: int(math.Round(p.R * 255)),
		G: int(math.Round(p.G * 255)),
		B: int(math.Round(p.B * 255)),
	}
}
