package commands

import (
	"fmt"
	"os"
	"strings"

	"envy/internal/auth"
	"envy/internal/storage"
)

func RunExport(projectName string) {
	password, err := auth.PromptPassword("Enter master password: ")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	projects, _, err := storage.Load(password)
	if err != nil {
		fmt.Printf("Error loading vault: %v\n", err)
		return
	}

	foundIdx := -1
	for i, p := range projects {
		if strings.EqualFold(p.Name, projectName) {
			foundIdx = i
			break
		}
	}

	if foundIdx < 0 {
		fmt.Printf("Project '%s' not found in vault.\n", projectName)
		return
	}

	foundProject := &projects[foundIdx]

	var content strings.Builder
	content.WriteString(fmt.Sprintf("# Exported from Envy - Project: %s (%s)\n", foundProject.Name, foundProject.Environment))

	for _, key := range foundProject.Keys {
		line := fmt.Sprintf("%s=%s\n", key.Key, key.Current.Value)
		content.WriteString(line)
	}

	fileName := ".env"
	if err := os.WriteFile(fileName, []byte(content.String()), 0o600); err != nil {
		fmt.Printf("Error writing .env file: %v\n", err)
		return
	}
	fmt.Println("Warning: Exported .env file contains secrets in plain text. Keep it secure!")
	fmt.Printf("Exported project '%s' to .env\n", foundProject.Name)
}
