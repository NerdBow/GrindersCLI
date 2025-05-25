package model

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/NerdBow/GrindersTUI/internal/keymap"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	IdWidth       int = 3
	DateWidth     int = 10
	DurationWidth int = 8
	NameWidth     int = 10
	CategoryWidth int = 10
	GoalWidth     int = 20
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type RecentLogsModel struct {
	token    string
	page     int
	logTable table.Model
	status   string
	logs     []Log
}

func RecentLogsModelInit(token string) *RecentLogsModel {
	columns := []table.Column{
		{Title: "Id", Width: IdWidth},
		{Title: "Date", Width: DateWidth},
		{Title: "Duration", Width: DurationWidth},
		{Title: "Name", Width: NameWidth},
		{Title: "Category", Width: CategoryWidth},
		{Title: "Goal", Width: GoalWidth},
	}
	return &RecentLogsModel{
		token:    token,
		page:     1,
		logTable: table.New(table.WithColumns(columns)),
	}
}

func (m *RecentLogsModel) Init() tea.Cmd {
	m.logTable.SetHeight(20)
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
		case key.Matches(msg, keymap.VimBinding.Left):
			m.page--
			if m.page == 0 {
				m.page = 1
			}
			return m, m.getRecentLogs()
		case key.Matches(msg, keymap.VimBinding.Right):
			m.page++
			return m, m.getRecentLogs()
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
	case tea.WindowSizeMsg:
		log.Printf("H: %d W: %d", msg.Height, msg.Width)
	}
	return m, nil
}

func (m *RecentLogsModel) View() string {
	b := strings.Builder{}
	b.WriteString(baseStyle.Render(m.logTable.View()))
	b.WriteRune('\n')
	b.WriteString(fmt.Sprintf("Page: %d", m.page))
	return b.String()
}

func (m *RecentLogsModel) getRecentLogs() tea.Cmd {
	return func() tea.Msg {
		url := os.Getenv("URL")

		url += "/user/log"

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
