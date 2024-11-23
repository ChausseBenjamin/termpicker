package picker

import "github.com/charmbracelet/bubbles/key"

type keybinds struct {
	next, prev key.Binding
}

func newKeybinds() keybinds {
	return keybinds{
		next: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("j", "previous slider"),
		),
		prev: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("k", "next slider"),
		),
	}
}

func Keys() []key.Binding {
	k := newKeybinds()
	return []key.Binding{k.next, k.prev}
}

func (m Model) AllKeys() []key.Binding {
	return append(Keys(), m.sliders[m.active].AllKeys()...)
}
