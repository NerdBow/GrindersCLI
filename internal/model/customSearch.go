package model

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/NerdBow/GrindersTUI/internal/keymap"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type CustomSearchModel struct {
	focusIndexRow int
	focusIndexCol int
	orderSettings []int
	inputs        []textinput.Model
	choices       [][]string
	querySettings GetLogsSettingsMsg
	status        string
}

func CustomSearchModelInit() *CustomSearchModel {
	textInputs := make([]textinput.Model, 3)
	textInputs[0] = textinput.New()
	textInputs[0].Placeholder = "Category"
	textInputs[0].TextStyle = textInputFocusedStyle
	textInputs[0].PromptStyle = textInputFocusedStyle

	textInputs[1] = textinput.New()
	textInputs[1].Placeholder = "Date Start"
	textInputs[1].TextStyle = textInputUnfocusedStyle
	textInputs[1].PromptStyle = textInputUnfocusedStyle

	textInputs[2] = textinput.New()
	textInputs[2].Placeholder = "Date End"
	textInputs[2].TextStyle = textInputUnfocusedStyle
	textInputs[2].PromptStyle = textInputUnfocusedStyle

	return &CustomSearchModel{
		focusIndexRow: 0,
		focusIndexCol: 0,
		orderSettings: []int{0, 0},
		inputs:        textInputs,
		choices:       [][]string{{"Date", "Duration"}, {"Descending", "Ascending"}},
		querySettings: GetLogsSettingsMsg{},
		status:        "",
	}
}

func (m *CustomSearchModel) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.inputs[0].Focus())
}

func (m *CustomSearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymap.VimBinding.Exit):
			return m, func() tea.Msg { return ModelMsg{CustomLogSearch, ViewLog, nil} }
		case key.Matches(msg, keymap.VimBinding.Select):
			err := m.writeSettings()
			if err != nil {
				m.status = err.Error()
				return m, nil
			}
			return m, func() tea.Msg { return ModelMsg{CustomLogSearch, RecentLogs, m.querySettings} }
		case key.Matches(msg, keymap.VimBinding.ChangeFocus):
			m.focusIndexRow = (m.focusIndexRow + 1) % (len(m.inputs) + len(m.choices))
			for i := range m.inputs {
				m.inputs[i].Blur()
				m.inputs[i].TextStyle = textInputUnfocusedStyle
				m.inputs[i].PromptStyle = textInputUnfocusedStyle
			}
			if m.focusIndexRow < len(m.inputs) {
				m.inputs[m.focusIndexRow].PromptStyle = textInputFocusedStyle
				m.inputs[m.focusIndexRow].TextStyle = textInputFocusedStyle
				return m, m.inputs[m.focusIndexRow].Focus()
			}
			if m.focusIndexRow >= len(m.inputs) {
				m.focusIndexCol = m.orderSettings[m.focusIndexRow-len(m.inputs)]
			}
		case key.Matches(msg, keymap.VimBinding.Right):
			if m.focusIndexRow < len(m.inputs) {
				break
			}
			m.focusIndexCol = (m.focusIndexCol + 1) % (len(m.choices[0]))
			m.orderSettings[m.focusIndexRow-len(m.inputs)] = m.focusIndexCol
		case key.Matches(msg, keymap.VimBinding.Left):
			if m.focusIndexRow < len(m.inputs) {
				break
			}
			m.focusIndexCol = (len(m.choices[0]) + m.focusIndexCol - 1) % (len(m.choices[0]))
			m.orderSettings[m.focusIndexRow-len(m.inputs)] = m.focusIndexCol
		}
	}
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m *CustomSearchModel) writeSettings() error {
	m.querySettings.Category = m.inputs[0].Value()

	if m.inputs[1].Value() != "" {
		dateStart, err := time.ParseInLocation(time.DateOnly, m.inputs[1].Value(), time.Now().Local().Location())
		if err != nil {
			return errors.New("Please input dates as YYYY-MM-DD")
		}
		m.querySettings.DateStart = dateStart.Unix()
	}
	if m.inputs[2].Value() != "" {
		dateEnd, err := time.ParseInLocation(time.DateOnly, m.inputs[2].Value(), time.Now().Local().Location())
		if err != nil {
			return errors.New("Please input dates as YYYY-MM-DD")
		}
		dateEnd = dateEnd.Add(time.Hour * 24)
		m.querySettings.DateEnd = dateEnd.Unix()
	}
	mapping := []string{"des", "asc"}

	m.querySettings.Order = fmt.Sprintf("%s_%s", m.choices[0][m.orderSettings[0]], mapping[m.orderSettings[1]])
	return nil
}

func (m *CustomSearchModel) View() string {
	b := strings.Builder{}
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		b.WriteByte('\n')
	}
	for i := range m.choices {
		if i == m.focusIndexRow-len(m.inputs) {
			b.WriteString(textInputFocusedStyle.Render("> "))
		} else {
			b.WriteString(textInputUnfocusedStyle.Render("> "))
		}
		for j := range m.choices[0] {
			if j == m.orderSettings[i] {
				b.WriteString(textInputFocusedStyle.Render(m.choices[i][j]))
			} else {
				b.WriteString(textInputUnfocusedStyle.Render(m.choices[i][j]))
			}
			b.WriteString("    ")
		}
		b.WriteByte('\n')
	}
	b.WriteByte('\n')
	b.WriteString(textInputFocusedStyle.Render(m.status))
	return b.String()
}
