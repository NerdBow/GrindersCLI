package main

import (
	"fmt"
	"os"

	"github.com/NerdBow/GrindersTUI/internal/model"
	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	currentState tea.Model
}

func initApp() App {
	return App{currentState: model.HomeModelInit()}
}

func (m App) Init() tea.Cmd {
	return nil
}

func (m App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := m.currentState.Update(msg)
	return model, cmd
}

func (m App) View() string {
	return m.currentState.View()
}

func main() {
	p := tea.NewProgram(initApp(), tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		fmt.Printf("There is an error: %v", err)
		os.Exit(1)
	}

}
