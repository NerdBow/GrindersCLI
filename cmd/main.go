package main

import (
	"fmt"
	"os"

	"github.com/NerdBow/GrindersTUI/internal/model"
	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	currentState   tea.Model
	homeModel      *model.HomeModel
	signInModel    *model.SignInModel
	createLogModel *model.CreateLogModel
	stopwatchModel *model.StopwatchModel
	token          string
}

func initApp() *App {
	return &App{
		homeModel:   model.HomeModelInit(),
		signInModel: model.SignInModelInit(),
	}
}

func (m *App) Init() tea.Cmd {
	m.currentState = m.signInModel
	return nil
}

func (m *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case model.UserTokenMsg:
		m.token = msg.Token
		m.signInModel = model.SignInModelInit()
		m.currentState = m.homeModel
	case model.SignOutMsg:
		m.token = ""
		m.homeModel = model.HomeModelInit()
		m.currentState = m.signInModel
	case model.HomeModelSwitch:
		switch msg {
		case model.CreateLogSwitch:
			m.createLogModel = model.CreateLogModelInit()
			m.currentState = m.createLogModel
		case model.ViewLogSwitch:
		case model.EditLogSwitch:
		case model.DeleteLogSwitch:
		}
	case model.CreateLogModelSwitch:
		switch msg {
		case model.HomeSwitch:
			m.homeModel = model.HomeModelInit()
			m.currentState = m.homeModel
		case model.TimerSwitch:
			m.stopwatchModel = model.StopwatchModelInit()
			m.currentState = m.stopwatchModel
		}

	}
	_, cmd := m.currentState.Update(msg)
	return m, cmd
}

func (m *App) View() string {
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
