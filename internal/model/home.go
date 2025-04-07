package model

import (
	"strings"

	"github.com/NerdBow/GrindersTUI/internal/keymap"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type SignOutMsg struct{}

type HomeModelSwitch uint8

const (
	CreateLogField uint8 = iota
	ViewLogsField
	EditLogField
	DeleteLogField
	SignOutField

	CreateLogSwitch HomeModelSwitch = iota
	ViewLogSwitch
	EditLogSwitch
	DeleteLogSwitch
)

type HomeModel struct {
	choices    []string
	focusIndex int
}

func HomeModelInit() *HomeModel {
	return &HomeModel{
		choices:    []string{"Create Log", "View Log(s)", "Edit Log", "Delete Log", "Sign Out"},
		focusIndex: 0,
	}
}

func (m *HomeModel) Init() tea.Cmd {
	return nil
}

func (m *HomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymap.VimBinding.Up):
			m.focusIndex = ((m.focusIndex-1)%len(m.choices) + len(m.choices)) % len(m.choices)
			return m, nil
		case key.Matches(msg, keymap.VimBinding.Down):
			m.focusIndex = (m.focusIndex + 1) % len(m.choices)
			return m, nil
		case key.Matches(msg, keymap.VimBinding.Exit):
			return m, tea.Quit
		case key.Matches(msg, keymap.VimBinding.Select):
			switch uint8(m.focusIndex) {
			case CreateLogField:
				return m, func() tea.Msg { return CreateLogSwitch }
			case ViewLogsField:
				return m, func() tea.Msg { return ViewLogSwitch }
			case EditLogField:
				return m, func() tea.Msg { return EditLogSwitch }
			case DeleteLogField:
				return m, func() tea.Msg { return DeleteLogSwitch }
			case SignOutField:
				return m, func() tea.Msg { return SignOutMsg{} }
			}
			return m, nil
		}

	}
	return m, nil
}

func (m *HomeModel) View() string {
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
