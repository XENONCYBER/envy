package tui

import (
	"strings"

	"envy/internal/config"

	"github.com/charmbracelet/lipgloss"
)

type BottomBar struct {
	width    int
	state    SessionState
	keymap   config.KeyMap
	bindings []KeyBinding
	styles   config.Styles
}

type KeyBinding struct {
	Key         string
	Description string
}

func NewBottomBar(width int, state SessionState, keymap config.KeyMap, bindings []KeyBinding, styles config.Styles) BottomBar {
	return BottomBar{
		width:    width,
		state:    state,
		keymap:   keymap,
		bindings: bindings,
		styles:   styles,
	}
}

func (b BottomBar) Render() string {
	modeText := " NORMAL "
	modeColor := b.styles.Success
	if b.state == StateInsert {
		modeText = " INSERT "
		modeColor = b.styles.Warning
	}

	modeStyle := lipgloss.NewStyle().
		Foreground(b.styles.Base).
		Background(modeColor).
		Bold(true)

	modeIndicator := modeStyle.Render(modeText)

	var shortcuts []string
	separatorText := "  │  "

	for _, binding := range b.bindings {
		keyStyle := lipgloss.NewStyle().
			Foreground(b.styles.Accent).
			Bold(true)

		descStyle := lipgloss.NewStyle().
			Foreground(b.styles.Text)

		shortcut := descStyle.Render(binding.Description+": ") + keyStyle.Render(binding.Key)
		shortcuts = append(shortcuts, shortcut)
	}

	shortcutsText := strings.Join(shortcuts, separatorText)

	content := " " + modeIndicator + separatorText + shortcutsText

	contentWidth := lipgloss.Width(content)
	remainingSpace := b.width - contentWidth
	if remainingSpace < 0 {
		remainingSpace = 0
	}
	padding := strings.Repeat(" ", remainingSpace)

	barStyle := lipgloss.NewStyle().
		Foreground(b.styles.Text).
		Background(b.styles.Surface0)

	return barStyle.Render(content + padding)
}

func GridViewBindings(km config.KeyMap, state SessionState) []KeyBinding {
	if state == StateInsert {
		return []KeyBinding{
			{Key: km.Tab, Description: "Mode"},
			{Key: km.Back, Description: "Exit"},
		}
	}

	return []KeyBinding{
		{Key: km.Enter, Description: "Open"},
		{Key: km.Search, Description: "Search"},
		{Key: km.Create, Description: "New"},
		{Key: km.Delete, Description: "Delete"},
		{Key: km.Quit, Description: "Quit"},
	}
}

func DetailViewBindings(km config.KeyMap) []KeyBinding {
	return []KeyBinding{
		{Key: km.Enter, Description: "Reveal"},
		{Key: km.Yank, Description: "Copy"},
		{Key: km.Edit, Description: "Edit Key"},
		{Key: km.EditProject, Description: "Edit Project"},
		{Key: km.History, Description: "History"},
		{Key: km.Delete, Description: "Delete"},
		{Key: km.Back, Description: "Back"},
	}
}

func CreateViewBindings(km config.KeyMap, state SessionState) []KeyBinding {
	if state == StateInsert {
		return []KeyBinding{
			{Key: km.Back, Description: "Exit"},
		}
	}

	return []KeyBinding{
		{Key: km.Search, Description: "Edit"},
		{Key: km.Add, Description: "Add Key"},
		{Key: km.Save, Description: "Save"},
		{Key: km.Quit, Description: "Cancel"},
	}
}

func EditViewBindings(km config.KeyMap) []KeyBinding {
	return []KeyBinding{
		{Key: km.Enter, Description: "Save"},
		{Key: km.Back, Description: "Cancel"},
	}
}

func EditProjectViewBindings(km config.KeyMap, state SessionState) []KeyBinding {
	if state == StateInsert {
		return []KeyBinding{
			{Key: km.Back, Description: "Exit"},
		}
	}

	return []KeyBinding{
		{Key: km.Search, Description: "Edit"},
		{Key: km.Add, Description: "Add Key"},
		{Key: km.Delete, Description: "Delete Key"},
		{Key: km.Save, Description: "Save"},
		{Key: km.Back, Description: "Cancel"},
	}
}

func HistoryViewBindings(km config.KeyMap) []KeyBinding {
	return []KeyBinding{
		{Key: "j/k", Description: "Navigate Keys"},
		{Key: km.Back, Description: "Close"},
	}
}
