package ui

import (
	"github.com/ChausseBenjamin/termpicker/internal/progress"
	lg "github.com/charmbracelet/lipgloss"
)

const (
	// Colors:
	textSel   = "#F2F1F0"
	textNorm  = "#A7AFB1"
	textFaint = "#6F797B"
	geomFg    = "#ACB3B5"

	TabSepLeft  = "["
	TabSepMid   = " | "
	TabSepRight = "]"

	PickerSelRune = ">"

	PromptPrefix      = "> "
	PromptPlaceholder = "Enter a color (ex: #b7416e)"

	SliderMinWidth = 22 // 1 ASCII change every 2.05 deg. avg
	SliderMaxWidth = 90 // 2 ASCII change per deg.

)

type sliderOpts struct {
	R, G, B    []progress.Option
	C, M, Y, K []progress.Option
	H, S, L    []progress.Option
}

type StyleSheet struct {
	TabSel       lg.Style
	TabNorm      lg.Style
	TabGeom      lg.Style
	SliderVal    lg.Style
	SliderLabel  lg.Style
	PickerCursor lg.Style
	Preview      lg.Style
	InputPrompt  lg.Style
	InputText    lg.Style
	Notice       lg.Style
	Quit         lg.Style
	Boxed        lg.Style
	Sliders      sliderOpts
}

var style StyleSheet

func Style() StyleSheet {
	return style
}

func init() {
	baseStyle := lg.NewStyle().
		Foreground(lg.Color(textNorm))

	baseSliderOpts := []progress.Option{
		progress.WithColorProfile(ColorProfile()),
		progress.WithoutPercentage(),
		// progress.WithBinaryFill(), // uncomment for legacy look
	}

	style = StyleSheet{
		TabSel: baseStyle.Inherit(lg.NewStyle().
			Foreground(lg.Color(textSel)).
			Underline(true).
			Bold(true)),

		TabNorm: baseStyle.Inherit(lg.NewStyle().
			Foreground(lg.Color(textFaint)).
			Underline(false).
			Bold(false)),

		TabGeom: baseStyle.Inherit(lg.NewStyle().
			Foreground(lg.Color(geomFg))),

		SliderVal: baseStyle,

		SliderLabel: baseStyle,

		PickerCursor: baseStyle.Inherit(lg.NewStyle().
			Bold(true)),

		Preview: baseStyle,

		InputPrompt: baseStyle.Inherit(lg.NewStyle().
			Bold(true)),

		InputText: baseStyle,

		Notice: baseStyle.Inherit(lg.NewStyle().
			Bold(true)),

		Quit: baseStyle.Inherit(lg.NewStyle().
			Foreground(lg.Color(textSel)).
			Bold(true)),

		Boxed: baseStyle.Inherit(lg.NewStyle().
			Border(lg.RoundedBorder())),

		Sliders: sliderOpts{
			// RGB
			R: append(baseSliderOpts, progress.WithGradient("#660000", "#ff0000")),
			G: append(baseSliderOpts, progress.WithGradient("#006600", "#00ff00")),
			B: append(baseSliderOpts, progress.WithGradient("#000066", "#0000ff")),

			// CMYK
			C: append(baseSliderOpts, progress.WithGradient("#006666", "#00ffff")),
			M: append(baseSliderOpts, progress.WithGradient("#660066", "#ff00ff")),
			Y: append(baseSliderOpts, progress.WithGradient("#666600", "#ffff00")),
			K: append(baseSliderOpts, progress.WithSolidFill("#000000")),

			// HSL
			H: append(baseSliderOpts, progress.WithDefaultGradient()),
			S: append(baseSliderOpts, progress.WithGradient("#a68e59", "#ffae00")),
			L: append(baseSliderOpts, progress.WithGradient("#222222", "#ffffff")),
		},
	}
}
