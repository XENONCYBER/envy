// Package storage handles encrypted storage of projects and API keys using AES-256-GCM encryption
package storage

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"envy/internal/config"
	"envy/internal/crypto"
	"envy/internal/domain"
)

const (
	schemaVersion = 1
)

// Store configuration: set by main before use
var storeConfig config.BackendConfig

func SetConfig(cfg config.BackendConfig) {
	storeConfig = cfg
}

func getLockPath() string {
	if storeConfig.LockPath != "" {
		return storeConfig.LockPath
	}
	return config.GetDefaultLockPath()
}

func getStorePath() string {
	if storeConfig.KeysPath != "" {
		return storeConfig.KeysPath
	}
	return config.GetDefaultKeysPath()
}

func IsFirstRun() (bool, error) {
	path := getStorePath()

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check if file exists: %w", err)
	}
	return false, nil
}

func Initialize(password string) error {
	salt, err := crypto.GenerateSalt()
	if err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	key := crypto.DeriveKey(password, salt)

	authHash := crypto.GenerateAuthHash(key)

	store := domain.Store{
		Version:  schemaVersion,
		Salt:     base64.StdEncoding.EncodeToString(salt),
		AuthHash: authHash,
		Projects: []domain.Project{},
	}

	path := getStorePath()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	return saveStore(store)
}

func saveStore(store domain.Store) error {
	lockPath := getLockPath()
	if err := os.MkdirAll(filepath.Dir(lockPath), 0o700); err != nil {
		return fmt.Errorf("failed to create lock directory: %w", err)
	}

	lock, err := AcquireLock(lockPath)
	if err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}
	defer lock.Release()
	return saveStoreUnlocked(store)
}

func Load(password string) ([]domain.Project, []byte, error) {
	path := getStorePath()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []domain.Project{}, nil, nil
		}
		return nil, nil, fmt.Errorf("failed to read storage file: %w", err)
	}

	var store domain.Store
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, nil, fmt.Errorf("failed to parse storage file (corrupted?): %w", err)
	}

	salt, err := base64.StdEncoding.DecodeString(store.Salt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode salt: %w", err)
	}

	key := crypto.DeriveKey(password, salt)

	if !crypto.VerifyAuthHash(key, store.AuthHash) {
		return nil, nil, fmt.Errorf("authentication failed: incorrect password")
	}

	decryptedProjects, err := decryptSecrets(store.Projects, key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decrypt secrets: %w", err)
	}

	return decryptedProjects, key, nil
}

func Save(projects []domain.Project, key []byte) error {
	lockPath := getLockPath()
	if err := os.MkdirAll(filepath.Dir(lockPath), 0o700); err != nil {
		return fmt.Errorf("failed to create lock directory: %w", err)
	}

	lock, err := AcquireLock(lockPath)
	if err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}
	defer lock.Release()

	path := getStorePath()

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read existing storage file: %w", err)
	}

	var existingStore domain.Store
	if err := json.Unmarshal(data, &existingStore); err != nil {
		return fmt.Errorf("failed to parse storage file: %w", err)
	}

	encryptedProjects, err := encryptSecrets(projects, key)
	if err != nil {
		return fmt.Errorf("failed to encrypt secrets: %w", err)
	}

	store := domain.Store{
		Version:  schemaVersion,
		Salt:     existingStore.Salt,
		AuthHash: existingStore.AuthHash,
		Projects: encryptedProjects,
	}

	return saveStoreUnlocked(store)
}

func encryptSecrets(projects []domain.Project, key []byte) ([]domain.Project, error) {
	encrypted := make([]domain.Project, len(projects))

	for i, project := range projects {
		encryptedKeys := make([]domain.APIKey, len(project.Keys))

		for j, apiKey := range project.Keys {
			encryptedCurrent, err := crypto.Encrypt([]byte(apiKey.Current.Value), key)
			if err != nil {
				return nil, fmt.Errorf("failed to encrypt current value for %s.%s: %w",
					project.Name, apiKey.Key, err)
			}

			encryptedHistory := make([]domain.SecretVersion, len(apiKey.History))
			for k, historyVersion := range apiKey.History {
				encryptedValue, err := crypto.Encrypt([]byte(historyVersion.Value), key)
				if err != nil {
					return nil, fmt.Errorf("failed to encrypt history value for %s.%s: %w",
						project.Name, apiKey.Key, err)
				}
				encryptedHistory[k] = domain.SecretVersion{
					Value:     encryptedValue,
					CreatedAt: historyVersion.CreatedAt,
					CreatedBy: historyVersion.CreatedBy,
				}
			}

			encryptedKeys[j] = domain.APIKey{
				Title: apiKey.Title,
				Key:   apiKey.Key,
				Current: domain.SecretVersion{
					Value:     encryptedCurrent,
					CreatedAt: apiKey.Current.CreatedAt,
					CreatedBy: apiKey.Current.CreatedBy,
				},
				History: encryptedHistory,
			}
		}

		encrypted[i] = domain.Project{
			Name:        project.Name,
			Environment: project.Environment,
			Keys:        encryptedKeys,
		}
	}

	return encrypted, nil
}

func decryptSecrets(projects []domain.Project, key []byte) ([]domain.Project, error) {
	decrypted := make([]domain.Project, len(projects))

	for i, project := range projects {
		decryptedKeys := make([]domain.APIKey, len(project.Keys))

		for j, apiKey := range project.Keys {
			decryptedCurrent, err := crypto.Decrypt(apiKey.Current.Value, key)
			if err != nil {
				return nil, fmt.Errorf("failed to decrypt current value for %s.%s: %w",
					project.Name, apiKey.Key, err)
			}

			decryptedHistory := make([]domain.SecretVersion, len(apiKey.History))
			for k, historyVersion := range apiKey.History {
				decryptedValue, err := crypto.Decrypt(historyVersion.Value, key)
				if err != nil {
					return nil, fmt.Errorf("failed to decrypt history value for %s.%s: %w",
						project.Name, apiKey.Key, err)
				}
				decryptedHistory[k] = domain.SecretVersion{
					Value:     string(decryptedValue),
					CreatedAt: historyVersion.CreatedAt,
					CreatedBy: historyVersion.CreatedBy,
				}
			}

			decryptedKeys[j] = domain.APIKey{
				Title: apiKey.Title,
				Key:   apiKey.Key,
				Current: domain.SecretVersion{
					Value:     string(decryptedCurrent),
					CreatedAt: apiKey.Current.CreatedAt,
					CreatedBy: apiKey.Current.CreatedBy,
				},
				History: decryptedHistory,
			}
		}

		decrypted[i] = domain.Project{
			Name:        project.Name,
			Environment: project.Environment,
			Keys:        decryptedKeys,
		}
	}

	return decrypted, nil
}

func saveStoreUnlocked(store domain.Store) error {
	path := getStorePath()

	storeJSON, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal store: %w", err)
	}

	tempPath := path + ".tmp"

	if err := os.WriteFile(tempPath, storeJSON, 0o600); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := os.Rename(tempPath, path); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

func CreateBackup() error {
	path := getStorePath()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read vault for backup: %w", err)
	}

	backupPath := fmt.Sprintf("%s.backup", path)
	if err := os.WriteFile(backupPath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write backup: %w", err)
	}

	return nil
}
