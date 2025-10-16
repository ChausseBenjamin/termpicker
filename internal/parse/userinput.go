package parse

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/ChausseBenjamin/termpicker/internal/colors"
)

var (
	errUnknownColorFormat = errors.New("unrecognized color format")
	errHexParsing         = errors.New("failed to parse hex color")
	errRGBParsing         = errors.New("failed to parse RGB color")
	errHSLParsing         = errors.New("failed to parse HSL color")
	errCMYKParsing        = errors.New("failed to parse CMYK color")
	errOKLCHParsing       = errors.New("failed to parse OKLCH color")
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
	if strings.Contains(s, "oklch") {
		return oklch(s)
	}
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

func oklch(s string) (colors.ColorSpace, error) {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "oklch(from ") {
		return oklchRelative(s)
	}
	// Absolute
	s = strings.TrimPrefix(s, "oklch(")
	s = strings.TrimSuffix(s, ")")
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == ' ' || r == '/'
	})
	if len(parts) < 3 {
		return nil, errors.Join(errOKLCHParsing, errors.New("not enough components"))
	}
	L, err := parseValue(parts[0])
	if err != nil {
		return nil, errors.Join(errOKLCHParsing, err)
	}
	C, err := parseValue(parts[1])
	if err != nil {
		return nil, errors.Join(errOKLCHParsing, err)
	}
	H, err := parseValue(parts[2])
	if err != nil {
		return nil, errors.Join(errOKLCHParsing, err)
	}
	return colors.OKLCH{L: L, C: C, H: H}, nil
}

func oklchRelative(s string) (colors.ColorSpace, error) {
	s = strings.TrimPrefix(s, "oklch(from ")
	s = strings.TrimSuffix(s, ")")
	parts := strings.SplitN(s, " ", 2)
	if len(parts) != 2 {
		return nil, errors.Join(errOKLCHParsing, errors.New("invalid relative format"))
	}
	originStr := parts[0]
	valuesStr := parts[1]
	origin, err := Color(originStr)
	if err != nil {
		return nil, errors.Join(errOKLCHParsing, err)
	}
	originPrecise := origin.ToPrecise()
	originOk := colors.OKLCH{}.FromPrecise(originPrecise).(colors.OKLCH)
	valueParts := strings.Fields(valuesStr)
	if len(valueParts) != 3 {
		return nil, errors.Join(errOKLCHParsing, errors.New("not enough values"))
	}
	L, err := parseComponent(valueParts[0], originOk)
	if err != nil {
		return nil, errors.Join(errOKLCHParsing, err)
	}
	C, err := parseComponent(valueParts[1], originOk)
	if err != nil {
		return nil, errors.Join(errOKLCHParsing, err)
	}
	H, err := parseComponent(valueParts[2], originOk)
	if err != nil {
		return nil, errors.Join(errOKLCHParsing, err)
	}
	return colors.OKLCH{L: L, C: C, H: H}, nil
}

func parseComponent(s string, origin colors.OKLCH) (float64, error) {
	s = strings.TrimSpace(s)
	switch s {
	case "l":
		return origin.L, nil
	case "c":
		return origin.C, nil
	case "h":
		return origin.H, nil
	default:
		return parseValue(s)
	}
}

func parseValue(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, errors.New("empty value")
	}
	if strings.HasSuffix(s, "%") {
		s = strings.TrimSuffix(s, "%")
		v, err := strconv.ParseFloat(s, 64)
		return v / 100, err
	}
	if strings.HasSuffix(s, "deg") {
		s = strings.TrimSuffix(s, "deg")
		return strconv.ParseFloat(s, 64)
	}
	if strings.HasSuffix(s, "rad") {
		s = strings.TrimSuffix(s, "rad")
		v, err := strconv.ParseFloat(s, 64)
		return v * 180 / math.Pi, err
	}
	if strings.HasSuffix(s, "turn") {
		s = strings.TrimSuffix(s, "turn")
		v, err := strconv.ParseFloat(s, 64)
		return v * 360, err
	}
	// Plain number
	return strconv.ParseFloat(s, 64)
}
