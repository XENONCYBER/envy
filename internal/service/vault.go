// Package service provides the business logic layer between TUI and storage
package service

import (
	"fmt"
	"time"

	"envy/internal/domain"
	"envy/internal/storage"
)

// VaultService defines the interface for vault operations
// This interface enables mocking for tests
type VaultService interface {
	GetProjects() []domain.Project
	GetProject(name, env string) (*domain.Project, error)
	CreateProject(project domain.Project) error
	UpdateProject(project domain.Project) error
	DeleteProject(name, env string) error

	AddKey(projectName, projectEnv string, key domain.APIKey) error
	UpdateKey(projectName, projectEnv, keyName string, newValue string) error
	DeleteKey(projectName, projectEnv, keyName string) error

	Save() error

	GetEncryptionKey() []byte
}

// Implements VaultService
type vaultService struct {
	projects      []domain.Project
	encryptionKey []byte
}

func NewVaultService(projects []domain.Project, encryptionKey []byte) VaultService {
	return &vaultService{
		projects:      projects,
		encryptionKey: encryptionKey,
	}
}

func (v *vaultService) GetProjects() []domain.Project {
	return v.projects
}

func (v *vaultService) GetProject(name, env string) (*domain.Project, error) {
	for i := range v.projects {
		if v.projects[i].Name == name && v.projects[i].Environment == env {
			return &v.projects[i], nil
		}
	}
	return nil, fmt.Errorf("project '%s' (%s) not found", name, env)
}

func (v *vaultService) CreateProject(project domain.Project) error {
	if err := domain.ValidateProjectName(project.Name); err != nil {
		return err
	}
	if err := domain.ValidateEnvironment(project.Environment); err != nil {
		return err
	}

	for _, p := range v.projects {
		if p.Name == project.Name && p.Environment == project.Environment {
			return fmt.Errorf("project '%s' (%s) already exists", project.Name, project.Environment)
		}
	}

	for _, key := range project.Keys {
		if err := domain.ValidateKeyName(key.Key); err != nil {
			return fmt.Errorf("invalid key '%s': %w", key.Key, err)
		}
	}

	v.projects = append(v.projects, project)
	return nil
}

func (v *vaultService) UpdateProject(project domain.Project) error {
	for i, p := range v.projects {
		if p.Name == project.Name && p.Environment == project.Environment {
			v.projects[i] = project
			return nil
		}
	}
	return fmt.Errorf("project '%s' (%s) not found", project.Name, project.Environment)
}

func (v *vaultService) DeleteProject(name, env string) error {
	for i, p := range v.projects {
		if p.Name == name && p.Environment == env {
			v.projects[i] = v.projects[len(v.projects)-1]
			v.projects = v.projects[:len(v.projects)-1]
			return nil
		}
	}
	return fmt.Errorf("project '%s' (%s) not found", name, env)
}

func (v *vaultService) AddKey(projectName, projectEnv string, key domain.APIKey) error {
	if err := domain.ValidateKeyName(key.Key); err != nil {
		return err
	}

	project, err := v.GetProject(projectName, projectEnv)
	if err != nil {
		return err
	}

	for _, k := range project.Keys {
		if k.Key == key.Key {
			return fmt.Errorf("key '%s' already exists in project", key.Key)
		}
	}

	project.Keys = append(project.Keys, key)
	return nil
}

func (v *vaultService) UpdateKey(projectName, projectEnv, keyName string, newValue string) error {
	project, err := v.GetProject(projectName, projectEnv)
	if err != nil {
		return err
	}

	for i, key := range project.Keys {
		if key.Key == keyName {
			project.Keys[i].History = append(project.Keys[i].History, key.Current)

			project.Keys[i].Current = domain.SecretVersion{
				Value:     newValue,
				CreatedAt: time.Now(),
				CreatedBy: "tui-edit",
			}
			return nil
		}
	}

	return fmt.Errorf("key '%s' not found in project '%s' (%s)", keyName, projectName, projectEnv)
}

func (v *vaultService) DeleteKey(projectName, projectEnv, keyName string) error {
	project, err := v.GetProject(projectName, projectEnv)
	if err != nil {
		return err
	}

	for i, key := range project.Keys {
		if key.Key == keyName {
			project.Keys[i] = project.Keys[len(project.Keys)-1]
			project.Keys = project.Keys[:len(project.Keys)-1]
			return nil
		}
	}

	return fmt.Errorf("key '%s' not found in project '%s' (%s)", keyName, projectName, projectEnv)
}

// Save persists all projects to storage
func (v *vaultService) Save() error {
	return storage.Save(v.projects, v.encryptionKey)
}

func (v *vaultService) GetEncryptionKey() []byte {
	return v.encryptionKey
}
