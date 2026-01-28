// Package domain: Stores all the structs necessary to store data and render the view.
package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	EnvProd  = "prod"
	EnvStage = "stage"
	EnvDev   = "dev"
)

type SecretVersion struct {
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
}

type APIKey struct {
	Title   string          `json:"title"`
	Key     string          `json:"key"`
	Current SecretVersion   `json:"current"`
	History []SecretVersion `json:"history"`
}

type Project struct {
	Name        string   `json:"name"`
	Environment string   `json:"environment"`
	Keys        []APIKey `json:"keys"`
}

type Store struct {
	Version  int       `json:"version"`
	Salt     string    `json:"salt"`
	AuthHash string    `json:"auth_hash"`
	Projects []Project `json:"projects"`
}

func ValidateProjectName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("project name cannot be empty")
	}
	if len(name) > 256 {
		return errors.New("project name too long (max 256 characters)")
	}
	return nil
}

func ValidateEnvironment(env string) error {
	env = strings.TrimSpace(env)
	if env != EnvProd && env != EnvDev && env != EnvStage {
		return fmt.Errorf("invalid environment '%s' (must be prod, dev, or stage)", env)
	}
	return nil
}

func ValidateKeyName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("key name cannot be empty")
	}
	if strings.ContainsAny(name, "=\n\r") {
		return errors.New("key name cannot contain =, newline, or carriage return")
	}
	if len(name) > 256 {
		return errors.New("key name too long (max 256 characters)")
	}
	return nil
}
