package model

import (
	"github.com/NerdBow/GrindersTUI/internal/keymap"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type CreateLogModelSwitch uint8

type EmptyFieldErrorMsg uint8

const (
	NameField     uint8 = 0
	CategoryField       = 1
	GoalField           = 2
	ConfirmButton       = iota

	HomeSwitch CreateLogModelSwitch = iota
	TimerSwitch

	UsernameFieldEmpty EmptyFieldErrorMsg = iota
	CategoryFieldEmpty
	GoalFieldEmpty
)

type CreateLogModel struct {
	focusIndex int
	inputs     []textinput.Model
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
	return m, nil
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

	return b.String()
}
