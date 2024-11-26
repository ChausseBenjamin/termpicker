package parse

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ChausseBenjamin/termpicker/internal/colors"
)

var (
	errUnknownColorFormat = errors.New("Unrecognized color format")
	errHexParsing         = errors.New("Failed to parse hex color")
	errRGBParsing         = errors.New("Failed to parse RGB color")
	errHSLParsing         = errors.New("Failed to parse HSL color")
	errCMYKParsing        = errors.New("Failed to parse CMYK color")
)

func sanitize(s string) string {
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.ReplaceAll(s, "%", "")
	s = strings.ReplaceAll(s, "°", "")
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	return s
}

func Color(s string) (colors.ColorSpace, error) {
	s = sanitize(s)
	switch {
	case strings.Contains(s, "#"):
		return hex(s)
	case strings.Contains(s, "rgb"):
		return rgb(s)
	case strings.Contains(s, "hsl"):
		return hsl(s)
	case strings.Contains(s, "cmyk"):
		return cmyk(s)
	default:
		return nil, errUnknownColorFormat
	}
}

func rgb(s string) (colors.ColorSpace, error) {
	var r, g, b int
	_, err := fmt.Sscanf(s, "rgb(%d,%d,%d)", &r, &g, &b)
	if err != nil {
		return nil, errors.Join(errRGBParsing, err)
	}
	return colors.RGB{R: r, G: g, B: b}, nil
}

func hex(s string) (colors.ColorSpace, error) {
	var r, g, b int
	_, err := fmt.Sscanf(s, "#%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return nil, errors.Join(errHexParsing, err)
	}
	return colors.RGB{R: r, G: g, B: b}, nil
}

func cmyk(s string) (colors.ColorSpace, error) {
	var c, m, y, k int
	_, err := fmt.Sscanf(s, "cmyk(%d,%d,%d,%d)", &c, &m, &y, &k)
	if err != nil {
		return nil, errors.Join(errCMYKParsing, err)
	}
	return colors.CMYK{C: c, M: m, Y: y, K: k}, nil
}

func hsl(str string) (colors.ColorSpace, error) {
	var h, s, l int
	_, err := fmt.Sscanf(str, "hsl(%d,%d,%d)", &h, &s, &l)
	if err != nil {
		return nil, errors.Join(errHSLParsing, err)
	}
	return colors.HSL{H: h, S: s, L: l}, nil
}
