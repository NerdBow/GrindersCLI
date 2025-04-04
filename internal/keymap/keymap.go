package keymap

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyBindings struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Exit   key.Binding
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

		Select: key.NewBinding(
			key.WithKeys("space", "enter"),
			key.WithHelp("space/enter", "select"),
		),

		Exit: key.NewBinding(
			key.WithKeys("esc", "q", "ctrl+c"),
			key.WithHelp("esc/q/ctrl+c", "quit"),
		),
	}
)
