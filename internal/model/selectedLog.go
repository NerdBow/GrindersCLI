package model

import (
	"fmt"
	"strings"

	"github.com/NerdBow/GrindersTUI/internal/keymap"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	BackField = iota
	EditField
	DeleteField
	TextField
)

type SelectedLogModel struct {
	log           Log
	previousModel int
	choices       []string
	focusIndex    int
	textField     textinput.Model
	token         string
}

func SelectedLogModelInit(log Log, previousModel int, token string) *SelectedLogModel {
	return &SelectedLogModel{log, previousModel, []string{"Back", "Edit", "Delete"}, 0, textinput.New(), token}
}

func (m *SelectedLogModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *SelectedLogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymap.VimBinding.Exit):
			return m, func() tea.Msg { return ModelMsg{SelectedLog, m.previousModel, nil} }
		case key.Matches(msg, keymap.VimBinding.ChangeFocus):
			m.focusIndex = (m.focusIndex + 1) % len(m.choices)
		case key.Matches(msg, keymap.VimBinding.Select):
			switch m.focusIndex {
			case BackField:
				m.textField.Blur()
				return m, func() tea.Msg { return ModelMsg{SelectedLog, m.previousModel, nil} }
			case EditField:
				m.textField.Blur()
			case DeleteField:
				return m, m.textField.Focus()
			case TextField:
				return m, nil
			}

		}
	}
	var cmd tea.Cmd
	m.textField, cmd = m.textField.Update(msg)
	return m, cmd
}

func (m *SelectedLogModel) View() string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("Id: %d", m.log.Id))
	b.WriteRune('\n')
	b.WriteString(fmt.Sprintf("Date: %s", m.log.DateString()))
	b.WriteRune('\n')
	b.WriteString(fmt.Sprintf("Name: %s", m.log.Name))
	b.WriteRune('\n')
	b.WriteString(fmt.Sprintf("Category: %s", m.log.Category))
	b.WriteRune('\n')
	b.WriteString(fmt.Sprintf("Goal: %s", m.log.Goal))
	b.WriteRune('\n')
	b.WriteString(fmt.Sprintf("Duration: %s", m.log.DurationString()))
	b.WriteRune('\n')
	return b.String()
}
