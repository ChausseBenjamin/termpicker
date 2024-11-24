package colors

import (
	"fmt"
	"math"
)

type HSL struct {
	H int // 0-360
	S int // 0-100
	L int // 0-100
}

func (h HSL) String() string {
	return fmt.Sprintf("hsl(%d, %d%%, %d%%)", h.H, h.S, h.L)
}

func (h HSL) ToPrecise() PreciseColor {
	// Normalize H, S, L
	hue := float64(h.H) / 360.0
	sat := float64(h.S) / 100.0
	light := float64(h.L) / 100.0

	var r, g, b float64

	if sat == 0 {
		// Achromatic case
		r, g, b = light, light, light
	} else {
		var q float64
		if light < 0.5 {
			q = light * (1 + sat)
		} else {
			q = light + sat - (light * sat)
		}
		p := 2*light - q
		r = hueToRGB(p, q, hue+1.0/3.0)
		g = hueToRGB(p, q, hue)
		b = hueToRGB(p, q, hue-1.0/3.0)
	}

	return PreciseColor{R: r, G: g, B: b}
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}

func (h HSL) FromPrecise(p PreciseColor) ColorSpace {
	r := p.R
	g := p.G
	b := p.B

	max := math.Max(math.Max(r, g), b)
	min := math.Min(math.Min(r, g), b)
	delta := max - min

	light := (max + min) / 2
	var sat, hue float64

	if delta == 0 {
		// Achromatic case
		hue, sat = 0, 0
	} else {
		if light < 0.5 {
			sat = delta / (max + min)
		} else {
			sat = delta / (2 - max - min)
		}

		switch max {
		case r:
			hue = (g-b)/delta + (6 * boolToFloat64(g < b))
		case g:
			hue = (b-r)/delta + 2
		case b:
			hue = (r-g)/delta + 4
		}
		hue /= 6
	}

	return HSL{
		H: int(math.Round(hue * 360)),
		S: int(math.Round(sat * 100)),
		L: int(math.Round(light * 100)),
	}
}

func boolToFloat64(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
