package quit

import (
	"github.com/ChausseBenjamin/termpicker/internal/ui"
	tea "github.com/charmbracelet/bubbletea/v2"
)

const byeMsg = "Goodbye!\n"

type Model struct {
	failed  bool
	exitMsg string
}

func Ok(msg string) Model   { return Model{exitMsg: msg} }
func Fail(msg string) Model { return Model{failed: true, exitMsg: msg} }

// Cancel is only relevant for oneshot where quiting without returning a
// color is non-ok behaviour which the program/script calling termpicker
// should be able to interpret
func Cancel() Model { return Model{failed: true} }

func (m Model) Failed() bool    { return m.failed }
func (m Model) ExitMsg() string { return m.exitMsg }

func (m Model) Init() tea.Cmd                       { return nil }
func (m Model) Update(tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m Model) View() string {
	if m.failed || m.exitMsg != "" {
		return "" // oneshot: payload is printed after exiting alt-screen
	}
	return ui.Style().Quit.Render(byeMsg)
}
