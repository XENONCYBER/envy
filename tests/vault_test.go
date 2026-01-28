package tests

import (
	"testing"
	"time"

	"envy/internal/domain"
	"envy/internal/service"
)

func createTestProject(name, env string, keys ...string) domain.Project {
	apiKeys := make([]domain.APIKey, len(keys))
	for i, key := range keys {
		apiKeys[i] = domain.APIKey{
			Title: key,
			Key:   key,
			Current: domain.SecretVersion{
				Value:     "secret-" + key,
				CreatedAt: time.Now(),
				CreatedBy: "test",
			},
			History: []domain.SecretVersion{},
		}
	}
	return domain.Project{
		Name:        name,
		Environment: env,
		Keys:        apiKeys,
	}
}

func TestNewVaultService(t *testing.T) {
	projects := []domain.Project{
		createTestProject("test-project", "dev", "API_KEY"),
	}
	key := []byte("0123456789abcdef0123456789abcdef")

	vault := service.NewVaultService(projects, key)

	if vault == nil {
		t.Fatal("NewVaultService() returned nil")
	}

	if len(vault.GetProjects()) != 1 {
		t.Errorf("GetProjects() returned %d projects, want 1", len(vault.GetProjects()))
	}
}

func TestGetProject(t *testing.T) {
	projects := []domain.Project{
		createTestProject("project1", "dev", "KEY1"),
		createTestProject("project1", "prod", "KEY2"),
		createTestProject("project2", "dev", "KEY3"),
	}
	vault := service.NewVaultService(projects, nil)

	// Find existing project
	proj, err := vault.GetProject("project1", "dev")
	if err != nil {
		t.Fatalf("GetProject() error: %v", err)
	}
	if proj.Name != "project1" || proj.Environment != "dev" {
		t.Errorf("GetProject() = %v, want project1/dev", proj)
	}

	// Find another environment
	proj, err = vault.GetProject("project1", "prod")
	if err != nil {
		t.Fatalf("GetProject() error: %v", err)
	}
	if proj.Environment != "prod" {
		t.Errorf("GetProject() environment = %q, want prod", proj.Environment)
	}

	// Non-existent project
	_, err = vault.GetProject("nonexistent", "dev")
	if err == nil {
		t.Error("GetProject() should return error for non-existent project")
	}
}

func TestCreateProject(t *testing.T) {
	vault := service.NewVaultService([]domain.Project{}, nil)

	proj := createTestProject("new-project", "dev", "API_KEY")
	err := vault.CreateProject(proj)
	if err != nil {
		t.Fatalf("CreateProject() error: %v", err)
	}

	if len(vault.GetProjects()) != 1 {
		t.Errorf("CreateProject() did not add project, count = %d", len(vault.GetProjects()))
	}

	// Try to create duplicate
	err = vault.CreateProject(proj)
	if err == nil {
		t.Error("CreateProject() should return error for duplicate project")
	}

	// Create same name but different environment (should succeed)
	projProd := createTestProject("new-project", "prod", "API_KEY")
	err = vault.CreateProject(projProd)
	if err != nil {
		t.Errorf("CreateProject() should allow same name with different env: %v", err)
	}
}

func TestCreateProjectValidation(t *testing.T) {
	vault := service.NewVaultService([]domain.Project{}, nil)

	// Invalid project name
	proj := createTestProject("", "dev", "KEY")
	err := vault.CreateProject(proj)
	if err == nil {
		t.Error("CreateProject() should reject empty project name")
	}

	// Invalid environment
	proj = createTestProject("valid-name", "invalid", "KEY")
	err = vault.CreateProject(proj)
	if err == nil {
		t.Error("CreateProject() should reject invalid environment")
	}

	// Invalid key name
	proj = domain.Project{
		Name:        "valid-name",
		Environment: "dev",
		Keys: []domain.APIKey{
			{Key: "INVALID=KEY", Current: domain.SecretVersion{Value: "test"}},
		},
	}
	err = vault.CreateProject(proj)
	if err == nil {
		t.Error("CreateProject() should reject invalid key name")
	}
}

func TestDeleteProject(t *testing.T) {
	projects := []domain.Project{
		createTestProject("project1", "dev", "KEY1"),
		createTestProject("project2", "dev", "KEY2"),
	}
	vault := service.NewVaultService(projects, nil)

	// Delete existing
	err := vault.DeleteProject("project1", "dev")
	if err != nil {
		t.Fatalf("DeleteProject() error: %v", err)
	}

	if len(vault.GetProjects()) != 1 {
		t.Errorf("DeleteProject() did not remove project, count = %d", len(vault.GetProjects()))
	}

	// Verify the right project was deleted
	_, err = vault.GetProject("project1", "dev")
	if err == nil {
		t.Error("project1 should have been deleted")
	}

	// Delete non-existent
	err = vault.DeleteProject("nonexistent", "dev")
	if err == nil {
		t.Error("DeleteProject() should return error for non-existent project")
	}
}

func TestAddKey(t *testing.T) {
	vault := service.NewVaultService([]domain.Project{
		createTestProject("project", "dev"),
	}, nil)

	key := domain.APIKey{
		Title: "NEW_KEY",
		Key:   "NEW_KEY",
		Current: domain.SecretVersion{
			Value:     "secret",
			CreatedAt: time.Now(),
			CreatedBy: "test",
		},
	}

	err := vault.AddKey("project", "dev", key)
	if err != nil {
		t.Fatalf("AddKey() error: %v", err)
	}

	proj, _ := vault.GetProject("project", "dev")
	if len(proj.Keys) != 1 {
		t.Errorf("AddKey() did not add key, count = %d", len(proj.Keys))
	}

	// Add duplicate key
	err = vault.AddKey("project", "dev", key)
	if err == nil {
		t.Error("AddKey() should return error for duplicate key")
	}

	// Add to non-existent project
	err = vault.AddKey("nonexistent", "dev", key)
	if err == nil {
		t.Error("AddKey() should return error for non-existent project")
	}
}

func TestUpdateKey(t *testing.T) {
	vault := service.NewVaultService([]domain.Project{
		createTestProject("project", "dev", "API_KEY"),
	}, nil)

	// Get original value
	proj, _ := vault.GetProject("project", "dev")
	originalValue := proj.Keys[0].Current.Value

	// Update
	err := vault.UpdateKey("project", "dev", "API_KEY", "new-value")
	if err != nil {
		t.Fatalf("UpdateKey() error: %v", err)
	}

	// Verify update
	proj, _ = vault.GetProject("project", "dev")
	if proj.Keys[0].Current.Value != "new-value" {
		t.Errorf("UpdateKey() did not update value, got %q", proj.Keys[0].Current.Value)
	}

	// Verify history
	if len(proj.Keys[0].History) != 1 {
		t.Errorf("UpdateKey() did not preserve history, count = %d", len(proj.Keys[0].History))
	}
	if proj.Keys[0].History[0].Value != originalValue {
		t.Errorf("UpdateKey() history value = %q, want %q", proj.Keys[0].History[0].Value, originalValue)
	}

	// Update non-existent key
	err = vault.UpdateKey("project", "dev", "NONEXISTENT", "value")
	if err == nil {
		t.Error("UpdateKey() should return error for non-existent key")
	}
}

func TestDeleteKey(t *testing.T) {
	vault := service.NewVaultService([]domain.Project{
		createTestProject("project", "dev", "KEY1", "KEY2"),
	}, nil)

	err := vault.DeleteKey("project", "dev", "KEY1")
	if err != nil {
		t.Fatalf("DeleteKey() error: %v", err)
	}

	proj, _ := vault.GetProject("project", "dev")
	if len(proj.Keys) != 1 {
		t.Errorf("DeleteKey() did not remove key, count = %d", len(proj.Keys))
	}

	// Delete non-existent key
	err = vault.DeleteKey("project", "dev", "NONEXISTENT")
	if err == nil {
		t.Error("DeleteKey() should return error for non-existent key")
	}
}

func TestGetEncryptionKey(t *testing.T) {
	expectedKey := []byte("0123456789abcdef0123456789abcdef")
	vault := service.NewVaultService([]domain.Project{}, expectedKey)

	key := vault.GetEncryptionKey()
	if string(key) != string(expectedKey) {
		t.Errorf("GetEncryptionKey() = %v, want %v", key, expectedKey)
	}
}
