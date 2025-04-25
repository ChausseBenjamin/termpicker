package slider

import "github.com/charmbracelet/bubbles/key"

type keybinds struct {
	incRegular key.Binding
	decRegular key.Binding
	incPrecise key.Binding
	decPrecise key.Binding
}

func newKeybinds() keybinds {
	return keybinds{
		incRegular: key.NewBinding(
			key.WithKeys("right", "l"),
		),
		decRegular: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("h/l", "Dec./Inc. 5%"),
		),
		incPrecise: key.NewBinding(
			key.WithKeys("shift+right", "L"),
		),
		decPrecise: key.NewBinding(
			key.WithKeys("shift+left", "H"),
			key.WithHelp("H/L", "Dec./Inc. (fine)"),
		),
	}
}

// Join all keybindings into a single slice
// a parent can use to know what Keys
// it's children have.
func Keys() []key.Binding {
	k := newKeybinds()
	return []key.Binding{
		k.incRegular,
		k.decRegular,
		k.incPrecise,
		k.decPrecise,
	}
}

// AllKeys returns key.Bindings for the Model
// and all of its active children. The parent
// can use this to generate help text.
func (m Model) AllKeys() [][]key.Binding {
	k := newKeybinds()
	return [][]key.Binding{{k.decRegular, k.decPrecise}}
}
