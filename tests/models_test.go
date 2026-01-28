package tests

import (
	"strings"
	"testing"

	"envy/internal/domain"
)

func TestValidateProjectName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid name", "my-project", false},
		{"valid with spaces", "  my-project  ", false},
		{"empty", "", true},
		{"only spaces", "   ", true},
		{"too long", strings.Repeat("a", 257), true},
		{"max length", strings.Repeat("a", 256), false},
		{"special chars", "project_123-test", false},
		{"unicode", "项目名称", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := domain.ValidateProjectName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateProjectName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateEnvironment(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"prod", "prod", false},
		{"dev", "dev", false},
		{"stage", "stage", false},
		{"prod with spaces", "  prod  ", false},
		{"invalid", "production", true},
		{"empty", "", true},
		{"uppercase", "PROD", true},
		{"mixed case", "Prod", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := domain.ValidateEnvironment(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEnvironment(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateKeyName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid key", "API_KEY", false},
		{"valid with spaces trimmed", "  API_KEY  ", false},
		{"empty", "", true},
		{"only spaces", "   ", true},
		{"contains equals", "KEY=value", true},
		{"contains newline", "KEY\nNAME", true},
		{"contains carriage return", "KEY\rNAME", true},
		{"too long", strings.Repeat("a", 257), true},
		{"max length", strings.Repeat("a", 256), false},
		{"lowercase", "api_key", false},
		{"mixed case", "ApiKey", false},
		{"with numbers", "KEY_123", false},
		{"just numbers", "12345", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := domain.ValidateKeyName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateKeyName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestEnvironmentConstants(t *testing.T) {
	if domain.EnvProd != "prod" {
		t.Errorf("EnvProd = %q, want %q", domain.EnvProd, "prod")
	}
	if domain.EnvStage != "stage" {
		t.Errorf("EnvStage = %q, want %q", domain.EnvStage, "stage")
	}
	if domain.EnvDev != "dev" {
		t.Errorf("EnvDev = %q, want %q", domain.EnvDev, "dev")
	}
}
