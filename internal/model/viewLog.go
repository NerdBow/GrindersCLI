package model

import (
	"github.com/NerdBow/GrindersTUI/internal/keymap"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

const (
	RecentLogsField int = iota
	IdSearchField
	CustomSearchField
)

type ViewLogModel struct {
	focusIndex int
	choices    []string
}

func ViewLogModelInit() *ViewLogModel {
	return &ViewLogModel{
		choices:    []string{"Recent Logs", "Id Search", "Custom Search"},
		focusIndex: 0,
	}
}

func (m *ViewLogModel) Init() tea.Cmd {
	return nil
}

func (m *ViewLogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *ViewLogModel) View() string {
	b := strings.Builder{}
	return b.String()
}
