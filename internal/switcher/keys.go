package switcher

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
)

const (
	cpHex  = "x"
	cpRGB  = "r"
	cpHSL  = "s"
	cpCMYK = "c"
)

type keybinds struct {
	next, prev, copy, help, insert, quit key.Binding
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
		copy: key.NewBinding(
			key.WithKeys(cpHex, cpRGB, cpHSL, cpCMYK),
			key.WithHelp(
				strings.Join([]string{cpHex, cpRGB, cpHSL, cpCMYK}, "/"),
				"copy color",
			),
		),
		help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		insert: key.NewBinding(
			key.WithKeys("i", ":"),
		),
		quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

func Keys() []key.Binding {
	k := newKeybinds()
	return []key.Binding{k.next, k.prev, k.copy, k.help, k.quit}
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

func (m Model) textInputKeys() []key.Binding {
	return []key.Binding{
		m.input.KeyMap.CharacterForward,
		m.input.KeyMap.CharacterBackward,
		m.input.KeyMap.WordForward,
		m.input.KeyMap.WordBackward,
		m.input.KeyMap.DeleteWordBackward,
		m.input.KeyMap.DeleteWordForward,
		m.input.KeyMap.DeleteAfterCursor,
		m.input.KeyMap.DeleteBeforeCursor,
		m.input.KeyMap.DeleteCharacterBackward,
		m.input.KeyMap.DeleteCharacterForward,
		m.input.KeyMap.LineStart,
		m.input.KeyMap.LineEnd,
		m.input.KeyMap.Paste,
		m.input.KeyMap.AcceptSuggestion,
		m.input.KeyMap.NextSuggestion,
		m.input.KeyMap.PrevSuggestion,
	}
}
