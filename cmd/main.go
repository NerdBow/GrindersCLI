package main

import (
	"fmt"
	"os"

	"github.com/NerdBow/GrindersTUI/internal/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

const (
	WORKTORESTRATIO int = 5
)

type App struct {
	currentState      tea.Model
	homeModel         *model.HomeModel
	signInModel       *model.SignInModel
	createLogModel    *model.CreateLogModel
	stopwatchModel    *model.StopwatchModel
	restTimerModel    *model.RestTimerModel
	viewLogModel      *model.ViewLogModel
	selectedLogModel  *model.SelectedLogModel
	editLogModel      *model.EditLogModel
	recentLogsModel   *model.RecentLogsModel
	customSearchModel *model.CustomSearchModel
	token             string
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
			case model.RecentLogs:
				m.recentLogsModel = model.RecentLogsModelInit(m.token)
				m.currentState = m.recentLogsModel
				return m, m.recentLogsModel.Init()
			case model.CustomLogSearch:
				m.customSearchModel = model.CustomSearchModelInit()
				m.currentState = m.customSearchModel
				return m, m.currentState.Init()
			}
		case model.CustomLogSearch:
			switch msg.NextModel {
			case model.ViewLog:
				m.currentState = m.viewLogModel
			}
		case model.RecentLogs:
			switch msg.NextModel {
			case model.ViewLog:
				m.currentState = m.viewLogModel
			case model.SelectedLog:
				switch other := msg.Other.(type) {
				case model.Log:
					m.selectedLogModel = model.SelectedLogModelInit(other, model.RecentLogs, m.token)
				}
				m.currentState = m.selectedLogModel
			}
		case model.SelectedLog:
			switch msg.NextModel {
			case model.RecentLogs:
				m.currentState = m.recentLogsModel
				return m, m.recentLogsModel.Init()
			case model.EditLog:
				switch other := msg.Other.(type) {
				case model.Log:
					m.editLogModel = model.EditLogModelInit(other, m.token)
				}
				m.currentState = m.editLogModel
				return m, m.editLogModel.Init()
			}
		case model.EditLog:
			switch msg.NextModel {
			case model.SelectedLog:
				switch other := msg.Other.(type) {
				case model.Log:
					m.selectedLogModel = model.SelectedLogModelInit(other, model.RecentLogs, m.token)
				}
				m.currentState = m.selectedLogModel
				return m, m.selectedLogModel.Init()
			}
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
	case tea.WindowSizeMsg:
		// TODO: have all the models take in a resize function so the dimension can be sent in from here
		// m.currentState.resize(msg.Height, msg.Width)
	}
	_, cmd := m.currentState.Update(msg)
	return m, cmd
}

func (m *App) View() string {
	return m.currentState.View()
}

func main() {
	if os.Getenv("URL") == "" {
		err := godotenv.Load(".env")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
	}

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	p := tea.NewProgram(initApp(), tea.WithAltScreen())
	_, err = p.Run()
	if err != nil {
		fmt.Printf("There is an error: %v", err)
		os.Exit(1)
	}
}
