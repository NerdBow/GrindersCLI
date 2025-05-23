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
	SelectedLog
	IdLogSearch
	CustomLogSearch
)

type ModelMsg struct {
	CurrentModel int
	NextModel    int
	Other        tea.Msg
}

type Log struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Goal     string `json:"goal"`
	Date     int64  `json:"date"`
	Duration int64  `json:"duration"`
	UserId   int    `json:"userId"`
}

func (l Log) ToStringArray() []string {
	return []string{strconv.FormatInt(l.Id, 10), l.DateString(), l.DurationString(), l.Name, l.Category, l.Goal}
}

func (l Log) DurationString() string {
	duration := time.Second * time.Duration(l.Duration)
	return fmt.Sprintf("%02d:%02d:%02d", int(duration.Hours())%60, int(duration.Minutes())%60, int(duration.Seconds())%60)
}

func (l Log) DateString() string {
	return time.Unix(int64(l.Date), 0).Format("2006-01-02")
}

func (l Log) IsEmpty() bool {
	return (l.Id == 0 && l.Name == "" && l.Category == "" && l.Goal == "" && l.Date == 0 && l.Duration == 0 && l.UserId == 0)
}
