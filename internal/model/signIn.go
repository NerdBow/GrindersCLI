package model

import (
	"strings"

	"github.com/NerdBow/GrindersTUI/internal/keymap"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	textInputFocusedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff88aa"))
	textInputUnfocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
)

type SignInModel struct {
	inputs     []textinput.Model
	focusIndex int
}

func SignInModelInit() SignInModel {
	m := SignInModel{
		inputs: make([]textinput.Model, 2),
	}
	usernameTextInput := textinput.New()
	usernameTextInput.CharLimit = 64
	usernameTextInput.Placeholder = "Username"
	usernameTextInput.PromptStyle = textInputFocusedStyle
	usernameTextInput.TextStyle = textInputFocusedStyle
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
	return textinput.Blink
}

func (m SignInModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymap.VimBinding.ChangeFocus):
			m.focusIndex = (m.focusIndex + 1) % 3

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
		case key.Matches(msg, keymap.VimBinding.Select):
			if m.focusIndex == 2 {
				return m, nil
			}
		case key.Matches(msg, keymap.VimBinding.Exit):
			return m, tea.Quit
		}

	}
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return m, tea.Batch(cmds...)

}

func (m SignInModel) View() string {
	var b strings.Builder
	for _, t := range m.inputs {
		b.WriteString(t.View())
		b.WriteRune('\n')
	}
	logInChoice := textInputUnfocusedStyle.Render("Log In")
	if m.focusIndex == 2 {
		logInChoice = textInputFocusedStyle.Render("Log In")
	}

	b.WriteString(logInChoice)
	return b.String()
}
