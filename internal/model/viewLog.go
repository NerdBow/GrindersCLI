package model

import (
	"strings"

	"github.com/NerdBow/GrindersTUI/internal/keymap"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymap.VimBinding.Select):
			switch m.focusIndex {
			case RecentLogsField:
				return m, func() tea.Msg { return ModelMsg{ViewLog, LogTable, nil} }
			case IdSearchField:
				return m, func() tea.Msg { return ModelMsg{ViewLog, IdLogSearch, nil} }
			case CustomSearchField:
				return m, func() tea.Msg { return ModelMsg{ViewLog, CustomLogSearch, nil} }
			}
		case key.Matches(msg, keymap.VimBinding.Up):
			m.focusIndex = ((m.focusIndex-1)%len(m.choices) + len(m.choices)) % len(m.choices)
			return m, nil
		case key.Matches(msg, keymap.VimBinding.Down):
			m.focusIndex = (m.focusIndex + 1) % len(m.choices)
			return m, nil
		case key.Matches(msg, keymap.VimBinding.Exit):
			return m, func() tea.Msg { return ModelMsg{ViewLog, Home, nil} }
		}
	}
	return m, nil
}

func (m *ViewLogModel) View() string {
	var b strings.Builder

	for i := range m.choices {
		if i == m.focusIndex {
			b.WriteString("> ")
		}

		b.WriteString(m.choices[i])

		if i < len(m.choices)-1 {
			b.WriteRune('\n')
		}
	}
	return b.String()
}
