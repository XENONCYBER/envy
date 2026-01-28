package tui

import (
	"strings"

	"envy/internal/config"
	"envy/internal/domain"
	"envy/internal/service"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type SessionState int

type ViewState int

const (
	StateNormal SessionState = iota
	StateInsert
)

const (
	ViewGrid ViewState = iota
	ViewDetail
	ViewCreate
	ViewEdit
	ViewEditProject
	ViewConfirm
)

type EnvOption int

type ConfirmAction int

const (
	ConfirmNone ConfirmAction = iota
	ConfirmDeleteProject
	ConfirmDeleteKey
)

const (
	EnvOptionDev EnvOption = iota
	EnvOptionProd
	EnvOptionStage
)

type SearchMode int

const (
	SearchAll SearchMode = iota
	SearchProjects
	SearchKeys
)

func (s SearchMode) String() string {
	switch s {
	case SearchProjects:
		return "Projects"
	case SearchKeys:
		return "Keys"
	default:
		return "All"
	}
}

func (s SearchMode) Next() SearchMode {
	switch s {
	case SearchAll:
		return SearchProjects
	case SearchProjects:
		return SearchKeys
	default:
		return SearchAll
	}
}

func (e EnvOption) String() string {
	switch e {
	case EnvOptionDev:
		return domain.EnvDev
	case EnvOptionProd:
		return domain.EnvProd
	case EnvOptionStage:
		return domain.EnvStage
	default:
		return domain.EnvDev
	}
}

type Model struct {
	vault    service.VaultService
	projects []domain.Project
	filtered []int

	searchInput textinput.Model
	searchMode  SearchMode
	keys        config.KeyMap
	styles      config.Styles

	state       SessionState
	currentView ViewState
	statusMsg   string

	selectedIdx  int
	scrollOffset int
	cols         int
	width        int
	height       int

	activeProject domain.Project
	detailCursor  int

	// form fields
	inputs      []textinput.Model // 0: project name, 1: key name, 2: key value
	focusIndex  int
	selectedEnv EnvOption
	pendingKeys []domain.APIKey

	// Edit mode fields
	editInput       textinput.Model
	editSidebarOpen bool

	revealedKey bool

	// History sidebar fields
	historySidebarOpen bool
	historyKeyIdx      int

	// Project edit mode fields
	editProjectName   textinput.Model
	editProjectKeyIdx int
	editProjectNewKey []textinput.Model // 0: key name, 1: key value
	editProjectFocus  int               // 0: name, 1: keys list, 2: new key name, 3: new key value, 4: add btn, 5: save btn

	// Confirmation dialog fields
	confirmAction  ConfirmAction
	confirmMessage string
	previousView   ViewState
}

func NewModel(projects []domain.Project, encryptionKey []byte, appConfig config.AppConfig) Model {
	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.Prompt = ""
	ti.Focus()
	ti.Width = 40

	var inputs []textinput.Model = make([]textinput.Model, 3)

	inputs[0] = textinput.New()
	inputs[0].Placeholder = "Project Name"
	inputs[0].Focus()
	inputs[0].Width = 40

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "Key Name (e.g. API_KEY)"
	inputs[1].Width = 40

	inputs[2] = textinput.New()
	inputs[2].Placeholder = "Value"
	inputs[2].Width = 40

	editInput := textinput.New()
	editInput.Placeholder = "New Value"
	editInput.Width = 50

	// Project edit inputs
	editProjectName := textinput.New()
	editProjectName.Placeholder = "Project Name"
	editProjectName.Prompt = ""
	editProjectName.Width = 40

	editProjectNewKey := make([]textinput.Model, 2)
	editProjectNewKey[0] = textinput.New()
	editProjectNewKey[0].Placeholder = "Key Name"
	editProjectNewKey[0].Prompt = ""
	editProjectNewKey[0].Width = 30
	editProjectNewKey[1] = textinput.New()
	editProjectNewKey[1].Placeholder = "Value"
	editProjectNewKey[1].Prompt = ""
	editProjectNewKey[1].Width = 30

	styles := config.NewStyles(appConfig.Theme)

	vault := service.NewVaultService(projects, encryptionKey)

	filteredIndices := make([]int, len(projects))
	for i := range projects {
		filteredIndices[i] = i
	}

	return Model{
		searchInput: ti,
		searchMode:  SearchAll,
		vault:       vault,
		projects:    projects,
		filtered:    filteredIndices,
		keys:        appConfig.Keys,
		styles:      styles,
		state:       StateNormal,
		currentView: ViewGrid,
		cols:        2,

		inputs:            inputs,
		focusIndex:        0,
		selectedEnv:       EnvOptionDev,
		editInput:         editInput,
		pendingKeys:       []domain.APIKey{},
		activeProject:     domain.Project{},
		editSidebarOpen:   false,
		editProjectName:   editProjectName,
		editProjectNewKey: editProjectNewKey,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) RefreshFiltered() {
	m.filtered = filterProjects(m.projects, m.searchInput.Value(), m.searchMode)
}

func filterProjects(all []domain.Project, query string, mode SearchMode) []int {
	if query == "" {
		indices := make([]int, len(all))
		for i := range all {
			indices[i] = i
		}
		return indices
	}

	var res []int
	q := strings.ToLower(query)

	for i, p := range all {
		matched := false

		if mode == SearchAll || mode == SearchProjects {
			if strings.Contains(strings.ToLower(p.Name), q) {
				matched = true
			}
		}

		if !matched && (mode == SearchAll || mode == SearchKeys) {
			for _, k := range p.Keys {
				if strings.Contains(strings.ToLower(k.Key), q) {
					matched = true
					break
				}
			}
		}

		if matched {
			res = append(res, i)
		}
	}
	return res
}

func (m *Model) GetFilteredProject(idx int) *domain.Project {
	if idx < 0 || idx >= len(m.filtered) {
		return nil
	}
	projectIdx := m.filtered[idx]
	if projectIdx < 0 || projectIdx >= len(m.projects) {
		return nil
	}
	return &m.projects[projectIdx]
}

func (m *Model) FilteredLen() int {
	return len(m.filtered)
}

func (m *Model) adjustScroll() {
	if len(m.filtered) == 0 {
		m.scrollOffset = 0
		return
	}

	cols := m.styles.GridCols
	visibleRows := m.styles.GridVisibleRows

	selectedRow := m.selectedIdx / cols

	if selectedRow >= m.scrollOffset+visibleRows {
		m.scrollOffset = selectedRow - visibleRows + 1
	}

	if selectedRow < m.scrollOffset {
		m.scrollOffset = selectedRow
	}

	totalRows := (len(m.filtered) + cols - 1) / cols
	maxScroll := totalRows - visibleRows
	if maxScroll < 0 {
		maxScroll = 0
	}
	if m.scrollOffset > maxScroll {
		m.scrollOffset = maxScroll
	}

	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}
}
