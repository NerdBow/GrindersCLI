package model

import (
	"github.com/NerdBow/GrindersTUI/internal/keymap"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const ()

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type RecentLogsModel struct {
	token      string
	page       int
	logTable   table.Model
	focusIndex int
	status     string
	logs       []Log
}

func RecentLogsModelInit(token string) *RecentLogsModel {
	columns := []table.Column{
		{Title: "Id", Width: 10},
		{Title: "Date", Width: 10},
		{Title: "Duration", Width: 8},
		{Title: "Name", Width: 25},
		{Title: "Category", Width: 25},
		{Title: "Goal", Width: 75},
	}
	return &RecentLogsModel{
		token:      token,
		page:       1,
		logTable:   table.New(table.WithColumns(columns)),
		focusIndex: 0,
	}
}

func (m *RecentLogsModel) Init() tea.Cmd {
	return nil
}

func (m *RecentLogsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymap.VimBinding.Exit):
			return m, func() tea.Msg { return ModelMsg{RecentLogs, ViewLog, nil} }
		}
	}
	return m, nil
}

func (m *RecentLogsModel) View() string {
	return baseStyle.Render(m.logTable.View()) + "\n"
}

}
