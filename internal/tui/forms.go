package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) viewCreate() string {
	containerWidth := 70
	if m.width < 80 {
		containerWidth = m.width - 10
	}
	if containerWidth > 80 {
		containerWidth = 80
	}

	inputWidth := containerWidth - 8

	focusColor := m.styles.Accent
	normalColor := m.styles.Surface1
	textColor := m.styles.Text
	dimColor := m.styles.Overlay0

	modeText := "NORMAL"
	modeColor := m.styles.Success
	if m.state == StateInsert {
		modeText = "INSERT"
		modeColor = m.styles.Warning
	}

	modeIndicator := lipgloss.NewStyle().
		Foreground(m.styles.Base).
		Background(modeColor).
		Bold(true).
		Padding(0, 1).
		Render(modeText)

	titleBar := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().
			Bold(true).
			Foreground(m.styles.Base).
			Background(focusColor).
			Padding(0, 2).
			Render("CREATE PROJECT"),
		" ",
		modeIndicator,
	)

	title := lipgloss.NewStyle().
		Width(containerWidth).
		Align(lipgloss.Center).
		Render(titleBar)

	getFieldStyle := func(focused bool) lipgloss.Style {
		borderColor := normalColor
		if focused && m.state == StateInsert {
			borderColor = focusColor
		} else if focused {
			borderColor = m.styles.Warning
		}
		return lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(0, 1).
			Width(inputWidth)
	}

	projectNameLabel := "Project Name"
	if m.focusIndex == 0 {
		if m.state == StateInsert {
			projectNameLabel = lipgloss.NewStyle().Foreground(focusColor).Bold(true).Render("› " + projectNameLabel + " (editing)")
		} else {
			projectNameLabel = lipgloss.NewStyle().Foreground(m.styles.Warning).Bold(true).Render("› " + projectNameLabel)
		}
	} else {
		projectNameLabel = lipgloss.NewStyle().Foreground(dimColor).Render("  " + projectNameLabel)
	}
	projectNameField := getFieldStyle(m.focusIndex == 0).Render(m.inputs[0].View())

	envLabel := "Environment"
	if m.focusIndex == 1 {
		envLabel = lipgloss.NewStyle().Foreground(m.styles.Warning).Bold(true).Render("› " + envLabel)
	} else {
		envLabel = lipgloss.NewStyle().Foreground(dimColor).Render("  " + envLabel)
	}

	devStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Foreground(textColor).
		BorderForeground(normalColor)

	prodStyle := devStyle.Copy()
	stageStyle := devStyle.Copy()

	switch m.selectedEnv {
	case EnvOptionDev:
		devStyle = devStyle.BorderForeground(m.styles.Success).Background(m.styles.Success).Foreground(m.styles.Base).Bold(true)
	case EnvOptionProd:
		prodStyle = prodStyle.BorderForeground(m.styles.Error).Background(m.styles.Error).Foreground(m.styles.Base).Bold(true)
	case EnvOptionStage:
		stageStyle = stageStyle.BorderForeground(m.styles.Warning).Background(m.styles.Warning).Foreground(m.styles.Base).Bold(true)
	}

	envBoxes := lipgloss.JoinHorizontal(
		lipgloss.Top,
		devStyle.Render(" DEV "),
		" ",
		prodStyle.Render(" PROD "),
		" ",
		stageStyle.Render(" STAGE "),
	)

	keyNameLabel := "Key Name"
	if m.focusIndex == 2 {
		if m.state == StateInsert {
			keyNameLabel = lipgloss.NewStyle().Foreground(focusColor).Bold(true).Render("› " + keyNameLabel + " (editing)")
		} else {
			keyNameLabel = lipgloss.NewStyle().Foreground(m.styles.Warning).Bold(true).Render("› " + keyNameLabel)
		}
	} else {
		keyNameLabel = lipgloss.NewStyle().Foreground(dimColor).Render("  " + keyNameLabel)
	}
	keyNameField := getFieldStyle(m.focusIndex == 2).Render(m.inputs[1].View())

	keyValueLabel := "Key Value"
	if m.focusIndex == 3 {
		if m.state == StateInsert {
			keyValueLabel = lipgloss.NewStyle().Foreground(focusColor).Bold(true).Render("› " + keyValueLabel + " (editing)")
		} else {
			keyValueLabel = lipgloss.NewStyle().Foreground(m.styles.Warning).Bold(true).Render("› " + keyValueLabel)
		}
	} else {
		keyValueLabel = lipgloss.NewStyle().Foreground(dimColor).Render("  " + keyValueLabel)
	}
	keyValueField := getFieldStyle(m.focusIndex == 3).Render(m.inputs[2].View())

	addKeyStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 2).
		Foreground(textColor).
		BorderForeground(normalColor)

	saveStyle := addKeyStyle.Copy()

	if m.focusIndex == 4 {
		addKeyStyle = addKeyStyle.
			BorderForeground(m.styles.Success).
			Foreground(m.styles.Success).
			Bold(true)
	}

	if m.focusIndex == 5 {
		saveStyle = saveStyle.
			BorderForeground(focusColor).
			Foreground(focusColor).
			Bold(true)
	}

	buttons := lipgloss.JoinHorizontal(
		lipgloss.Top,
		addKeyStyle.Render("+ Add"),
		"  ",
		saveStyle.Render("Save"),
	)

	var pendingSection string
	if len(m.pendingKeys) > 0 {
		pendingTitle := lipgloss.NewStyle().
			Foreground(m.styles.Success).
			Render(fmt.Sprintf("Keys: %d", len(m.pendingKeys)))

		var keysList []string
		for i, k := range m.pendingKeys {
			if i >= 3 {
				keysList = append(keysList, lipgloss.NewStyle().Foreground(dimColor).Render(fmt.Sprintf("  +%d more", len(m.pendingKeys)-3)))
				break
			}
			keysList = append(keysList, lipgloss.NewStyle().Foreground(dimColor).Render("  • "+k.Key))
		}

		pendingSection = pendingTitle + "\n" + strings.Join(keysList, "\n")
	}

	statusSection := ""
	if m.statusMsg != "" {
		statusSection = lipgloss.NewStyle().
			Foreground(m.styles.Error).
			Render(" ! " + m.statusMsg)
	}

	var formParts []string

	formParts = append(formParts, title)
	formParts = append(formParts, "")
	formParts = append(formParts, projectNameLabel)
	formParts = append(formParts, projectNameField)
	formParts = append(formParts, "")
	formParts = append(formParts, envLabel)
	formParts = append(formParts, envBoxes)
	formParts = append(formParts, "")
	formParts = append(formParts, keyNameLabel)
	formParts = append(formParts, keyNameField)
	formParts = append(formParts, keyValueLabel)
	formParts = append(formParts, keyValueField)
	formParts = append(formParts, "")
	formParts = append(formParts, buttons)

	if pendingSection != "" {
		formParts = append(formParts, "")
		formParts = append(formParts, pendingSection)
	}

	if statusSection != "" {
		formParts = append(formParts, "")
		formParts = append(formParts, statusSection)
	}

	form := lipgloss.JoinVertical(lipgloss.Left, formParts...)

	container := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(focusColor).
		Padding(1, 2).
		Width(containerWidth).
		Render(form)

	bindings := CreateViewBindings(m.keys, m.state)
	bottomBar := NewBottomBar(m.width, m.state, m.keys, bindings, m.styles)

	contentHeight := m.height - 3
	centeredContainer := lipgloss.Place(
		m.width,
		contentHeight,
		lipgloss.Center,
		lipgloss.Center,
		container,
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		centeredContainer,
		bottomBar.Render(),
	)
}

// WARNING: This function is kept for compatibility but shouldn't be called directly
func (m Model) viewEdit() string {
	// This should not be called anymore - edit is now a sidebar
	return m.viewDetail()
}
