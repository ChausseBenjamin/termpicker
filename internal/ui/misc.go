package ui

import (
	lg "github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func ColorProfile() termenv.Profile {
	return termenv.TrueColor
}

func init() {
	r := lg.DefaultRenderer()
	r.SetColorProfile(ColorProfile())
	lg.SetDefaultRenderer(r)
}
