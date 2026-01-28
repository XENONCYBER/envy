// Package config handles application configuration via Lua
package config

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetDefaultDataDir returns the default directory for storing application data (keys.json, .lock)
// - Linux: ~/.envy/
// - macOS: ~/.envy/
// - Windows: %APPDATA%\envy\
func GetDefaultDataDir() string {
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			home, _ := os.UserHomeDir()
			appData = filepath.Join(home, "AppData", "Roaming")
		}
		return filepath.Join(appData, "envy")
	default:
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".envy")
	}
}

// GetDefaultConfigDir returns the default directory for storing configuration files (config.lua)
// - Linux: ~/.config/envy/
// - macOS: ~/Library/Application Support/envy/
// - Windows: %APPDATA%\envy\
func GetDefaultConfigDir() string {
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			home, _ := os.UserHomeDir()
			appData = filepath.Join(home, "AppData", "Roaming")
		}
		return filepath.Join(appData, "envy")
	case "darwin":
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "Library", "Application Support", "envy")
	default:
		// Respect XDG_CONFIG_HOME if set
		xdgConfig := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfig != "" {
			return filepath.Join(xdgConfig, "envy")
		}
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".config", "envy")
	}
}

func GetDefaultKeysPath() string {
	return filepath.Join(GetDefaultDataDir(), "keys.json")
}

func GetDefaultLockPath() string {
	return filepath.Join(GetDefaultDataDir(), ".lock")
}

func GetDefaultConfigPath() string {
	return filepath.Join(GetDefaultConfigDir(), "config.lua")
}

func EnsureDirectories() error {
	dataDir := GetDefaultDataDir()
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		return err
	}

	configDir := GetDefaultConfigDir()
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return err
	}

	return nil
}
