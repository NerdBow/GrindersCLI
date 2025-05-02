package main

import (
	"fmt"
	"os"

	"github.com/NerdBow/GrindersTUI/internal/model"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	WORKTORESTRATIO int = 5
)

type App struct {
	currentState   tea.Model
	homeModel      *model.HomeModel
	signInModel    *model.SignInModel
	createLogModel *model.CreateLogModel
	stopwatchModel *model.StopwatchModel
	restTimerModel *model.RestTimerModel
	viewLogModel   *model.ViewLogModel
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
	case model.ModelMsg:
		switch msg.CurrentModel {
		case model.SignIn:
			switch other := msg.Other.(type) {
			case model.UserTokenMsg:
				m.token = other.Token
			}

			switch msg.NextModel {
			case model.Home:
				m.signInModel = model.SignInModelInit()
				m.currentState = m.homeModel
			}
		case model.Home:
			switch msg.Other.(type) {
			case model.SignOutMsg:
				m.token = ""
			}

			switch msg.NextModel {
			case model.CreateLog:
				m.createLogModel = model.CreateLogModelInit()
				m.currentState = m.createLogModel
			case model.ViewLog:
				m.viewLogModel = model.ViewLogModelInit()
				m.currentState = m.viewLogModel
			case model.EditLog:
			case model.DeleteLog:
			case model.SignIn:
				m.currentState = m.signInModel
			}
		case model.CreateLog:
			switch msg.NextModel {
			case model.Home:
				m.homeModel = model.HomeModelInit()
				m.currentState = m.homeModel
			case model.Stopwatch:
				name, category, goal := m.createLogModel.GetLogInfo()
				m.stopwatchModel = model.StopwatchModelInit(name, category, goal, m.token)
				m.currentState = m.stopwatchModel
			}
		case model.ViewLog:
			switch msg.NextModel {
			case model.Home:
				m.homeModel = model.HomeModelInit()
				m.currentState = m.homeModel
			}
		case model.EditLog:
		case model.DeleteLog:
		case model.Stopwatch:
			switch msg.NextModel {
			case model.CreateLog:
				m.stopwatchModel = model.StopwatchModelInit("", "", "", "")
				m.currentState = m.createLogModel
			case model.RestTimer:
				m.restTimerModel = model.RestTimerModelInit(m.stopwatchModel.GetWorkTime(), WORKTORESTRATIO)
				m.currentState = m.restTimerModel
			}
		case model.RestTimer:
			switch msg.NextModel {
			case model.Stopwatch:
				m.currentState = m.stopwatchModel
			}
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
