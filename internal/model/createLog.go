package model

import (
	"github.com/NerdBow/GrindersTUI/internal/keymap"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type EmptyFieldErrorMsg uint8

const (
	NameField     uint8 = 0
	CategoryField       = 1
	GoalField           = 2
	ConfirmButton       = iota

	UsernameFieldEmpty EmptyFieldErrorMsg = iota
	CategoryFieldEmpty
	GoalFieldEmpty
)

type CreateLogModel struct {
	focusIndex   int
	inputs       []textinput.Model
	errorMessage string
}

func CreateLogModelInit() *CreateLogModel {
	m := CreateLogModel{
		inputs: make([]textinput.Model, 3),
	}
	for i := range m.inputs {
		t := textinput.New()
		t.CharLimit = 100
		t.PromptStyle = textInputUnfocusedStyle
		t.TextStyle = textInputUnfocusedStyle
		switch uint8(i) {
		case NameField:
			t.Placeholder = "Name"
			t.PromptStyle = textInputFocusedStyle
			t.TextStyle = textInputFocusedStyle
			t.Focus()
		case CategoryField:
			t.Placeholder = "Category"
		case GoalField:
			t.Placeholder = "Goal"
		}
		m.inputs[i] = t
	}
	return &m
}

func (m *CreateLogModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *CreateLogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymap.VimBinding.ChangeFocus):
			m.focusIndex = (m.focusIndex + 1) % 4

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = textInputFocusedStyle
					m.inputs[i].TextStyle = textInputFocusedStyle
					continue
				}
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = textInputUnfocusedStyle
				m.inputs[i].TextStyle = textInputUnfocusedStyle
			}
			return m, tea.Batch(cmds...)
		case key.Matches(msg, keymap.VimBinding.Exit):
			return m, func() tea.Msg { return ModelMsg{CreateLog, Home, nil} }
		case key.Matches(msg, keymap.VimBinding.Select):
			switch uint8(m.focusIndex) {
			case ConfirmButton:
				if m.IsInputsEmpty() {
					m.errorMessage = "There can be no empty fields!"
					return m, nil
				}
				return m, func() tea.Msg { return ModelMsg{CreateLog, Stopwatch, nil} }
			}
		}
	}
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m *CreateLogModel) View() string {
	var b strings.Builder
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		b.WriteRune('\n')
	}
	confirmChoice := textInputUnfocusedStyle.Render("Confirm")
	if m.focusIndex == ConfirmButton {
		confirmChoice = textInputFocusedStyle.Render("Confirm")
	}

	b.WriteString(confirmChoice)
	b.WriteRune('\n')
	b.WriteRune('\n')

	b.WriteString(textInputFocusedStyle.Render(m.errorMessage))

	return b.String()
}

func (m *CreateLogModel) IsInputsEmpty() bool {
	for _, t := range m.inputs {
		if t.Value() == "" {
			return true
		}
	}
	return false
}

func (m *CreateLogModel) GetLogInfo() (string, string, string) {
	return m.inputs[0].Value(), m.inputs[1].Value(), m.inputs[2].Value()
}
