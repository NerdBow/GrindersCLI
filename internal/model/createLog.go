package model

import (
	"github.com/NerdBow/GrindersTUI/internal/keymap"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type CreateLogModel struct {
}

func CreateLogModelInit() *CreateLogModel {
	return &CreateLogModel{}
}

func (m *CreateLogModel) Init() tea.Cmd {
	return nil
}

func (m *CreateLogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *CreateLogModel) View() string {
	return ""
}
