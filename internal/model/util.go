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
	RecentLogs
	IdLogSearch
	CustomLogSearch
)

type ModelMsg struct {
	CurrentModel int
	NextModel    int
	Other        tea.Msg
}

type Log struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Goal     string `json:"goal"`
	Date     int    `json:"date"`
	Duration int    `json:"duration"`
	UserId   int    `json:"userId"`
}

