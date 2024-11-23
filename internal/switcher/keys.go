package switcher

import "github.com/charmbracelet/bubbles/key"

type keybinds struct {
	next, prev, quit key.Binding
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
		quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

func Keys() []key.Binding {
	k := newKeybinds()
	return []key.Binding{k.next, k.prev}
}

func (m Model) AllKeys() []key.Binding {
	return append(Keys(), m.pickers[m.active].AllKeys()...)
}
