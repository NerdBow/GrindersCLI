package model

import (
	"fmt"
	"strings"

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
	}
}

func (m *CustomSearchModel) Init() tea.Cmd {
	return nil
}

func (m *CustomSearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymap.VimBinding.Exit):
			return m, func() tea.Msg { return ModelMsg{CustomLogSearch, ViewLog, nil} }
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
	b.WriteString(fmt.Sprintf("Row: %d Col: %d ", m.focusIndexRow, m.focusIndexCol))
	return b.String()
}
