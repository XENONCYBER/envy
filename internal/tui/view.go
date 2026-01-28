// Package tui contains the files for rendering the view of the TUI application.
package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.cols == 0 {
		m.cols = 1
	}

	switch m.currentView {
	case ViewCreate:
		return m.viewCreate()
	case ViewEdit:
		return m.viewDetail()
	case ViewDetail:
		return m.viewDetail()
	case ViewEditProject:
		return m.viewEditProject()
	case ViewConfirm:
		return m.viewConfirm()
	default:
		return m.viewGrid()
	}
}

func (m Model) viewGrid() string {
	logoStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.styles.Accent).
		MarginTop(1).
		MarginBottom(1)
	logo := logoStyle.Render("███████╗███╗   ██╗██╗   ██╗██╗   ██╗\n██╔════╝████╗  ██║██║   ██║╚██╗ ██╔╝\n█████╗  ██╔██╗ ██║██║   ██║ ╚████╔╝ \n██╔══╝  ██║╚██╗██║╚██╗ ██╔╝  ╚██╔╝  \n███████╗██║ ╚████║ ╚████╔╝    ██║   \n╚══════╝╚═╝  ╚═══╝  ╚═══╝     ╚═╝   ")

	modeStyle := lipgloss.NewStyle().
		Foreground(m.styles.Base).
		Background(m.styles.Overlay0).
		Padding(0, 1).
		Bold(true)

	if m.state == StateInsert {
		modeStyle = modeStyle.Background(m.styles.Accent)
	}

	modeBadge := modeStyle.Render(m.searchMode.String())

	searchContent := lipgloss.JoinHorizontal(lipgloss.Center, modeBadge, " ", m.searchInput.View())

	searchBoxStyle := m.styles.SearchStyle
	if m.state == StateInsert {
		searchBoxStyle = m.styles.ActiveSearchStyle
	}

	searchBar := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, searchBoxStyle.Render(searchContent))

	cols := m.styles.GridCols
	visibleRows := m.styles.GridVisibleRows

	totalRows := (len(m.filtered) + cols - 1) / cols

	startRow := m.scrollOffset
	endRow := m.scrollOffset + visibleRows
	if endRow > totalRows {
		endRow = totalRows
	}

	var rows []string

	for row := startRow; row < endRow; row++ {
		var currentRow []string

		for col := 0; col < cols; col++ {
			i := row*cols + col
			if i >= len(m.filtered) {
				break
			}

			p := m.GetFilteredProject(i)
			if p == nil {
				continue
			}
			cardStyle := m.styles.CardStyle
			if i == m.selectedIdx {
				cardStyle = m.styles.SelectedCardStyle
			}

			badge := m.styles.RenderEnvironmentBadge(p.Environment)

			keysPrev := ""
			for j, k := range p.Keys {
				if j > 1 {
					keysPrev += m.styles.DimStyle.Render(fmt.Sprintf("  + %d more", len(p.Keys)-2))
					break
				}
				keysPrev += m.styles.DimStyle.Render("  • "+k.Key) + "\n"
			}

			content := lipgloss.JoinVertical(lipgloss.Left,
				lipgloss.JoinHorizontal(lipgloss.Top, m.styles.TitleStyle.Render(p.Name), "  ", badge),
				lipgloss.NewStyle().Foreground(m.styles.Surface1).Render(strings.Repeat("─", 32)),
				keysPrev,
			)

			currentRow = append(currentRow, cardStyle.Render(content))
		}

		if len(currentRow) > 0 {
			rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, currentRow...))
		}
	}

	grid := lipgloss.JoinVertical(lipgloss.Center, rows...)

	mainContent := lipgloss.JoinVertical(
		lipgloss.Center,
		logo,
		searchBar,
		"",
		grid,
	)

	bindings := GridViewBindings(m.keys, m.state)
	bottomBar := NewBottomBar(m.width, m.state, m.keys, bindings, m.styles)

	contentHeight := m.height - 3 // Reserve 3 lines for bottom bar
	contentArea := lipgloss.NewStyle().
		Height(contentHeight).
		Render(mainContent)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		contentArea,
		bottomBar.Render(),
	)
}

func (m Model) viewDetail() string {
	totalWidth := m.width - 4
	totalHeight := m.height - 6 // Reserve space for bottom bar
	if totalWidth < 60 {
		totalWidth = 60
	}
	if totalHeight < 20 {
		totalHeight = 20
	}

	var combined string

	if m.editSidebarOpen {
		mainWidth := int(float64(totalWidth) * 0.6)
		sidebarWidth := totalWidth - mainWidth - 2

		mainContent := m.renderDetailMain(mainWidth, totalHeight)
		sidebarContent := m.renderEditSidebar(sidebarWidth, totalHeight)

		combined = lipgloss.JoinHorizontal(
			lipgloss.Top,
			mainContent,
			sidebarContent,
		)

		bindings := EditViewBindings(m.keys)
		bottomBar := NewBottomBar(m.width, StateNormal, m.keys, bindings, m.styles)

		contentHeight := m.height - 3
		contentArea := lipgloss.NewStyle().Height(contentHeight).Render(combined)

		return lipgloss.JoinVertical(
			lipgloss.Left,
			contentArea,
			bottomBar.Render(),
		)
	}

	if m.historySidebarOpen {
		mainWidth := int(float64(totalWidth) * 0.6)
		sidebarWidth := totalWidth - mainWidth - 2

		mainContent := m.renderDetailMain(mainWidth, totalHeight)
		sidebarContent := m.renderHistorySidebar(sidebarWidth, totalHeight)

		combined = lipgloss.JoinHorizontal(
			lipgloss.Top,
			mainContent,
			sidebarContent,
		)

		bindings := HistoryViewBindings(m.keys)
		bottomBar := NewBottomBar(m.width, StateNormal, m.keys, bindings, m.styles)

		contentHeight := m.height - 3
		contentArea := lipgloss.NewStyle().Height(contentHeight).Render(combined)

		return lipgloss.JoinVertical(
			lipgloss.Left,
			contentArea,
			bottomBar.Render(),
		)
	}

	mainContent := m.renderDetailMain(totalWidth, totalHeight)

	bindings := DetailViewBindings(m.keys)
	bottomBar := NewBottomBar(m.width, StateNormal, m.keys, bindings, m.styles)

	contentHeight := m.height - 3
	contentArea := lipgloss.NewStyle().Height(contentHeight).Render(mainContent)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		contentArea,
		bottomBar.Render(),
	)
}

func (m Model) renderDetailMain(width, height int) string {
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.styles.Base).
		Background(m.styles.Accent).
		Padding(0, 2).
		Render(" " + m.activeProject.Name + " ")

	badge := m.styles.RenderEnvironmentBadge(m.activeProject.Environment)

	status := ""
	if m.statusMsg != "" {
		status = "  " + lipgloss.NewStyle().Foreground(m.styles.Warning).Bold(true).Render(m.statusMsg)
	}

	var keyRows []string

	for i, apiKey := range m.activeProject.Keys {
		cursor := "   "
		style := m.styles.DimStyle
		value := " ••••••••••••••••••••"

		if i == m.detailCursor {
			cursor = " ››"
			style = m.styles.TitleStyle
			if m.revealedKey {
				value = " " + apiKey.Current.Value
			}
		}

		keyName := apiKey.Key
		if len(keyName) > 30 {
			keyName = keyName[:27] + "..."
		}

		row := fmt.Sprintf("%s %-32s  %s", cursor, keyName, value)
		keyRows = append(keyRows, style.Render(row))
	}

	keysContent := lipgloss.JoinVertical(lipgloss.Left, keyRows...)

	content := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Center, header, " ", badge, status),
		"",
		keysContent,
	)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.styles.Accent).
		Width(width).
		Height(height).
		Padding(2).
		Render(content)
}

func (m Model) renderEditSidebar(width, height int) string {
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.styles.Base).
		Background(m.styles.Warning).
		Padding(0, 2).
		Render(" EDIT VALUE ")

	keyName := ""
	if len(m.activeProject.Keys) > 0 && m.detailCursor < len(m.activeProject.Keys) {
		keyName = m.activeProject.Keys[m.detailCursor].Key
	}

	subtext := m.styles.DimStyle.Render(" Editing: " + keyName)

	inputBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.styles.Warning).
		Padding(1, 2).
		Width(width - 8).
		Render(m.editInput.View())

	content := lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		subtext,
		"",
		inputBox,
	)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.styles.Warning).
		Width(width).
		Height(height).
		Padding(2).
		Render(content)
}

func (m Model) renderHistorySidebar(width, height int) string {
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.styles.Base).
		Background(m.styles.Accent).
		Padding(0, 2).
		Render(" HISTORY ")

	keyName := ""
	var historyEntries []string

	if len(m.activeProject.Keys) > 0 && m.historyKeyIdx < len(m.activeProject.Keys) {
		apiKey := m.activeProject.Keys[m.historyKeyIdx]
		keyName = apiKey.Key

		currentStyle := lipgloss.NewStyle().
			Foreground(m.styles.Success).
			Bold(true)
		currentTime := apiKey.Current.CreatedAt.Format("2006-01-02 15:04")
		currentValue := apiKey.Current.Value
		if len(currentValue) > width-16 {
			currentValue = currentValue[:width-19] + "..."
		}

		historyEntries = append(historyEntries,
			currentStyle.Render("Current:"),
			m.styles.DimStyle.Render("  "+currentTime),
			lipgloss.NewStyle().Foreground(m.styles.Text).Render("  "+currentValue),
			"",
		)

		if len(apiKey.History) > 0 {
			historyEntries = append(historyEntries,
				lipgloss.NewStyle().Foreground(m.styles.Overlay0).Bold(true).Render("Previous:"),
			)

			for i, entry := range apiKey.History {
				if i >= 5 {
					remaining := len(apiKey.History) - 5
					historyEntries = append(historyEntries,
						m.styles.DimStyle.Render(fmt.Sprintf("  + %d more entries", remaining)),
					)
					break
				}
				entryTime := entry.CreatedAt.Format("2006-01-02 15:04")
				entryValue := entry.Value
				if len(entryValue) > width-16 {
					entryValue = entryValue[:width-19] + "..."
				}
				historyEntries = append(historyEntries,
					m.styles.DimStyle.Render("  "+entryTime),
					lipgloss.NewStyle().Foreground(m.styles.Overlay0).Render("  "+entryValue),
					"",
				)
			}
		} else {
			historyEntries = append(historyEntries,
				m.styles.DimStyle.Italic(true).Render("  No previous history"),
			)
		}
	}

	subtext := m.styles.DimStyle.Render(" Key: " + keyName)

	historyContent := lipgloss.JoinVertical(lipgloss.Left, historyEntries...)

	content := lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		subtext,
		"",
		historyContent,
	)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.styles.Accent).
		Width(width).
		Height(height).
		Padding(2).
		Render(content)
}

func (m Model) viewConfirm() string {
	var bgView string
	switch m.previousView {
	case ViewGrid:
		bgView = m.viewGrid()
	case ViewDetail:
		bgView = m.viewDetail()
	default:
		bgView = m.viewGrid()
	}

	dialogWidth := 50
	dialogHeight := 7

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.styles.Base).
		Background(m.styles.Error).
		Padding(0, 2).
		Render(" CONFIRM ")

	message := lipgloss.NewStyle().
		Foreground(m.styles.Text).
		Width(dialogWidth - 6).
		Align(lipgloss.Center).
		Render(m.confirmMessage)

	options := lipgloss.NewStyle().
		Foreground(m.styles.Overlay0).
		Width(dialogWidth - 6).
		Align(lipgloss.Center).
		Render("Press (y) to confirm or (n) to cancel")

	dialogContent := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		message,
		"",
		options,
	)

	dialog := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.styles.Error).
		Padding(1, 2).
		Width(dialogWidth).
		Height(dialogHeight).
		Render(dialogContent)

	centeredDialog := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		dialog,
	)

	// MAYBE: For now, just show the dialog (could overlay on bg in future)
	_ = bgView
	return centeredDialog
}

func (m Model) viewEditProject() string {
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
			Render("EDIT PROJECT"),
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
	if m.editProjectFocus == 0 {
		if m.state == StateInsert {
			projectNameLabel = lipgloss.NewStyle().Foreground(focusColor).Bold(true).Render("› " + projectNameLabel + " (editing)")
		} else {
			projectNameLabel = lipgloss.NewStyle().Foreground(m.styles.Warning).Bold(true).Render("› " + projectNameLabel)
		}
	} else {
		projectNameLabel = lipgloss.NewStyle().Foreground(dimColor).Render("  " + projectNameLabel)
	}
	projectNameField := getFieldStyle(m.editProjectFocus == 0).Render(m.editProjectName.View())

	keysLabel := "Keys"
	if m.editProjectFocus == 1 {
		keysLabel = lipgloss.NewStyle().Foreground(m.styles.Warning).Bold(true).Render("› " + keysLabel + " (select to delete)")
	} else {
		keysLabel = lipgloss.NewStyle().Foreground(dimColor).Render("  " + keysLabel)
	}

	var keyRows []string
	for i, k := range m.activeProject.Keys {
		cursor := "  "
		style := lipgloss.NewStyle().Foreground(dimColor)
		if m.editProjectFocus == 1 && i == m.editProjectKeyIdx {
			cursor = "› "
			style = lipgloss.NewStyle().Foreground(textColor).Bold(true)
		}
		keyRows = append(keyRows, style.Render(cursor+k.Key))
	}
	if len(keyRows) == 0 {
		keyRows = append(keyRows, lipgloss.NewStyle().Foreground(dimColor).Italic(true).Render("  (no keys)"))
	}

	keysListStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(normalColor).
		Padding(0, 1).
		Width(inputWidth).
		MaxHeight(6)
	if m.editProjectFocus == 1 {
		keysListStyle = keysListStyle.BorderForeground(m.styles.Warning)
	}
	keysList := keysListStyle.Render(strings.Join(keyRows, "\n"))

	newKeyNameLabel := "New Key Name"
	if m.editProjectFocus == 2 {
		if m.state == StateInsert {
			newKeyNameLabel = lipgloss.NewStyle().Foreground(focusColor).Bold(true).Render("› " + newKeyNameLabel + " (editing)")
		} else {
			newKeyNameLabel = lipgloss.NewStyle().Foreground(m.styles.Warning).Bold(true).Render("› " + newKeyNameLabel)
		}
	} else {
		newKeyNameLabel = lipgloss.NewStyle().Foreground(dimColor).Render("  " + newKeyNameLabel)
	}
	newKeyNameField := getFieldStyle(m.editProjectFocus == 2).Render(m.editProjectNewKey[0].View())

	newKeyValueLabel := "New Key Value"
	if m.editProjectFocus == 3 {
		if m.state == StateInsert {
			newKeyValueLabel = lipgloss.NewStyle().Foreground(focusColor).Bold(true).Render("› " + newKeyValueLabel + " (editing)")
		} else {
			newKeyValueLabel = lipgloss.NewStyle().Foreground(m.styles.Warning).Bold(true).Render("› " + newKeyValueLabel)
		}
	} else {
		newKeyValueLabel = lipgloss.NewStyle().Foreground(dimColor).Render("  " + newKeyValueLabel)
	}
	newKeyValueField := getFieldStyle(m.editProjectFocus == 3).Render(m.editProjectNewKey[1].View())

	addKeyStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 2).
		Foreground(textColor).
		BorderForeground(normalColor)

	saveStyle := addKeyStyle.Copy()

	if m.editProjectFocus == 4 {
		addKeyStyle = addKeyStyle.
			BorderForeground(m.styles.Success).
			Foreground(m.styles.Success).
			Bold(true)
	}

	if m.editProjectFocus == 5 {
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

	statusSection := ""
	if m.statusMsg != "" {
		statusSection = lipgloss.NewStyle().
			Foreground(m.styles.Warning).
			Render(" " + m.statusMsg)
	}

	var formParts []string
	formParts = append(formParts, title)
	formParts = append(formParts, "")
	formParts = append(formParts, projectNameLabel)
	formParts = append(formParts, projectNameField)
	formParts = append(formParts, "")
	formParts = append(formParts, keysLabel)
	formParts = append(formParts, keysList)
	formParts = append(formParts, "")
	formParts = append(formParts, newKeyNameLabel)
	formParts = append(formParts, newKeyNameField)
	formParts = append(formParts, newKeyValueLabel)
	formParts = append(formParts, newKeyValueField)
	formParts = append(formParts, "")
	formParts = append(formParts, buttons)

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

	bindings := EditProjectViewBindings(m.keys, m.state)
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
