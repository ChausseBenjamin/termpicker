package switcher

import "github.com/charmbracelet/bubbles/key"

type keybinds struct {
	next, prev, cpHex, cpRgb, cpHsl, cpCmyk, help, quit key.Binding
}

func newKeybinds() keybinds {
	return keybinds{
		next: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next picker"),
		),
		prev: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev. picker"),
		),
		cpHex: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("x", "copy hex"),
		),
		cpRgb: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "copy rgb"),
		),
		cpHsl: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "copy hsl"),
		),
		cpCmyk: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "copy cmyk"),
		),
		help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

func Keys() []key.Binding {
	k := newKeybinds()
	return []key.Binding{k.next, k.prev, k.cpHex, k.cpRgb, k.cpHsl, k.cpCmyk, k.help, k.quit}
}

func shortKeys() [][]key.Binding {
	keys := make([][]key.Binding, 2)
	rows := 2
	cRow := 0
	for i := 0; i < len(Keys()); i++ {
		keys[cRow] = append(keys[cRow], Keys()[i])
		cRow++
		if cRow == rows {
			cRow = 0
		}
	}
	return keys
}

func (m Model) AllKeys() [][]key.Binding {
	keys := make([][]key.Binding, len(m.pickers[m.active].AllKeys())+1)
	keys[0] = Keys()
	copy(keys[1:], m.pickers[m.active].AllKeys())
	return keys
	// return append(m.pickers[m.active].AllKeys(), Keys())
}
