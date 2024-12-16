package switcher

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
)

const (
	cpHex   = "x"
	cpRGB   = "r"
	cpHSL   = "s"
	cpCMYK  = "c"
	cpEscFG = "f"
	cpEscBG = "b"
)

type keybinds struct {
	next, prev, copy, help, insert, esc, confirm, quit key.Binding
}

func newKeybinds() keybinds {
	cpKeys := []string{cpHex, cpRGB, cpHSL, cpCMYK, cpEscBG, cpEscFG}
	return keybinds{
		next: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next picker"),
		),
		prev: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev picker"),
		),
		copy: key.NewBinding(
			key.WithKeys(cpKeys...),
			key.WithHelp(
				strings.Join(cpKeys, "/"),
				"copy color",
			),
		),
		help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		insert: key.NewBinding(
			key.WithKeys("i", ":"),
			key.WithHelp("i", "manual input"),
		),
		esc: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "exit manual input"),
			key.WithDisabled(),
		),
		confirm: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm manual input"),
			key.WithDisabled(),
		),
		quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

func Keys() []key.Binding {
	k := newKeybinds()
	return []key.Binding{k.next, k.prev, k.copy, k.insert, k.esc, k.confirm, k.help, k.quit}
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
