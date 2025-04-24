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
	done     bool
	running  bool
	text     string
	duration time.Duration
}

func RestTimerModelInit(workTime time.Duration, ratio int) *RestTimerModel {
	return &RestTimerModel{
		done:     false,
		running:  false,
		duration: time.Duration(workTime.Nanoseconds() / int64(ratio)),
		text:     "Press select to start your break.",
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

	if m.running {
		durationTime := fmt.Sprintf("%02d:%02d:%02d", int(m.duration.Hours())%60, int(m.duration.Minutes())%60, int(m.duration.Seconds())%60)
		b.WriteString(durationTime)
		return b.String()
	}

	b.WriteString(m.text)
	return b.String()
}
