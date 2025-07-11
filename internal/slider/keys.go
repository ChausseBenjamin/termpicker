package slider

import "github.com/charmbracelet/bubbles/v2/key"

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
			key.WithHelp("l", "+5% slider"),
		),
		decRegular: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("h", "-5% slider"),
		),
		incPrecise: key.NewBinding(
			key.WithKeys("shift+right", "L"),
			key.WithHelp("L", "+1 slider"),
		),
		decPrecise: key.NewBinding(
			key.WithKeys("shift+left", "H"),
			key.WithHelp("H", "-1 slider"),
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
	return [][]key.Binding{Keys()}
}
