package model

import (
	"fmt"
	"strings"
	"time"

	// "github.com/NerdBow/GrindersTUI/internal/keymap"
	// "github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type RestTimerModel struct {
	duration time.Duration
}

func RestTimerModelInit(workTime time.Duration, ratio int) *RestTimerModel {
	return &RestTimerModel{
		duration: time.Duration(float64(workTime.Nanoseconds()) / float64(ratio)),
	}
}

func (m *RestTimerModel) Init() tea.Cmd {
	return nil
}

func (m *RestTimerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (m *RestTimerModel) View() string {
	b := strings.Builder{}

	durationTime := fmt.Sprintf("%02d:%02d:%02d", int(m.duration.Hours())%60, int(m.duration.Minutes())%60, int(m.duration.Seconds())%60)
	b.WriteString(durationTime)

	return b.String()
}
