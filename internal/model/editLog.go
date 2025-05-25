package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

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

func (m *EditLogModel) SyncTextInputs() {
	m.inputs[0].SetValue(m.log.DateString())
	m.inputs[1].SetValue(m.log.Name)
	m.inputs[2].SetValue(m.log.Category)
	m.inputs[3].SetValue(m.log.Goal)
	m.inputs[4].SetValue(m.log.DurationString())
}

// GetEditLog takes the changed values of the textinputs and creates a logs of them.
// If any error occurs then a non-empty string will be returned.
func (m *EditLogModel) GetEditLog() (Log, string) {
	changed := false
	l := Log{}
	l.Id = m.log.Id

	if m.log.Name != m.inputs[1].Value() {
		l.Name = m.inputs[1].Value()
		changed = true
	}

	if m.log.Category != m.inputs[2].Value() {
		l.Category = m.inputs[2].Value()
		changed = true
	}

	if m.log.Goal != m.inputs[3].Value() {
		l.Goal = m.inputs[3].Value()
		changed = true
	}

	if m.log.DateString() != m.inputs[0].Value() {
		t, err := time.ParseInLocation(time.DateOnly, m.inputs[0].Value(), time.Now().Local().Location())
		if err != nil {
			return Log{}, "Please format your date as YYYY-MM-DD"
		}
		l.Date = t.Unix()
		changed = true
	}

	if m.log.DurationString() != m.inputs[4].Value() {
		unitDurations := strings.Split(m.inputs[4].Value(), ":")
		if len(unitDurations) != 3 {
			return Log{}, "Please format your duration as HH:MM:SS"
		}
		d, err := time.ParseDuration(fmt.Sprintf("%sh%sm%ss", unitDurations[0], unitDurations[1], unitDurations[2]))
		if err != nil {
			return Log{}, "Please format your duration as HH:MM:SS"
		}
		l.Duration = int64(d.Seconds())
		changed = true
	}

	if !changed {
		return Log{}, ""
	}

	return l, ""
}

func (m *EditLogModel) EditLog(editLog Log) tea.Cmd {
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
