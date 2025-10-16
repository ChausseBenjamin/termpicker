package colors

import (
	"math"
	"testing"
)

const (
	PCmaxDelta     = 1e-2
	AssertTemplate = "Testing '%v'. Expected %v to become: %v Got: %v"
)

type equivalentColors struct {
	name  string
	pc    PreciseColor
	rgb   RGB
	cmyk  CMYK
	hsl   HSL
	oklch OKLCH
}

func pcDeltaOk(a, b PreciseColor) bool {
	return math.Abs(a.R-b.R) < PCmaxDelta &&
		math.Abs(a.G-b.G) < PCmaxDelta &&
		math.Abs(a.B-b.B) < PCmaxDelta
}

func getEquivalents() []equivalentColors {
	return []equivalentColors{
		// Black & White {{{
		{
			"Pure White",
			PreciseColor{1, 1, 1},
			RGB{255, 255, 255},
			CMYK{0, 0, 0, 0},
			HSL{0, 0, 100},
			OKLCH{1, 0, 0},
		},
		{
			"Pure Black",
			PreciseColor{0, 0, 0},
			RGB{0, 0, 0},
			CMYK{0, 0, 0, 100},
			HSL{0, 0, 0},
			OKLCH{0, 0, 0},
		},
		// }}}
		// Pure RGB {{{
		{
			"Red",
			PreciseColor{1, 0, 0},
			RGB{255, 0, 0},
			CMYK{0, 100, 100, 0},
			HSL{0, 100, 50},
			OKLCH{0.628, 0.258, 29.23},
		},
		{
			"Green",
			PreciseColor{0, 1, 0},
			RGB{0, 255, 0},
			CMYK{100, 0, 100, 0},
			HSL{120, 100, 50},
			OKLCH{0.866, 0.295, 142.51},
		},
		{
			"Blue",
			PreciseColor{0, 0, 1},
			RGB{0, 0, 255},
			CMYK{100, 100, 0, 0},
			HSL{240, 100, 50},
			OKLCH{0.452, 0.313, 264.06},
		},
		// }}}
		// Pure CMYK {{{
		{
			"Cyan",
			PreciseColor{0, 1, 1},
			RGB{0, 255, 255},
			CMYK{100, 0, 0, 0},
			HSL{180, 100, 50},
			OKLCH{0.905, 0.155, 194.80},
		},
		{
			"Magenta",
			PreciseColor{1, 0, 1},
			RGB{255, 0, 255},
			CMYK{0, 100, 0, 0},
			HSL{300, 100, 50},
			OKLCH{0.702, 0.323, 328.36},
		},
		{
			"Yellow",
			PreciseColor{1, 1, 0},
			RGB{255, 255, 0},
			CMYK{0, 0, 100, 0},
			HSL{60, 100, 50},
			OKLCH{0.968, 0.211, 109.78},
		},
		// note: Black is already tested
		// }}}
		// TODO: add less pure colors to test luminance and saturation better
	}
}

func TestToPreciseColor(t *testing.T) {
	for _, ce := range getEquivalents() {
		target := ce.pc
		for _, cs := range []ColorSpace{ce.rgb, ce.cmyk, ce.hsl, ce.oklch} {
			pc := cs.ToPrecise()
			if !pcDeltaOk(pc, target) {
				t.Errorf(AssertTemplate, ce.name, cs, target, pc)
			}
		}
	}
}

func TestToRgb(t *testing.T) {
	for _, ce := range getEquivalents() {
		target := ce.rgb
		for _, cs := range []ColorSpace{ce.pc, ce.cmyk, ce.hsl, ce.oklch} {
			rgb := RGB{}.FromPrecise(cs.ToPrecise()).(RGB)
			if rgb != target {
				t.Errorf(AssertTemplate, ce.name, cs, target, rgb)
			}
		}
	}
}

func TestToCmyk(t *testing.T) {
	for _, ce := range getEquivalents() {
		target := ce.cmyk
		for _, cs := range []ColorSpace{ce.pc, ce.rgb, ce.hsl, ce.oklch} {
			cmyk := CMYK{}.FromPrecise(cs.ToPrecise()).(CMYK)
			if cmyk != target {
				t.Errorf(AssertTemplate, ce.name, cs, target, cmyk)
			}
		}
	}
}

func TestToHsl(t *testing.T) {
	for _, ce := range getEquivalents() {
		target := ce.hsl
		for _, cs := range []ColorSpace{ce.pc, ce.rgb, ce.cmyk, ce.oklch} {
			hsl := HSL{}.FromPrecise(cs.ToPrecise()).(HSL)
			if hsl != target {
				t.Errorf(AssertTemplate, ce.name, cs, target, hsl)
			}
		}
	}
}

func TestToOKLCH(t *testing.T) {
	for _, ce := range getEquivalents() {
		target := ce.oklch
		for _, cs := range []ColorSpace{ce.pc, ce.rgb, ce.cmyk, ce.hsl} {
			oklch := OKLCH{}.FromPrecise(cs.ToPrecise()).(OKLCH)
			delta := 1e-2
			if math.Abs(oklch.L-target.L) > delta || math.Abs(oklch.C-target.C) > delta || math.Abs(oklch.H-target.H) > delta {
				t.Errorf(AssertTemplate, ce.name, cs, target, oklch)
			}
		}
	}
}
