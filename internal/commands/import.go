package commands

import (
	"fmt"
	"os"
	"strings"
	"time"

	"envy/internal/auth"
	"envy/internal/domain"
	"envy/internal/storage"

	"github.com/hashicorp/go-envparse"
)

func RunImport(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	envMap, err := envparse.Parse(file)
	if err != nil {
		fmt.Printf("Error parsing .env file: %v\n", err)
		return
	}

	if len(envMap) == 0 {
		fmt.Println("File is empty or contains no valid keys.")
		return
	}

	name, err := auth.PromptText("Enter Project Name: ")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	name = strings.TrimSpace(name)

	if err := domain.ValidateProjectName(name); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	env, err := auth.PromptText("Environment (prod/dev/stage) [dev]: ")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	env = strings.TrimSpace(env)
	if env == "" {
		env = domain.EnvDev
	}

	if err := domain.ValidateEnvironment(env); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	firstRun, err := storage.IsFirstRun()
	if err != nil {
		fmt.Printf("Error checking vault status: %v\n", err)
		return
	}

	var password string
	var projects []domain.Project
	var key []byte

	if firstRun {
		fmt.Println("No vault found. Creating new vault...")
		password, err = auth.PromptNewPassword()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if err := storage.Initialize(password); err != nil {
			fmt.Printf("Error initializing vault: %v\n", err)
			return
		}

		projects, key, err = storage.Load(password)
		if err != nil {
			fmt.Printf("Error loading vault: %v\n", err)
			return
		}
	} else {
		password, err = auth.PromptPassword("Enter master password: ")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		projects, key, err = storage.Load(password)
		if err != nil {
			fmt.Printf("Error loading vault: %v\n", err)
			return
		}
	}

	for i, p := range projects {
		if p.Name == name && p.Environment == env {
			fmt.Printf("Project '%s' (%s) already exists. Overwrite? [y/N]: ", name, env)
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" {
				fmt.Println("Import cancelled.")
				return
			}
			projects = append(projects[:i], projects[i+1:]...)
			break
		}
	}

	var newKeys []domain.APIKey
	for k, v := range envMap {
		if err := domain.ValidateKeyName(k); err != nil {
			fmt.Printf("Warning: Skipping invalid key '%s': %v\n", k, err)
			continue
		}

		newKeys = append(newKeys, domain.APIKey{
			Title: k,
			Key:   k,
			Current: domain.SecretVersion{
				Value:     v,
				CreatedAt: time.Now(),
				CreatedBy: "cli-import",
			},
			History: make([]domain.SecretVersion, 0),
		})
	}

	newProject := domain.Project{
		Name:        name,
		Environment: env,
		Keys:        newKeys,
	}

	projects = append(projects, newProject)

	if err := storage.Save(projects, key); err != nil {
		fmt.Printf("Error saving to vault: %v\n", err)
		return
	}

	fmt.Printf("\nSuccess! Imported project '%s' (%s) with %d keys.\n", name, env, len(newKeys))
	fmt.Println("Run 'envy' to view your keys in the TUI.")
}
