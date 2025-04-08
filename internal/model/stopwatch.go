package model

import (
	"github.com/NerdBow/GrindersTUI/internal/keymap"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"strings"

const (
	TimerField = iota
	RestField
	FinishLogField
)

type StopwatchModel struct {
	sw         stopwatch.Model
	focusIndex int
}
}

func (m *StopwatchModel) Init() tea.Cmd {
	return nil
}

func (m *StopwatchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
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
