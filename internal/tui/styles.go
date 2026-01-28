// Package tui - styles.go
// DEPRECATED: Global styles are deprecated. Use config.Styles from Model instead.
// This file is kept for backward compatibility during migration.
package tui

import "github.com/charmbracelet/lipgloss"

// DEPRECATED: These color variables are deprecated.
// Use config.Styles from the Model instead.
var (
	Base     = lipgloss.Color("#1e1e2e")
	Text     = lipgloss.Color("#cdd6f4")
	Mauve    = lipgloss.Color("#cba6f7")
	Red      = lipgloss.Color("#f38ba8")
	Green    = lipgloss.Color("#a6e3a1")
	Surface0 = lipgloss.Color("#313244")
	Surface1 = lipgloss.Color("#45475a")
	Overlay0 = lipgloss.Color("#6c7086")
	Yellow   = lipgloss.Color("#f9e2af")
)
