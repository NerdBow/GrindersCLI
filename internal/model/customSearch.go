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
	focusIndex int
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
