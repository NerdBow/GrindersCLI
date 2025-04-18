package model

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
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
	UsernameField           = 0
	PasswordField           = 1
	SignInButton            = 2
)

type SignInModel struct {
	inputs       []textinput.Model
	focusIndex   int
	errorMessage string
}

type UserTokenMsg struct {
	Token string `json:"token"`
}

type SignInErrorMsg struct {
	Message string `json:"message"`
}

type SystemErrorMsg string

func SignInModelInit() *SignInModel {
	m := SignInModel{
		inputs: make([]textinput.Model, 2),
	}
	usernameTextInput := textinput.New()
	usernameTextInput.CharLimit = 64
	usernameTextInput.Placeholder = "Username"
	usernameTextInput.PromptStyle = textInputFocusedStyle
	usernameTextInput.TextStyle = textInputFocusedStyle
	usernameTextInput.Focus()

	m.inputs[UsernameField] = usernameTextInput

	passwordTextInput := textinput.New()
	passwordTextInput.CharLimit = 64
	passwordTextInput.Placeholder = "Password"
	passwordTextInput.EchoMode = textinput.EchoPassword
	passwordTextInput.EchoCharacter = '•'

	m.inputs[PasswordField] = passwordTextInput

	return &m
}

func (m *SignInModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *SignInModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.focusIndex == SignInButton {
				return m, m.GetToken(m.inputs[UsernameField].Value(), m.inputs[PasswordField].Value())
			}
		case key.Matches(msg, keymap.VimBinding.Exit):
			return m, tea.Quit
		}
	case SystemErrorMsg:
		m.errorMessage = string(msg)
		return m, nil
	case SignInErrorMsg:
		m.errorMessage = msg.Message
		return m, nil
	}
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return m, tea.Batch(cmds...)

}

func (m *SignInModel) View() string {
	var b strings.Builder
	for _, t := range m.inputs {
		b.WriteString(t.View())
		b.WriteRune('\n')
	}
	logInChoice := textInputUnfocusedStyle.Render("Sign In")
	if m.focusIndex == SignInButton {
		logInChoice = textInputFocusedStyle.Render("Sign In")
	}
	b.WriteString(logInChoice)

	b.WriteRune('\n')
	b.WriteRune('\n')
	b.WriteString(textInputFocusedStyle.Render(m.errorMessage))

	return b.String()
}

func (m *SignInModel) GetToken(username string, password string) tea.Cmd {
	return func() tea.Msg {
		url := os.Getenv("URL")
		url = "http://localhost:8080/user/signin"
		request := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{username, password}

		jsonBytes, err := json.Marshal(request)

		if err != nil {
			return SystemErrorMsg(err.Error())
		}

		r, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBytes))

		if err != nil {
			return SystemErrorMsg(err.Error())
		}

		switch r.StatusCode {
		case http.StatusOK:
			jsonBytes, err = io.ReadAll(r.Body)
			if err != nil {
				return SystemErrorMsg(err.Error())
			}

			var msg UserTokenMsg

			err = json.Unmarshal(jsonBytes, &msg)

			if err != nil {
				return SystemErrorMsg(err.Error())
			}

			return ModelMsg{
				SignIn,
				Home,
				msg,
			}

		default:
			jsonBytes, err = io.ReadAll(r.Body)
			if err != nil {
				return SystemErrorMsg(err.Error())
			}

			var errMsg SignInErrorMsg

			err = json.Unmarshal(jsonBytes, &errMsg)

			if err != nil {
				return SystemErrorMsg(err.Error())
			}

			return errMsg
		}
	}
}
