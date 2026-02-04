// Package config handles application configuration including themes for future Lua integration
package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
)

// Theme defines all colors and styling used in the application
// This will be configurable via Lua in the future
type Theme struct {
	// Base colors
	Base     string `json:"base"`
	Text     string `json:"text"`
	Accent   string `json:"accent"`
	Surface0 string `json:"surface0"`
	Surface1 string `json:"surface1"`
	Overlay0 string `json:"overlay0"`

	// Semantic colors
	Success string `json:"success"`
	Warning string `json:"warning"`
	Error   string `json:"error"`

	// Environment badge colors
	ProdBg  string `json:"prod_bg"`
	DevBg   string `json:"dev_bg"`
	StageBg string `json:"stage_bg"`

	// History section colors
	CurrentBg  string `json:"current_bg"`
	PreviousBg string `json:"previous_bg"`

	// Grid layout
	GridCols        int `json:"grid_cols"`
	GridVisibleRows int `json:"grid_visible_rows"`

	// Card dimensions
	CardWidth  int `json:"card_width"`
	CardHeight int `json:"card_height"`
}

// Default theme: Catppuccin Mocha
func DefaultTheme() Theme {
	return Theme{
		Base:     "#1e1e2e",
		Text:     "#cdd6f4",
		Accent:   "#cba6f7",
		Surface0: "#313244",
		Surface1: "#45475a",
		Overlay0: "#6c7086",

		// Semantic colors
		Success: "#a6e3a1",
		Warning: "#f9e2af",
		Error:   "#f38ba8",

		// Environment badges
		ProdBg:  "#f38ba8", // Red
		DevBg:   "#a6e3a1", // Green
		StageBg: "#f9e2af", // Yellow

		// History section backgrounds
		CurrentBg:  "#a6e3a1", // Green
		PreviousBg: "#f9e2af", // Yellow

		// Grid layout
		GridCols:        3,
		GridVisibleRows: 2,

		// Card dimensions
		CardWidth:  38,
		CardHeight: 9,
	}
}

func LoadTheme() Theme {
	configPath := getThemePath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		return DefaultTheme()
	}

	var theme Theme
	if err := json.Unmarshal(data, &theme); err != nil {
		return DefaultTheme()
	}

	theme = mergeWithDefaults(theme)
	return theme
}

func SaveTheme(theme Theme) error {
	configPath := getThemePath()

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(theme, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0o644)
}

func getThemePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "theme.json"
	}
	return filepath.Join(home, ".config", "envy", "theme.json")
}

func mergeWithDefaults(theme Theme) Theme {
	defaults := DefaultTheme()

	if theme.Base == "" {
		theme.Base = defaults.Base
	}
	if theme.Text == "" {
		theme.Text = defaults.Text
	}
	if theme.Accent == "" {
		theme.Accent = defaults.Accent
	}
	if theme.Surface0 == "" {
		theme.Surface0 = defaults.Surface0
	}
	if theme.Surface1 == "" {
		theme.Surface1 = defaults.Surface1
	}
	if theme.Overlay0 == "" {
		theme.Overlay0 = defaults.Overlay0
	}
	if theme.Success == "" {
		theme.Success = defaults.Success
	}
	if theme.Warning == "" {
		theme.Warning = defaults.Warning
	}
	if theme.Error == "" {
		theme.Error = defaults.Error
	}
	if theme.ProdBg == "" {
		theme.ProdBg = defaults.ProdBg
	}
	if theme.DevBg == "" {
		theme.DevBg = defaults.DevBg
	}
	if theme.StageBg == "" {
		theme.StageBg = defaults.StageBg
	}
	if theme.CurrentBg == "" {
		theme.CurrentBg = defaults.CurrentBg
	}
	if theme.PreviousBg == "" {
		theme.PreviousBg = defaults.PreviousBg
	}
	if theme.GridCols == 0 {
		theme.GridCols = defaults.GridCols
	}
	if theme.GridVisibleRows == 0 {
		theme.GridVisibleRows = defaults.GridVisibleRows
	}
	if theme.CardWidth == 0 {
		theme.CardWidth = defaults.CardWidth
	}
	if theme.CardHeight == 0 {
		theme.CardHeight = defaults.CardHeight
	}

	return theme
}

type Styles struct {
	Base     lipgloss.Color
	Text     lipgloss.Color
	Accent   lipgloss.Color
	Surface0 lipgloss.Color
	Surface1 lipgloss.Color
	Overlay0 lipgloss.Color
	Success  lipgloss.Color
	Warning  lipgloss.Color
	Error    lipgloss.Color

	// History section colors
	CurrentBg  lipgloss.Color
	PreviousBg lipgloss.Color

	// Component styles
	SearchStyle       lipgloss.Style
	ActiveSearchStyle lipgloss.Style
	CardStyle         lipgloss.Style
	SelectedCardStyle lipgloss.Style
	TitleStyle        lipgloss.Style
	DimStyle          lipgloss.Style
	ProdBadge         lipgloss.Style
	DevBadge          lipgloss.Style
	StageBadge        lipgloss.Style

	// Grid layout
	GridCols        int
	GridVisibleRows int
	CardWidth       int
	CardHeight      int
	FullCardWidth   int
}

func NewStyles(theme Theme) Styles {
	s := Styles{
		Base:       lipgloss.Color(theme.Base),
		Text:       lipgloss.Color(theme.Text),
		Accent:     lipgloss.Color(theme.Accent),
		Surface0:   lipgloss.Color(theme.Surface0),
		Surface1:   lipgloss.Color(theme.Surface1),
		Overlay0:   lipgloss.Color(theme.Overlay0),
		Success:    lipgloss.Color(theme.Success),
		Warning:    lipgloss.Color(theme.Warning),
		Error:      lipgloss.Color(theme.Error),
		CurrentBg:  lipgloss.Color(theme.CurrentBg),
		PreviousBg: lipgloss.Color(theme.PreviousBg),

		GridCols:        theme.GridCols,
		GridVisibleRows: theme.GridVisibleRows,
		CardWidth:       theme.CardWidth,
		CardHeight:      theme.CardHeight,
		FullCardWidth:   theme.CardWidth + 6,
	}

	s.SearchStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.Surface1).
		Padding(0, 2).
		Width(60)

	s.ActiveSearchStyle = s.SearchStyle.Copy().
		BorderForeground(s.Accent).
		Border(lipgloss.ThickBorder())

	s.CardStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.Surface1).
		Padding(1, 2).
		Margin(0, 1).
		Width(theme.CardWidth).
		Height(theme.CardHeight)

	s.SelectedCardStyle = s.CardStyle.Copy().
		BorderForeground(s.Accent).
		Border(lipgloss.ThickBorder())

	s.TitleStyle = lipgloss.NewStyle().Bold(true).Foreground(s.Accent).MarginBottom(1)
	s.DimStyle = lipgloss.NewStyle().Foreground(s.Overlay0)

	s.ProdBadge = lipgloss.NewStyle().
		Background(lipgloss.Color(theme.ProdBg)).
		Foreground(s.Base).
		Padding(0, 2).Bold(true).MarginLeft(2)

	s.DevBadge = lipgloss.NewStyle().
		Background(lipgloss.Color(theme.DevBg)).
		Foreground(s.Base).
		Padding(0, 2).Bold(true).MarginLeft(2)

	s.StageBadge = lipgloss.NewStyle().
		Background(lipgloss.Color(theme.StageBg)).
		Foreground(s.Base).
		Padding(0, 2).Bold(true).MarginLeft(2)

	return s
}

func (s Styles) RenderEnvironmentBadge(env string) string {
	switch env {
	case "prod":
		return s.ProdBadge.Render(" PROD ")
	case "stage":
		return s.StageBadge.Render(" STAGE ")
	default:
		return s.DevBadge.Render(" DEV ")
	}
}
