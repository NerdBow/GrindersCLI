package model

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type SignInModel struct {
	inputs      []textinput.Model
	cursorIndex uint8
	cursor      cursor.Mode
}

func SignInModelInit() SignInModel {
	m := SignInModel{
		inputs: make([]textinput.Model, 2),
	}
	usernameTextInput := textinput.New()
	usernameTextInput.CharLimit = 64
	usernameTextInput.Placeholder = "Username"
	usernameTextInput.Focus()

	m.inputs[0] = usernameTextInput

	passwordTextInput := textinput.New()
	passwordTextInput.CharLimit = 64
	passwordTextInput.Placeholder = "Password"
	passwordTextInput.EchoMode = textinput.EchoPassword
	passwordTextInput.EchoCharacter = '•'

	m.inputs[1] = passwordTextInput

	return m
}

func (m SignInModel) Init() tea.Cmd {
	return nil
}

func (m SignInModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m SignInModel) View() string {
	var b strings.Builder
	return ""
}
