package model

import (
	"fmt"
	"strings"

	"github.com/NerdBow/GrindersTUI/internal/keymap"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type SelectedLogModel struct {
	log Log
}

func SelectedLogModelInit(log Log) *SelectedLogModel {
	return &SelectedLogModel{log}
}

func (m *SelectedLogModel) Init() tea.Cmd {
	return nil
}

func (m *SelectedLogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *SelectedLogModel) View() string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("Id: %d", m.log.Id))
	b.WriteRune('\n')
	b.WriteString(fmt.Sprintf("Date: %s", m.log.DateString()))
	b.WriteRune('\n')
	b.WriteString(fmt.Sprintf("Name: %s", m.log.Name))
	b.WriteRune('\n')
	b.WriteString(fmt.Sprintf("Category: %s", m.log.Category))
	b.WriteRune('\n')
	b.WriteString(fmt.Sprintf("Goal: %s", m.log.Goal))
	b.WriteRune('\n')
	b.WriteString(fmt.Sprintf("Duration: %s", m.log.DurationString()))
	b.WriteRune('\n')
	return b.String()
}
