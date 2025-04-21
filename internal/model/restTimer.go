package model

import (
	"strings"
	"time"

	"github.com/NerdBow/GrindersTUI/internal/keymap"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type RestTimerModel struct {
	duration time.Duration
}

func (m *RestTimerModel) Init() tea.Cmd {
	return nil
}

func (m *RestTimerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (m *RestTimerModel) View() string {
	b := strings.Builder{}
	return ""
}
