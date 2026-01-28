package tui

import (
	"time"

	"envy/internal/domain"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type statusClearMsg struct{}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Hard safety exit
		if msg.String() == m.keys.ForceQuit {
			return m, tea.Quit
		}

		switch m.currentView {
		case ViewGrid:
			return m.updateGrid(msg)
		case ViewDetail:
			return m.updateDetail(msg)
		case ViewCreate:
			return m.updateCreate(msg)
		case ViewEdit:
			return m.updateEdit(msg)
		case ViewEditProject:
			return m.updateEditProject(msg)
		case ViewConfirm:
			return m.updateConfirm(msg)
		}

	case statusClearMsg:
		m.statusMsg = ""
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.cols = (m.width - 4) / m.styles.FullCardWidth
		if m.cols < 1 {
			m.cols = 1
		}
	}
	return m, nil
}

func (m Model) updateGrid(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	k := msg.String()

	if m.state == StateInsert {
		if m.keys.MatchesKey(k, m.keys.Enter, m.keys.Back) {
			m.state = StateNormal
			m.searchInput.Blur()
			return m, nil
		}
		if k == m.keys.Tab {
			m.searchMode = m.searchMode.Next()
			m.filtered = filterProjects(m.projects, m.searchInput.Value(), m.searchMode)
			m.selectedIdx = 0
			m.scrollOffset = 0
			return m, nil
		}
		m.searchInput, cmd = m.searchInput.Update(msg)
		m.filtered = filterProjects(m.projects, m.searchInput.Value(), m.searchMode)

		m.selectedIdx = 0
		m.scrollOffset = 0
		return m, cmd
	}

	switch k {
	case m.keys.Create:
		m.currentView = ViewCreate
		m.focusIndex = 0
		m.state = StateNormal
		m.selectedEnv = EnvOptionDev
		m.pendingKeys = []domain.APIKey{}
		for i := range m.inputs {
			m.inputs[i].SetValue("")
			m.inputs[i].Blur()
		}
		return m, nil

	case m.keys.Enter:
		if len(m.filtered) > 0 {
			m.currentView = ViewDetail
			if p := m.GetFilteredProject(m.selectedIdx); p != nil {
				m.activeProject = *p
			}
			m.detailCursor = 0
			m.revealedKey = false
			m.editSidebarOpen = false
			m.historySidebarOpen = false
		}

	case m.keys.Quit:
		return m, tea.Quit

	case m.keys.Search, "/":
		m.state = StateInsert
		m.searchInput.Focus()
		m.scrollOffset = 0
		return m, textinput.Blink

	case m.keys.Delete:
		if len(m.filtered) > 0 && m.selectedIdx < len(m.filtered) {
			if p := m.GetFilteredProject(m.selectedIdx); p != nil {
				m.confirmAction = ConfirmDeleteProject
				m.confirmMessage = "Delete project '" + p.Name + "' (" + p.Environment + ")?"
				m.previousView = ViewGrid
				m.currentView = ViewConfirm
			}
		}
		return m, nil
	}

	cols := m.styles.GridCols
	if m.keys.IsNavigationLeft(k) {
		if m.selectedIdx%cols > 0 {
			m.selectedIdx--
			m.adjustScroll()
		}
	} else if m.keys.IsNavigationRight(k) {
		if m.selectedIdx%cols < cols-1 && m.selectedIdx < len(m.filtered)-1 {
			m.selectedIdx++
			m.adjustScroll()
		}
	} else if m.keys.IsNavigationDown(k) {
		if m.selectedIdx+cols < len(m.filtered) {
			m.selectedIdx += cols
			m.adjustScroll()
		}
	} else if m.keys.IsNavigationUp(k) {
		if m.selectedIdx-cols >= 0 {
			m.selectedIdx -= cols
			m.adjustScroll()
		}
	}

	return m, nil
}

func (m Model) updateDetail(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	k := msg.String()

	if m.editSidebarOpen {
		return m.updateEdit(msg)
	}

	if m.historySidebarOpen {
		return m.updateHistorySidebar(msg)
	}

	switch k {
	case m.keys.Back:
		m.currentView = ViewGrid
		m.activeProject = domain.Project{}
		return m, nil

	case m.keys.History:
		if len(m.activeProject.Keys) == 0 {
			m.statusMsg = "No keys to show history"
			return m, nil
		}
		m.historySidebarOpen = true
		m.historyKeyIdx = m.detailCursor
		return m, nil

	case m.keys.Edit:
		if len(m.activeProject.Keys) == 0 {
			m.statusMsg = "No keys to edit"
			return m, nil
		}
		if m.detailCursor >= len(m.activeProject.Keys) {
			m.detailCursor = len(m.activeProject.Keys) - 1
		}
		m.editSidebarOpen = true
		m.editInput.SetValue(m.activeProject.Keys[m.detailCursor].Current.Value)
		m.editInput.Focus()
		return m, textinput.Blink

	case m.keys.EditProject:
		m.currentView = ViewEditProject
		m.editProjectName.SetValue(m.activeProject.Name)
		m.editProjectKeyIdx = 0
		m.editProjectFocus = 0
		m.editProjectNewKey[0].SetValue("")
		m.editProjectNewKey[1].SetValue("")
		m.state = StateNormal
		return m, nil

	case m.keys.Yank:
		if len(m.activeProject.Keys) == 0 {
			m.statusMsg = "No keys to copy"
			return m, nil
		}
		if m.detailCursor >= len(m.activeProject.Keys) {
			m.detailCursor = len(m.activeProject.Keys) - 1
		}
		val := m.activeProject.Keys[m.detailCursor].Current.Value
		if err := clipboard.WriteAll(val); err != nil {
			m.statusMsg = "Failed to copy"
		} else {
			m.statusMsg = "Copied (will clear in 30s)"
			go func() {
				time.Sleep(30 * time.Second)
				clipboard.WriteAll("")
			}()
		}
		return m, tea.Tick(time.Second, func(time.Time) tea.Msg {
			return statusClearMsg{}
		})

	case m.keys.Space, m.keys.Enter:
		m.revealedKey = !m.revealedKey

	case m.keys.Delete:
		if len(m.activeProject.Keys) > 0 && m.detailCursor < len(m.activeProject.Keys) {
			keyName := m.activeProject.Keys[m.detailCursor].Key
			m.confirmAction = ConfirmDeleteKey
			m.confirmMessage = "Delete key '" + keyName + "'?"
			m.previousView = ViewDetail
			m.currentView = ViewConfirm
		}
		return m, nil
	}

	if m.keys.IsNavigationDown(k) {
		if m.detailCursor < len(m.activeProject.Keys)-1 {
			m.detailCursor++
			m.revealedKey = false
		}
	} else if m.keys.IsNavigationUp(k) {
		if m.detailCursor > 0 {
			m.detailCursor--
			m.revealedKey = false
		}
	}

	return m, nil
}

func (m Model) updateCreate(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	k := msg.String()

	m.statusMsg = ""

	if m.state == StateInsert {
		if k == m.keys.Back {
			m.state = StateNormal
			m.blurCurrentInput()
			return m, nil
		}

		if m.focusIndex == 0 || m.focusIndex == 2 || m.focusIndex == 3 {
			inputIdx := m.getInputIndex()
			if inputIdx >= 0 {
				var cmd tea.Cmd
				m.inputs[inputIdx], cmd = m.inputs[inputIdx].Update(msg)
				return m, cmd
			}
		}
		return m, nil
	}

	switch k {
	case m.keys.Quit:
		m.currentView = ViewGrid
		m.state = StateNormal
		return m, nil

	case m.keys.Save:
		return m.saveNewProject()

	case m.keys.Add:
		return m.addPendingKey()

	case m.keys.Search:
		if m.focusIndex == 0 || m.focusIndex == 2 || m.focusIndex == 3 {
			m.state = StateInsert
			return m, m.focusCurrentInput()
		}
		return m, nil
	}

	// Focus indices for create form:
	// 0: Project Name input
	// 1: Environment selector (dev/prod/stage)
	// 2: Key Name input
	// 3: Key Value input
	// 4: Add Key button
	// 5: Save Project button

	if k == m.keys.Tab {
		m.blurCurrentInput()
		m.focusIndex++
		if m.focusIndex > 5 {
			m.focusIndex = 0
		}
		return m, nil
	} else if k == m.keys.ShiftTab {
		m.blurCurrentInput()
		m.focusIndex--
		if m.focusIndex < 0 {
			m.focusIndex = 5
		}
		return m, nil
	}

	if m.keys.IsNavigationDown(k) {
		m.blurCurrentInput()
		m.focusIndex++
		if m.focusIndex > 5 {
			m.focusIndex = 0
		}
		return m, nil
	} else if m.keys.IsNavigationUp(k) {
		m.blurCurrentInput()
		m.focusIndex--
		if m.focusIndex < 0 {
			m.focusIndex = 5
		}
		return m, nil
	}

	if m.focusIndex == 1 {
		if m.keys.IsNavigationRight(k) || k == m.keys.Space {
			m.selectedEnv++
			if m.selectedEnv > EnvOptionStage {
				m.selectedEnv = EnvOptionDev
			}
			return m, nil
		} else if m.keys.IsNavigationLeft(k) {
			m.selectedEnv--
			if m.selectedEnv < EnvOptionDev {
				m.selectedEnv = EnvOptionStage
			}
			return m, nil
		}
	}

	if k == m.keys.Enter {
		switch m.focusIndex {
		case 0, 2, 3:
			m.state = StateInsert
			return m, m.focusCurrentInput()

		case 4:
			kName := m.inputs[1].Value()
			kVal := m.inputs[2].Value()

			if kName == "" {
				m.statusMsg = "Key name cannot be empty"
				return m, nil
			}

			if err := domain.ValidateKeyName(kName); err != nil {
				m.statusMsg = err.Error()
				return m, nil
			}

			m.pendingKeys = append(m.pendingKeys, domain.APIKey{
				Title: kName,
				Key:   kName,
				Current: domain.SecretVersion{
					Value:     kVal,
					CreatedAt: time.Now(),
					CreatedBy: "tui",
				},
				History: []domain.SecretVersion{},
			})

			m.inputs[1].SetValue("")
			m.inputs[2].SetValue("")
			m.blurCurrentInput()
			m.focusIndex = 2

		case 5:
			name := m.inputs[0].Value()
			env := m.selectedEnv.String()

			if name == "" {
				m.statusMsg = "Project name cannot be empty"
				return m, nil
			}

			if err := domain.ValidateProjectName(name); err != nil {
				m.statusMsg = err.Error()
				return m, nil
			}

			if len(m.pendingKeys) == 0 {
				m.statusMsg = "Add at least one key before saving"
				return m, nil
			}

			newProj := domain.Project{
				Name:        name,
				Environment: env,
				Keys:        m.pendingKeys,
			}

			if err := m.vault.CreateProject(newProj); err != nil {
				m.statusMsg = err.Error()
				return m, nil
			}

			if err := m.vault.Save(); err != nil {
				m.statusMsg = "ERROR: Failed to save"
				return m, nil
			}

			// Refresh projects from vault
			m.projects = m.vault.GetProjects()
			m.RefreshFiltered()
			m.currentView = ViewGrid
			m.state = StateNormal
		}
		return m, nil
	}

	return m, nil
}

func (m *Model) blurCurrentInput() {
	inputIdx := m.getInputIndex()
	if inputIdx >= 0 {
		m.inputs[inputIdx].Blur()
	}
}

func (m *Model) focusCurrentInput() tea.Cmd {
	inputIdx := m.getInputIndex()
	if inputIdx >= 0 {
		return m.inputs[inputIdx].Focus()
	}
	return nil
}

func (m Model) getInputIndex() int {
	switch m.focusIndex {
	case 0:
		return 0 // Project Name
	case 2:
		return 1 // Key Name
	case 3:
		return 2 // Key Value
	default:
		return -1 // Not an input field
	}
}

func (m Model) updateEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	k := msg.String()

	switch k {
	case m.keys.Back:
		m.editSidebarOpen = false
		m.editInput.Blur()
		return m, nil

	case m.keys.Enter:
		if len(m.activeProject.Keys) == 0 {
			m.statusMsg = "No keys to edit"
			m.editSidebarOpen = false
			return m, nil
		}
		if m.detailCursor >= len(m.activeProject.Keys) {
			m.detailCursor = len(m.activeProject.Keys) - 1
		}

		newValue := m.editInput.Value()
		keyName := m.activeProject.Keys[m.detailCursor].Key

		if err := m.vault.UpdateKey(m.activeProject.Name, m.activeProject.Environment, keyName, newValue); err != nil {
			m.statusMsg = "ERROR: " + err.Error()
			return m, nil
		}

		if err := m.vault.Save(); err != nil {
			m.statusMsg = "ERROR: Failed to save"
			return m, nil
		}

		m.projects = m.vault.GetProjects()
		if proj, err := m.vault.GetProject(m.activeProject.Name, m.activeProject.Environment); err == nil {
			m.activeProject = *proj
		}
		m.RefreshFiltered()

		m.statusMsg = "Updated"
		m.editSidebarOpen = false
		m.editInput.Blur()
		return m, nil
	}

	var cmd tea.Cmd
	m.editInput, cmd = m.editInput.Update(msg)
	return m, cmd
}

func (m Model) updateConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	k := msg.String()

	switch k {
	case "y", "Y":
		switch m.confirmAction {
		case ConfirmDeleteProject:
			if len(m.filtered) > 0 && m.selectedIdx < len(m.filtered) {
				p := m.GetFilteredProject(m.selectedIdx)
				if p == nil {
					m.currentView = m.previousView
					m.confirmAction = ConfirmNone
					return m, nil
				}
				if err := m.vault.DeleteProject(p.Name, p.Environment); err != nil {
					m.statusMsg = "ERROR: " + err.Error()
					m.currentView = m.previousView
					m.confirmAction = ConfirmNone
					return m, nil
				}

				if err := m.vault.Save(); err != nil {
					m.statusMsg = "ERROR: Failed to save"
					m.currentView = m.previousView
					m.confirmAction = ConfirmNone
					return m, nil
				}

				m.projects = m.vault.GetProjects()
				m.RefreshFiltered()

				if m.selectedIdx >= len(m.filtered) {
					m.selectedIdx = len(m.filtered) - 1
				}
				if m.selectedIdx < 0 {
					m.selectedIdx = 0
				}

				m.statusMsg = "Project deleted"
			}

		case ConfirmDeleteKey:
			if len(m.activeProject.Keys) > 0 && m.detailCursor < len(m.activeProject.Keys) {
				keyName := m.activeProject.Keys[m.detailCursor].Key
				if err := m.vault.DeleteKey(m.activeProject.Name, m.activeProject.Environment, keyName); err != nil {
					m.statusMsg = "ERROR: " + err.Error()
					m.currentView = m.previousView
					m.confirmAction = ConfirmNone
					return m, nil
				}

				if err := m.vault.Save(); err != nil {
					m.statusMsg = "ERROR: Failed to save"
					m.currentView = m.previousView
					m.confirmAction = ConfirmNone
					return m, nil
				}

				m.projects = m.vault.GetProjects()
				if proj, err := m.vault.GetProject(m.activeProject.Name, m.activeProject.Environment); err == nil {
					m.activeProject = *proj
				}
				m.RefreshFiltered()

				if m.detailCursor >= len(m.activeProject.Keys) {
					m.detailCursor = len(m.activeProject.Keys) - 1
				}
				if m.detailCursor < 0 {
					m.detailCursor = 0
				}

				m.statusMsg = "Key deleted"
			}
		}

		m.currentView = m.previousView
		m.confirmAction = ConfirmNone
		return m, nil

	case "n", "N", m.keys.Back:
		m.currentView = m.previousView
		m.confirmAction = ConfirmNone
		return m, nil
	}

	return m, nil
}

// addPendingKey: adds a key to the pending keys list (for create form)
func (m Model) addPendingKey() (tea.Model, tea.Cmd) {
	kName := m.inputs[1].Value()
	kVal := m.inputs[2].Value()

	if kName == "" {
		m.statusMsg = "Key name cannot be empty"
		return m, nil
	}

	if err := domain.ValidateKeyName(kName); err != nil {
		m.statusMsg = err.Error()
		return m, nil
	}

	m.pendingKeys = append(m.pendingKeys, domain.APIKey{
		Title: kName,
		Key:   kName,
		Current: domain.SecretVersion{
			Value:     kVal,
			CreatedAt: time.Now(),
			CreatedBy: "tui",
		},
		History: []domain.SecretVersion{},
	})

	m.inputs[1].SetValue("")
	m.inputs[2].SetValue("")
	m.statusMsg = "Key added"
	return m, nil
}

func (m Model) saveNewProject() (tea.Model, tea.Cmd) {
	name := m.inputs[0].Value()
	env := m.selectedEnv.String()

	if name == "" {
		m.statusMsg = "Project name cannot be empty"
		return m, nil
	}

	if err := domain.ValidateProjectName(name); err != nil {
		m.statusMsg = err.Error()
		return m, nil
	}

	if len(m.pendingKeys) == 0 {
		m.statusMsg = "Add at least one key before saving"
		return m, nil
	}

	newProj := domain.Project{
		Name:        name,
		Environment: env,
		Keys:        m.pendingKeys,
	}

	if err := m.vault.CreateProject(newProj); err != nil {
		m.statusMsg = err.Error()
		return m, nil
	}

	if err := m.vault.Save(); err != nil {
		m.statusMsg = "ERROR: Failed to save"
		return m, nil
	}

	m.projects = m.vault.GetProjects()
	m.RefreshFiltered()
	m.currentView = ViewGrid
	m.state = StateNormal
	return m, nil
}

// updateEditProject handles the project edit view
func (m Model) updateEditProject(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	k := msg.String()

	if m.state == StateInsert {
		if k == m.keys.Back {
			m.state = StateNormal
			m.editProjectName.Blur()
			m.editProjectNewKey[0].Blur()
			m.editProjectNewKey[1].Blur()
			return m, nil
		}

		var cmd tea.Cmd
		switch m.editProjectFocus {
		case 0: // Project name
			m.editProjectName, cmd = m.editProjectName.Update(msg)
		case 2: // New key name
			m.editProjectNewKey[0], cmd = m.editProjectNewKey[0].Update(msg)
		case 3: // New key value
			m.editProjectNewKey[1], cmd = m.editProjectNewKey[1].Update(msg)
		}
		return m, cmd
	}

	switch k {
	case m.keys.Back, m.keys.Quit:
		m.currentView = ViewDetail
		m.state = StateNormal
		return m, nil

	case m.keys.Save:
		return m.saveProjectChanges()

	case m.keys.Add:
		return m.addKeyToProject()

	case m.keys.Delete:
		if m.editProjectFocus == 1 && len(m.activeProject.Keys) > 0 {
			if m.editProjectKeyIdx < len(m.activeProject.Keys) {
				keyName := m.activeProject.Keys[m.editProjectKeyIdx].Key
				m.confirmAction = ConfirmDeleteKey
				m.confirmMessage = "Delete key '" + keyName + "'?"
				m.previousView = ViewEditProject
				m.currentView = ViewConfirm
			}
		}
		return m, nil

	case m.keys.Search:
		if m.editProjectFocus == 0 || m.editProjectFocus == 2 || m.editProjectFocus == 3 {
			m.state = StateInsert
			switch m.editProjectFocus {
			case 0:
				return m, m.editProjectName.Focus()
			case 2:
				return m, m.editProjectNewKey[0].Focus()
			case 3:
				return m, m.editProjectNewKey[1].Focus()
			}
		}
		return m, nil

	case m.keys.Tab:
		m.editProjectFocus++
		if m.editProjectFocus > 5 {
			m.editProjectFocus = 0
		}
		return m, nil

	case m.keys.ShiftTab:
		m.editProjectFocus--
		if m.editProjectFocus < 0 {
			m.editProjectFocus = 5
		}
		return m, nil
	}

	if m.keys.IsNavigationDown(k) {
		if m.editProjectFocus == 1 && m.editProjectKeyIdx < len(m.activeProject.Keys)-1 {
			m.editProjectKeyIdx++
		} else {
			m.editProjectFocus++
			if m.editProjectFocus > 5 {
				m.editProjectFocus = 0
			}
		}
		return m, nil
	}

	if m.keys.IsNavigationUp(k) {
		if m.editProjectFocus == 1 && m.editProjectKeyIdx > 0 {
			m.editProjectKeyIdx--
		} else {
			m.editProjectFocus--
			if m.editProjectFocus < 0 {
				m.editProjectFocus = 5
			}
		}
		return m, nil
	}

	if k == m.keys.Enter {
		switch m.editProjectFocus {
		case 0, 2, 3:
			m.state = StateInsert
			switch m.editProjectFocus {
			case 0:
				return m, m.editProjectName.Focus()
			case 2:
				return m, m.editProjectNewKey[0].Focus()
			case 3:
				return m, m.editProjectNewKey[1].Focus()
			}
		case 4: // Add key button
			return m.addKeyToProject()
		case 5: // Save button
			return m.saveProjectChanges()
		}
	}

	return m, nil
}

func (m Model) addKeyToProject() (tea.Model, tea.Cmd) {
	kName := m.editProjectNewKey[0].Value()
	kVal := m.editProjectNewKey[1].Value()

	if kName == "" {
		m.statusMsg = "Key name cannot be empty"
		return m, nil
	}

	if err := domain.ValidateKeyName(kName); err != nil {
		m.statusMsg = err.Error()
		return m, nil
	}

	// Add key to vault
	newKey := domain.APIKey{
		Title: kName,
		Key:   kName,
		Current: domain.SecretVersion{
			Value:     kVal,
			CreatedAt: time.Now(),
			CreatedBy: "tui",
		},
		History: []domain.SecretVersion{},
	}

	if err := m.vault.AddKey(m.activeProject.Name, m.activeProject.Environment, newKey); err != nil {
		m.statusMsg = err.Error()
		return m, nil
	}

	if err := m.vault.Save(); err != nil {
		m.statusMsg = "ERROR: Failed to save"
		return m, nil
	}

	if proj, err := m.vault.GetProject(m.activeProject.Name, m.activeProject.Environment); err == nil {
		m.activeProject = *proj
	}
	m.projects = m.vault.GetProjects()
	m.RefreshFiltered()

	m.editProjectNewKey[0].SetValue("")
	m.editProjectNewKey[1].SetValue("")
	m.statusMsg = "Key added"
	return m, nil
}

// updateHistorySidebar handles the history sidebar view
func (m Model) updateHistorySidebar(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	k := msg.String()

	switch k {
	case m.keys.Back:
		m.historySidebarOpen = false
		return m, nil
	}

	if m.keys.IsNavigationUp(k) {
		if m.historyKeyIdx > 0 {
			m.historyKeyIdx--
		}
		return m, nil
	}
	if m.keys.IsNavigationDown(k) {
		if m.historyKeyIdx < len(m.activeProject.Keys)-1 {
			m.historyKeyIdx++
		}
		return m, nil
	}

	return m, nil
}

// saveProjectChanges: saves project name changes
func (m Model) saveProjectChanges() (tea.Model, tea.Cmd) {
	newName := m.editProjectName.Value()

	if newName == "" {
		m.statusMsg = "Project name cannot be empty"
		return m, nil
	}

	if err := domain.ValidateProjectName(newName); err != nil {
		m.statusMsg = err.Error()
		return m, nil
	}

	if newName != m.activeProject.Name {
		updatedProject := m.activeProject
		updatedProject.Name = newName

		if err := m.vault.UpdateProject(updatedProject); err != nil {
			m.statusMsg = err.Error()
			return m, nil
		}

		if err := m.vault.Save(); err != nil {
			m.statusMsg = "ERROR: Failed to save"
			return m, nil
		}

		m.activeProject.Name = newName
	}

	m.projects = m.vault.GetProjects()
	if proj, err := m.vault.GetProject(m.activeProject.Name, m.activeProject.Environment); err == nil {
		m.activeProject = *proj
	}
	m.RefreshFiltered()

	m.statusMsg = "Saved"
	m.currentView = ViewDetail
	m.state = StateNormal
	return m, nil
}
