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
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	TimerField = iota
	RestField
	FinishLogField
)


type PostLogErrorMsg struct {
	Message string `json:"message"`
}
type LogIdMsg struct {
	Id int64 `json:"id"`
}
type StopwatchModel struct {
	sw          stopwatch.Model
	focusIndex  int
	logName     string
	logCategory string
	logGoal     string
	token       string
	status      string
}

func StopwatchModelInit(logName string, logCategory string, logGoal string, token string) *StopwatchModel {
	return &StopwatchModel{
		sw:          stopwatch.New(),
		focusIndex:  0,
		logName:     logName,
		logCategory: logCategory,
		logGoal:     logGoal,
		token:       token,
	}
}

func (m *StopwatchModel) Init() tea.Cmd {
	return nil
}

func (m *StopwatchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymap.VimBinding.ChangeFocus):
			m.focusIndex = (m.focusIndex + 1) % 3
		case key.Matches(msg, keymap.VimBinding.Select):
			switch m.focusIndex {
			case TimerField:
				cmd := m.sw.Toggle()
				return m, cmd
			case RestField:
				return m, nil //TODO: this is for later
			case FinishLogField:
				cmds := make([]tea.Cmd, 2)
				cmds[0] = m.sw.Stop()
				return m, tea.Batch(cmds...)
			}
		case key.Matches(msg, keymap.VimBinding.Exit):
			return nil, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.sw, cmd = m.sw.Update(msg)
	return m, cmd
}

func (m *StopwatchModel) View() string {
	b := strings.Builder{}
	durationTime := fmt.Sprintf("%02d:%02d:%02d", int(m.sw.Elapsed().Hours())%60, int(m.sw.Elapsed().Minutes())%60, int(m.sw.Elapsed().Seconds())%60)
	renderedTime := textInputUnfocusedStyle.Render(durationTime)

	if m.focusIndex == TimerField {
		renderedTime = textInputFocusedStyle.Render(durationTime)
	}
	b.WriteString(renderedTime)

	b.WriteRune('\n')

	restText := textInputUnfocusedStyle.Render("Rest")

	if m.focusIndex == RestField {
		restText = textInputFocusedStyle.Render("Rest")
	}

	b.WriteString(restText)

	b.WriteRune('\n')

	confirmText := textInputUnfocusedStyle.Render("Finish Log")
	if m.focusIndex == FinishLogField {
		confirmText = textInputFocusedStyle.Render("Finish Log")
	}
	b.WriteString(confirmText)
	return b.String()
}

func (m *StopwatchModel) postLog() tea.Cmd {
	return func() tea.Msg {
		url := os.Getenv("URL")

		url = "http://localhost:8080/user/log"

		log := struct {
			Date     int64  `json:"date"`
			Duration int64  `json:"duration"`
			Name     string `json:"name"`
			Category string `json:"category"`
			Goal     string `json:"goal"`
		}{
			time.Now().Unix(),
			int64(m.sw.Elapsed().Seconds()),
			m.logName,
			m.logCategory,
			m.logGoal,
		}

		jsonBytes, err := json.Marshal(log)

		if err != nil {
			return SystemErrorMsg(err.Error())
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
		req.Header.Add("Authorization", "Bearer "+m.token)
		res, err := http.DefaultClient.Do(req)

		if err != nil {
			return SystemErrorMsg(err.Error())
		}

		var msg LogIdMsg

		switch res.StatusCode {
		case http.StatusOK:
			jsonBytes, err = io.ReadAll(res.Body)
			if err != nil {
				return SystemErrorMsg(err.Error())
			}
			err = json.Unmarshal(jsonBytes, &msg)

			if err != nil {
				return SystemErrorMsg(err.Error())
			}
			return msg
		default:
			jsonBytes, err = io.ReadAll(res.Body)
			if err != nil {
				return SystemErrorMsg(err.Error())
			}

			var errMsg PostLogErrorMsg
			err = json.Unmarshal(jsonBytes, &errMsg)

			if err != nil {
				return SystemErrorMsg(err.Error())
			}

			return errMsg
		}
	}
}
