package model

import (
	"fmt"
	"strings"

	"github.com/NerdBow/GrindersTUI/internal/keymap"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CustomSearchModel struct {
	focusIndexRow int
	focusIndexCol int
	inputs        []textinput.Model
	choices       [][]string
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
		inputs:        textInputs,
		choices:       [][]string{{"Date", "Duration"}, {"Descending", "Ascending"}},
	}
}

func (m *CustomSearchModel) Init() tea.Cmd {
	return nil
}

func (m *CustomSearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *CustomSearchModel) View() string {
	b := strings.Builder{}
	return b.String()
}
