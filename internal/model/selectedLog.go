package model

import (
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

const (
	BackField = iota
	EditField
	DeleteField
	TextField
)

type ConfirmDeletion bool

type DeletionStatusMsg struct {
	Result bool `json:"result"`
}

type SelectedLogModel struct {
	log           Log
	previousModel int
	choices       []string
	focusIndex    int
	textField     textinput.Model
	token         string
	status        string
}

func SelectedLogModelInit(log Log, previousModel int, token string) *SelectedLogModel {
	return &SelectedLogModel{
		log:           log,
		previousModel: previousModel,
		choices:       []string{"Back", "Edit", "Delete"},
		focusIndex:    0,
		textField:     textinput.New(),
		token:         token,
		status:        "",
	}
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
			m.textField.Blur()
			m.focusIndex = (m.focusIndex + 1) % len(m.choices)
		case key.Matches(msg, keymap.VimBinding.Select):
			switch m.focusIndex {
			case BackField:
				return m, func() tea.Msg { return ModelMsg{SelectedLog, m.previousModel, nil} }
			case EditField:
				return m, func() tea.Msg { return ModelMsg{SelectedLog, EditLog, m.log} }
			case DeleteField:
				m.textField.Placeholder = "Type the log id of this log in order to delete it."
				m.focusIndex = TextField
				return m, m.textField.Focus()
			case TextField:
				return m, m.CheckTypedId(m.textField.Value(), m.log.Id)
			}
		}
	case DeletionStatusMsg:
		m.status = fmt.Sprintf("Log %d has been deleted", m.log.Id)
	case SystemErrorMsg:
		m.status = string(msg)
	case PostLogErrorMsg:
		m.status = msg.Message
	case ConfirmDeletion:
		if msg {
			return m, m.DeleteLog(m.log.Id)
		}
		m.status = "Invalid log id"
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
	b.WriteRune('\n')
	for i := range m.choices {
		if i != m.focusIndex {
			b.WriteString(textInputUnfocusedStyle.Render(m.choices[i]))
		} else {
			b.WriteString(textInputFocusedStyle.Render(m.choices[i]))
		}
		b.WriteString("   ")
	}
	b.WriteRune('\n')
	b.WriteString(textInputFocusedStyle.Render(m.status))
	b.WriteRune('\n')
	b.WriteString(m.textField.View())
	return b.String()
}

func (m *SelectedLogModel) CheckTypedId(typedId string, logId int64) tea.Cmd {
	return func() tea.Msg {
		id, err := strconv.ParseInt(typedId, 10, 64)
		if err != nil {
			return ConfirmDeletion(false)
		}
		if id != logId {
			return ConfirmDeletion(false)
		}
		return ConfirmDeletion(true)
	}
}

func (m *SelectedLogModel) DeleteLog(logId int64) tea.Cmd {
	return func() tea.Msg {
		url := os.Getenv("URL")
		url += "/user/log"
		url = fmt.Sprintf("%s?id=%d", url, logId)
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		req.Header.Add("Authorization", "Bearer "+m.token)
		if err != nil {
			return SystemErrorMsg(err.Error())
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return SystemErrorMsg(err.Error())
		}

		switch res.StatusCode {
		case http.StatusOK:
			msg := DeletionStatusMsg{}
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
