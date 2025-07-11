package notices

import (
	"log/slog"
	"strings"
	"time"

	"github.com/ChausseBenjamin/termpicker/internal/util"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/hashicorp/go-uuid"
)

const (
	expiryDelay = 3 // seconds
)

type NoticeExpiryMsg string

type Model struct {
	// Notices is a map of UUIDs pointing to a messages
	Notices map[string]string
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) View() string {
	noticeStr := ""
	for _, v := range m.Notices {
		noticeStr += v + "\n"
	}
	return strings.TrimRight(noticeStr, "\n")
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case NoticeExpiryMsg:
		delete(m.Notices, string(msg))
		return m, nil
	}
	return m, nil
}

func New() Model {
	return Model{
		Notices: make(map[string]string),
	}
}

func (m Model) New(msg string) tea.Cmd {
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		slog.Error("Failed to generate UUID", util.ErrKey, err)
	}
	m.Notices[uuid] = msg

	return func() tea.Msg {
		time.Sleep(expiryDelay * time.Second)
		return NoticeExpiryMsg(uuid)
	}
}

func (m Model) Reset(uuid string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(expiryDelay * time.Second)
		return NoticeExpiryMsg(uuid)
	}
}
