package keymap

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyBindings struct {
	Up          key.Binding
	Down        key.Binding
	Left        key.Binding
	Right       key.Binding
	Select      key.Binding
	Exit        key.Binding
	ChangeFocus key.Binding
}

var (
	VimBinding = KeyBindings{
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("↑/k", "move up"),
		),

		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("↓/j", "move down"),
		),

		Left: key.NewBinding(
			key.WithKeys("h", "left"),
			key.WithHelp("←/h", "move left"),
		),

		Right: key.NewBinding(
			key.WithKeys("l", "right"),
			key.WithHelp("→/l", "move right"),
		),

		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),

		Exit: key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("esc/ctrl+c", "quit"),
		),

		ChangeFocus: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "change focused section"),
		),
	}
)
