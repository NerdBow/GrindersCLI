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
	logTable   table.Model
	focusIndex int
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
