package model

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/NerdBow/GrindersTUI/internal/keymap"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const ()

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type RecentLogsModel struct {
	token      string
	page       int
	logTable   table.Model
	focusIndex int
	status     string
	logs       []Log
}

func RecentLogsModelInit(token string) *RecentLogsModel {
	columns := []table.Column{
		{Title: "Id", Width: 10},
		{Title: "Date", Width: 10},
		{Title: "Duration", Width: 8},
		{Title: "Name", Width: 25},
		{Title: "Category", Width: 25},
		{Title: "Goal", Width: 75},
	}
	return &RecentLogsModel{
		token:      token,
		page:       1,
		logTable:   table.New(table.WithColumns(columns)),
		focusIndex: 0,
	}
}

func (m *RecentLogsModel) Init() tea.Cmd {
	return m.getRecentLogs()
}

func (m *RecentLogsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymap.VimBinding.Exit):
			return m, func() tea.Msg { return ModelMsg{RecentLogs, ViewLog, nil} }
		case key.Matches(msg, keymap.VimBinding.Select):
			return m, func() tea.Msg { return ModelMsg{RecentLogs, SelectedLog, m.logs[m.logTable.Cursor()]} }
		case key.Matches(msg, keymap.VimBinding.Up):
			m.logTable.MoveUp(1)
			return m, nil
		case key.Matches(msg, keymap.VimBinding.Down):
			m.logTable.MoveDown(1)
			return m, nil
		}
	case []Log:
		m.logs = msg
		rows := make([]table.Row, 0, 20)
		for _, log := range m.logs {
			rows = append(rows, table.Row(log.ToStringArray()))
		}
		m.logTable.SetRows(rows)

	case SystemErrorMsg:
	case PostLogErrorMsg:
	}
	return m, nil
}

func (m *RecentLogsModel) View() string {
	return baseStyle.Render(m.logTable.View()) + "\n"
}

func (m *RecentLogsModel) getRecentLogs() tea.Cmd {
	return func() tea.Msg {
		url := os.Getenv("URL")

		url = "http://localhost:8080/user/log"

		req, err := http.NewRequest("GET", fmt.Sprintf("%s?order=DATE_DES&page=%d", url, m.page), nil)

		if err != nil {
			return SystemErrorMsg(err.Error())
		}

		req.Header.Add("Authorization", "Bearer "+m.token)

		res, err := http.DefaultClient.Do(req)

		if err != nil {
			return SystemErrorMsg(err.Error())
		}

		logs := make([]Log, 0, 20)

		switch res.StatusCode {
		case http.StatusOK:
			jsonBytes, err := io.ReadAll(res.Body)
			if err != nil {
				return SystemErrorMsg(err.Error())
			}
			err = json.Unmarshal(jsonBytes, &logs)

			if err != nil {
				return SystemErrorMsg(err.Error())
			}
			return logs
		default:
			jsonBytes, err := io.ReadAll(res.Body)
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
