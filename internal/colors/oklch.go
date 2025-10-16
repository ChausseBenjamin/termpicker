package colors

import (
	"fmt"
	"math"
)

type OKLCH struct {
	L float64 // Lightness 0-1
	C float64 // Chroma 0-0.5 (unbounded but typically)
	H float64 // Hue 0-360 degrees
}

func (o OKLCH) String() string {
	return fmt.Sprintf("oklch(%.1f%% %.3f %.2f)", o.L*100, o.C, o.H)
}

func (o OKLCH) ToPrecise() PreciseColor {
	// Convert OKLCH to Oklab first
	hRad := o.H * math.Pi / 180.0
	a := o.C * math.Cos(hRad)
	b := o.C * math.Sin(hRad)

	// Convert Oklab to linear RGB using the conversion matrices from Wikipedia
	// M2^-1 matrix (Oklab to l'm's')
	lPrime := 0.9999999984505197*o.L + 0.39633779217376786*a + 0.2158037580607588*b
	mPrime := 1.0000000088817607*o.L - 0.10556134232365635*a - 0.06385417477170591*b
	sPrime := 1.0000000546724108*o.L - 0.08948418209496575*a - 1.2914855378640917*b

	// Apply cube (inverse of cube root)
	l := lPrime * lPrime * lPrime
	m := mPrime * mPrime * mPrime
	s := sPrime * sPrime * sPrime

	// M1^-1 matrix (lms to XYZ)
	x := 1.2268798733741557*l - 0.5578149965554813*m + 0.28139105017721583*s
	y := -0.04057576262431372*l + 1.1122868293970594*m - 0.07171106666151701*s
	z := -0.07637294974672142*l - 0.4214933239627914*m + 1.5869240244272418*s

	// Convert XYZ to linear RGB (sRGB matrix)
	rLinear := 3.2406254773200533*x - 1.5372079722103187*y - 0.4986285986588718*z
	gLinear := -0.9689307147293197*x + 1.8757560608852415*y + 0.041517523842953964*z
	bLinear := 0.055710120445510616*x - 0.2040259135167538*y + 1.0569715142428784*z

	// Apply gamma correction (sRGB curve)
	r := linearToSRGB(rLinear)
	g := linearToSRGB(gLinear)
	bVal := linearToSRGB(bLinear)

	// Clamp to [0,1] range
	r = math.Max(0, math.Min(1, r))
	g = math.Max(0, math.Min(1, g))
	bVal = math.Max(0, math.Min(1, bVal))

	return PreciseColor{R: r, G: g, B: bVal}
}

func (o OKLCH) FromPrecise(p PreciseColor) ColorSpace {
	// Convert RGB to linear RGB
	rLinear := srgbToLinear(p.R)
	gLinear := srgbToLinear(p.G)
	bLinear := srgbToLinear(p.B)

	// Convert linear RGB to XYZ (inverse sRGB matrix)
	x := 0.4124564390896922*rLinear + 0.357576077643909*gLinear + 0.18043748326639894*bLinear
	y := 0.21267285140562253*rLinear + 0.715152155287818*gLinear + 0.07217499330655958*bLinear
	z := 0.019333895582329317*rLinear + 0.11919202588130297*gLinear + 0.9503040785363677*bLinear

	// Convert XYZ to lms using M1 matrix
	l := 0.8189330101*x + 0.3618667424*y - 0.1288597137*z
	m := 0.0329845436*x + 0.9293118715*y + 0.0361456387*z
	s := 0.0482003018*x + 0.2643662691*y + 0.6338517070*z

	// Apply cube root
	lPrime := math.Cbrt(l)
	mPrime := math.Cbrt(m)
	sPrime := math.Cbrt(s)

	// Convert to Oklab using M2 matrix
	lightness := 0.2104542553*lPrime + 0.7936177850*mPrime - 0.0040720468*sPrime
	a := 1.9779984951*lPrime - 2.4285922050*mPrime + 0.4505937099*sPrime
	b := 0.0259040371*lPrime + 0.7827717662*mPrime - 0.8086757660*sPrime

	// Convert Oklab to OKLCH
	chroma := math.Sqrt(a*a + b*b)
	hue := 0.0
	if chroma >= 1e-4 {
		hue = math.Atan2(b, a) * 180.0 / math.Pi
		if hue < 0 {
			hue += 360
		}
	}

	return OKLCH{
		L: lightness,
		C: chroma,
		H: hue,
	}
}

func linearToSRGB(linear float64) float64 {
	if linear <= 0.0031308 {
		return 12.92 * linear
	}
	return 1.055*math.Pow(linear, 1.0/2.4) - 0.055
}

func srgbToLinear(srgb float64) float64 {
	if srgb <= 0.04045 {
		return srgb / 12.92
	}
	return math.Pow((srgb+0.055)/1.055, 2.4)
}
