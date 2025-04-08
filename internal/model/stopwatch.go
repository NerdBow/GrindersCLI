package model

import (
	"fmt"
	"strings"

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

type StopwatchModel struct {
	sw         stopwatch.Model
	focusIndex int
}

func StopwatchModelInit() *StopwatchModel {
	return &StopwatchModel{
		sw:         stopwatch.New(),
		focusIndex: 0,
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
