package model

import (
	"fmt"
	"strconv"
	"time"

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

func (l Log) ToStringArray() []string {
	return []string{strconv.Itoa(l.Id), l.DateString(), l.DurationString(), l.Name, l.Category, l.Goal}
}

func (l Log) DurationString() string {
	duration := time.Second * time.Duration(l.Duration)
	return fmt.Sprintf("%02d:%02d:%02d", int(duration.Hours())%60, int(duration.Minutes())%60, int(duration.Seconds())%60)
}

func (l Log) DateString() string {
	return time.Unix(int64(l.Date), 0).Format("2006-01-02")
}
