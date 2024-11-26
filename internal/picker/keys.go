package picker

import (
	"github.com/charmbracelet/bubbles/key"
)

type keybinds struct {
	next, prev key.Binding
}

func newKeybinds() keybinds {
	return keybinds{
		next: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("j", "prev slider"),
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

func (m Model) AllKeys() [][]key.Binding {
	keys := make([][]key.Binding, len(m.sliders[m.active].AllKeys())+1)
	keys[0] = Keys()
	copy(keys[1:], m.sliders[m.active].AllKeys())
	return keys
}

func (m Model) textInputKeys() []key.Binding {
	return []key.Binding{
		m.input.Keymap.CharacterForward,
		m.input.Keymap.CharacterBackward,
		m.input.Keymap.WordForward,
		m.input.Keymap.WordBackward,
		m.input.Keymap.DeleteWordBackward,
		m.input.Keymap.DeleteWordForward,
		m.input.Keymap.DeleteAfterCursor,
		m.input.Keymap.DeleteBeforeCursor,
		m.input.Keymap.DeleteCharacterBackward,
		m.input.Keymap.DeleteCharacterForward,
		m.input.Keymap.LineStart,
		m.input.Keymap.LineEnd,
		m.input.Keymap.Paste,
		m.input.Keymap.AcceptSuggestion,
		m.input.Keymap.NextSuggestion,
		m.input.Keymap.PrevSuggestion,
		m.input.Keymap.textinput.DefaultKeyMap,
	}
}
