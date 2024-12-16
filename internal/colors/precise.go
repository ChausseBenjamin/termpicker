package colors

import (
	"fmt"
	"math"
	"strings"
)

const esc = "\\X1B"

type ColorSpace interface {
	ToPrecise() PreciseColor
	FromPrecise(PreciseColor) ColorSpace
}

func (c PreciseColor) String() string {
	return fmt.Sprintf("PC(%.4f, %.4f, %.4f)", c.R, c.G, c.B)
}

// PreciseColor is a color with floating point values for red, green, and blue.
// The extra precision minimizes rounding errors when converting between different
// color spaces. It is used as an intermediate representation when converting between
// different color spaces.
type PreciseColor struct {
	R, G, B float64
}

func (c PreciseColor) ToPrecise() PreciseColor {
	return c
}

func (c PreciseColor) FromPrecise(p PreciseColor) ColorSpace {
	return p
}

func Hex(cs ColorSpace) string {
	p := cs.ToPrecise()

	return strings.ToUpper(fmt.Sprintf("#%02x%02x%02x",
		int(math.Round(p.R*255)),
		int(math.Round(p.G*255)),
		int(math.Round(p.B*255)),
	))
}

func EscapedSeq(cs ColorSpace, fg bool) string {
	p := cs.ToPrecise()
	mod := 38 // fg by default
	if !fg {
		mod += 10
	}

	r := int(math.Round(p.R * 255))
	g := int(math.Round(p.G * 255))
	b := int(math.Round(p.B * 255))

	return fmt.Sprintf("%s[%d;2;%d;%d;%dm",
		esc, mod, r, g, b,
	)
}
