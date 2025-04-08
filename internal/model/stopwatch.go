package model

import (
	"github.com/NerdBow/GrindersTUI/internal/keymap"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type StopwatchModel struct {
	sw stopwatch.Model
}

func (m *StopwatchModel) Init() tea.Cmd {
	return nil
}

func (m *StopwatchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *StopwatchModel) View() string {
	return ""
}
