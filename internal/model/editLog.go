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

type EditLogModel struct{}

func (m *EditLogModel) Init() tea.Cmd {
	return nil
}

func (m *EditLogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *EditLogModel) View() string {
	b := strings.Builder{}
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
