/*
	If there are no flags with the command then run the tui application.
	If flags are present execute functions and no need to run tui application.
*/

// Package commands stores all the command line commands for the application.
package commands

import (
	"fmt"
	"os"

	"envy/internal/auth"
	"envy/internal/config"
	"envy/internal/storage"
	"envy/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	importFile  string
	exportProj  string
	showVersion bool
)

var appConfig config.AppConfig

var RootCmd = &cobra.Command{
	Use:   "envy",
	Short: "Envy: Secure encrypted vault for API Keys and Secrets",
	Long:  `Envy is a secure, encrypted terminal-based secret manager with both CLI and TUI interfaces.`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		appConfig = config.LoadAppConfig()

		storage.SetConfig(appConfig.Backend)

		if err := config.EnsureDataDir(appConfig.Backend); err != nil {
			fmt.Printf("Warning: failed to create data directory: %v\n", err)
		}
	},

	Run: func(cmd *cobra.Command, args []string) {
		if importFile != "" {
			RunImport(importFile)
			return
		}

		if exportProj != "" {
			RunExport(exportProj)
			return
		}

		firstRun, err := storage.IsFirstRun()
		if err != nil {
			fmt.Printf("Error checking vault status: %v\n", err)
			os.Exit(1)
		}

		var password string

		if firstRun {
			fmt.Println("Welcome to Envy - Secure Secret Manager")
			fmt.Println("No vault found. Let's create one!")
			fmt.Println("")

			password, err = auth.PromptNewPassword()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			if err := storage.Initialize(password); err != nil {
				fmt.Printf("Error initializing vault: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Vault created successfully!")
			fmt.Println("")
		} else {
			password, err = auth.PromptPassword("Enter master password: ")
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		}

		projects, key, err := storage.Load(password)
		if err != nil {
			fmt.Printf("Error loading vault: %v\n", err)
			os.Exit(1)
		}

		p := tea.NewProgram(tui.NewModel(projects, key, appConfig), tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running TUI: %v\n", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	RootCmd.PersistentFlags().StringVarP(&importFile, "import", "i", "", "Import .env file into vault")
	RootCmd.PersistentFlags().StringVarP(&exportProj, "export", "t", "", "Export project to .env file")

	RootCmd.Flags().BoolVar(&showVersion, "version", false, "Show version information")

	if err := RootCmd.ParseFlags(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if showVersion {
		fmt.Println("Envy v1.0.0") // TODO: Set actual version from build tags
		os.Exit(0)
	}

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
