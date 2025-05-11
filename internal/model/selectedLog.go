package model

import (
	"github.com/NerdBow/GrindersTUI/internal/keymap"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
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
	return b.String()
}
