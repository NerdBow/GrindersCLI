package model

import (
	tea "github.com/charmbracelet/bubbletea"
)

const (
	SignIn = iota
	Home
	CreateLog
	ViewLog
	EditLog
	DeleteLog
	Stopwatch
	RestTimer
)

type ModelMsg struct {
	CurrentModel int
	NextModel    int
	Other        tea.Msg
}
