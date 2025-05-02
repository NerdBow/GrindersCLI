package model

import (
	"github.com/NerdBow/GrindersTUI/internal/keymap"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

const ()

type RecentLogsModel struct {
	logTable   table.Model
	focusIndex int
}
