package parse

import (
	"math"
	"testing"

	"github.com/ChausseBenjamin/termpicker/internal/colors"
)

func TestColorParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected colors.ColorSpace
		hasError bool
	}{
		// Hex formats
		{"hex basic", "#ff0000", colors.RGB{R: 255, G: 0, B: 0}, false},
		{"hex lowercase", "#00ff00", colors.RGB{R: 0, G: 255, B: 0}, false},
		{"hex uppercase", "#0000FF", colors.RGB{R: 0, G: 0, B: 255}, false},
		{"hex mixed case", "#AbCdEf", colors.RGB{R: 171, G: 205, B: 239}, false},

		// RGB formats
		{"rgb basic", "rgb(255,0,0)", colors.RGB{R: 255, G: 0, B: 0}, false},
		{"rgb green", "rgb(0,255,0)", colors.RGB{R: 0, G: 255, B: 0}, false},
		{"rgb blue", "rgb(0,0,255)", colors.RGB{R: 0, G: 0, B: 255}, false},
		{"rgb black", "rgb(0,0,0)", colors.RGB{R: 0, G: 0, B: 0}, false},
		{"rgb white", "rgb(255,255,255)", colors.RGB{R: 255, G: 255, B: 255}, false},

		// HSL formats
		{"hsl red", "hsl(0,100,50)", colors.HSL{H: 0, S: 100, L: 50}, false},
		{"hsl green", "hsl(120,100,50)", colors.HSL{H: 120, S: 100, L: 50}, false},
		{"hsl blue", "hsl(240,100,50)", colors.HSL{H: 240, S: 100, L: 50}, false},
		{"hsl gray", "hsl(0,0,50)", colors.HSL{H: 0, S: 0, L: 50}, false},
		{"hsl white", "hsl(0,0,100)", colors.HSL{H: 0, S: 0, L: 100}, false},
		{"hsl black", "hsl(0,0,0)", colors.HSL{H: 0, S: 0, L: 0}, false},

		// CMYK formats
		{"cmyk red", "cmyk(0,100,100,0)", colors.CMYK{C: 0, M: 100, Y: 100, K: 0}, false},
		{"cmyk green", "cmyk(100,0,100,0)", colors.CMYK{C: 100, M: 0, Y: 100, K: 0}, false},
		{"cmyk blue", "cmyk(100,100,0,0)", colors.CMYK{C: 100, M: 100, Y: 0, K: 0}, false},
		{"cmyk black", "cmyk(0,0,0,100)", colors.CMYK{C: 0, M: 0, Y: 0, K: 100}, false},
		{"cmyk white", "cmyk(0,0,0,0)", colors.CMYK{C: 0, M: 0, Y: 0, K: 0}, false},

		// OKLCH absolute formats
		{"oklch basic", "oklch(0.5 0.2 120)", colors.OKLCH{L: 0.5, C: 0.2, H: 120}, false},
		{"oklch percent L", "oklch(50% 0.2 120)", colors.OKLCH{L: 0.5, C: 0.2, H: 120}, false},
		{"oklch percent C", "oklch(0.5 20% 120)", colors.OKLCH{L: 0.5, C: 0.2, H: 120}, false},
		{"oklch deg", "oklch(0.5 0.2 120deg)", colors.OKLCH{L: 0.5, C: 0.2, H: 120}, false},
		{"oklch rad", "oklch(0.5 0.2 2rad)", colors.OKLCH{L: 0.5, C: 0.2, H: 114.59155902616465}, false},
		{"oklch turn", "oklch(0.5 0.2 0.5turn)", colors.OKLCH{L: 0.5, C: 0.2, H: 180}, false},
		{"oklch with alpha", "oklch(0.5 0.2 120 / 0.8)", colors.OKLCH{L: 0.5, C: 0.2, H: 120}, false},

		// OKLCH relative formats
		{"oklch relative red", "oklch(from #ff0000 l c h)", colors.OKLCH{L: 0.627987, C: 0.257640, H: 29.227136}, false},
		{"oklch relative modified", "oklch(from #ff0000 0.8 0.4 h)", colors.OKLCH{L: 0.8, C: 0.4, H: 29.227136}, false},

		// Error cases
		{"invalid format", "invalid", nil, true},
		{"empty string", "", nil, true},
		{"malformed hex", "#xyz", nil, true},
		{"malformed rgb", "rgb(abc,def,ghi)", nil, true},
		{"malformed hsl", "hsl(abc,def,ghi)", nil, true},
		{"malformed cmyk", "cmyk(abc,def,ghi,jkl)", nil, true},
		{"malformed oklch", "oklch(abc def ghi)", nil, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Color(test.input)
			if test.hasError {
				if err == nil {
					t.Errorf("Expected error for %s, got none", test.input)
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error for %s: %v", test.input, err)
				return
			}

			// Type-specific assertions
			switch expected := test.expected.(type) {
			case colors.RGB:
				actual, ok := result.(colors.RGB)
				if !ok {
					t.Errorf("Expected RGB for %s, got %T", test.input, result)
					return
				}
				if actual != expected {
					t.Errorf("For %s, expected %v, got %v", test.input, expected, actual)
				}
			case colors.HSL:
				actual, ok := result.(colors.HSL)
				if !ok {
					t.Errorf("Expected HSL for %s, got %T", test.input, result)
					return
				}
				if actual != expected {
					t.Errorf("For %s, expected %v, got %v", test.input, expected, actual)
				}
			case colors.CMYK:
				actual, ok := result.(colors.CMYK)
				if !ok {
					t.Errorf("Expected CMYK for %s, got %T", test.input, result)
					return
				}
				if actual != expected {
					t.Errorf("For %s, expected %v, got %v", test.input, expected, actual)
				}
			case colors.OKLCH:
				actual, ok := result.(colors.OKLCH)
				if !ok {
					t.Errorf("Expected OKLCH for %s, got %T", test.input, result)
					return
				}
				delta := 1e-3
				if math.Abs(actual.L-expected.L) > delta || math.Abs(actual.C-expected.C) > delta || math.Abs(actual.H-expected.H) > delta {
					t.Errorf("For %s, expected L=%.6f C=%.6f H=%.6f, got L=%.6f C=%.6f H=%.6f", test.input, expected.L, expected.C, expected.H, actual.L, actual.C, actual.H)
				}
			default:
				t.Errorf("Unsupported expected type: %T", expected)
			}
		})
	}
}
