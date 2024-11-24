package switcher

import "github.com/charmbracelet/bubbles/key"

type keybinds struct {
	next, prev, cpHex, cpRgb, cpHsl, cpCmyk, quit key.Binding
}

func newKeybinds() keybinds {
	return keybinds{
		next: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next picker"),
		),
		prev: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "previous picker"),
		),
		cpHex: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("x", "yank/copy hex value"),
		),
		cpRgb: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "yank/copy RGB value"),
		),
		cpHsl: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "yank/copy HSL value"),
		),
		cpCmyk: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "yank/copy CMYK value"),
		),
		quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

func Keys() []key.Binding {
	k := newKeybinds()
	return []key.Binding{k.next, k.prev, k.cpHex, k.cpRgb, k.cpHsl, k.cpCmyk, k.quit}
}

func (m Model) AllKeys() []key.Binding {
	return append(Keys(), m.pickers[m.active].AllKeys()...)
}
