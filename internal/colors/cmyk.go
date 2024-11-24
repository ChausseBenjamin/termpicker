package colors

import (
	"fmt"
	"math"
)

type CMYK struct {
	C int // 0-100
	M int // 0-100
	Y int // 0-100
	K int // 0-100
}

func (c CMYK) String() string {
	return fmt.Sprintf("cmyk(%d%%, %d%%, %d%%, %d%%)", c.C, c.M, c.Y, c.K)
}

func (c CMYK) ToPrecise() PreciseColor {
	return PreciseColor{
		R: (1 - float64(c.C)/100) * (1 - float64(c.K)/100),
		G: (1 - float64(c.M)/100) * (1 - float64(c.K)/100),
		B: (1 - float64(c.Y)/100) * (1 - float64(c.K)/100),
	}
}

func (c CMYK) FromPrecise(p PreciseColor) ColorSpace {
	// Extract RGB components from the PreciseColor
	r := p.R
	g := p.G
	b := p.B

	// Calculate the K (key/black) component
	k := 1 - math.Max(math.Max(r, g), b)

	// Avoid division by zero when K is 1 (pure black)
	if k == 1 {
		return CMYK{C: 0, M: 0, Y: 0, K: 100}
	}

	// Calculate the CMY components based on the remaining color values
	cyan := (1 - r - k) / (1 - k)
	magenta := (1 - g - k) / (1 - k)
	yellow := (1 - b - k) / (1 - k)

	// Scale to 0-100 and return
	return CMYK{
		C: int(math.Round(cyan * 100)),
		M: int(math.Round(magenta * 100)),
		Y: int(math.Round(yellow * 100)),
		K: int(math.Round(k * 100)),
	}
}
