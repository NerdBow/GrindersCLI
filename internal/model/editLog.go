package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/NerdBow/GrindersTUI/internal/keymap"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ConfirmEdit bool

type EditStatusMsg struct {
	Result bool `json:"result"`
}

type EditLogModel struct {
	log          Log
	choices      []string
	focusIndex   int
	inputs       [5]textinput.Model
	token        string
	status       string
	confirmCount uint8
}

func EditLogModelInit(log Log, token string) *EditLogModel {
	var textInputs [5]textinput.Model
	for i := range textInputs {
		textInputs[i] = textinput.New()
	}
	return &EditLogModel{
		log:          log,
		choices:      []string{"Date", "Name", "Category", "Goal", "Duration", "Back", "Confirm"},
		focusIndex:   0,
		inputs:       textInputs,
		token:        token,
		status:       "",
		confirmCount: 0,
	}
}

func (m *EditLogModel) Init() tea.Cmd {
	m.inputs[0].PromptStyle = textInputFocusedStyle
	m.inputs[0].TextStyle = textInputFocusedStyle
	m.SyncTextInputs()
	return nil
}

func (m *EditLogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *EditLogModel) View() string {
	b := strings.Builder{}
	for i := range m.inputs {
		b.WriteString(fmt.Sprintf("%s: %s", m.choices[i], m.inputs[i].View()))
		b.WriteRune('\n')
	}
	b.WriteRune('\n')
	if m.choices[m.focusIndex] == "Back" {
		b.WriteString(textInputFocusedStyle.Render("Back"))
	} else {
		b.WriteString(textInputUnfocusedStyle.Render("Back"))
	}

	b.WriteString("    ")
	if m.choices[m.focusIndex] == "Confirm" {
		b.WriteString(textInputFocusedStyle.Render("Confirm"))
	} else {
		b.WriteString(textInputUnfocusedStyle.Render("Confirm"))
	}

	b.WriteString("    ")
	b.WriteRune('\n')
	b.WriteRune('\n')
	b.WriteString(textInputFocusedStyle.Render(m.status))
	return b.String()
}

func (m *EditLogModel) EditLog(logId int, editLog Log) tea.Cmd {
	return func() tea.Msg {
		url := os.Getenv("URL")
		url += "/user/log"

		jsonBytes, err := json.Marshal(editLog)
		if err != nil {
			return SystemErrorMsg(err.Error())
		}

		req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonBytes))
		if err != nil {
			return SystemErrorMsg(err.Error())
		}

		req.Header.Add("Authorization", "Bearer "+m.token)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return SystemErrorMsg(err.Error())
		}

		switch res.StatusCode {
		case http.StatusOK:
			msg := EditStatusMsg{}
			jsonBytes, err := io.ReadAll(res.Body)
			if err != nil {
				return SystemErrorMsg(err.Error())
			}

			err = json.Unmarshal(jsonBytes, &msg)
			if err != nil {
				return SystemErrorMsg(err.Error())
			}
			return msg
		default:
			msg := PostLogErrorMsg{}
			jsonBytes, err := io.ReadAll(res.Body)
			if err != nil {
				return SystemErrorMsg(err.Error())
			}

			err = json.Unmarshal(jsonBytes, &msg)
			if err != nil {
				return SystemErrorMsg(err.Error())
			}
			return msg
		}
	}
}
