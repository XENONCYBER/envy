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
	Long: `Envy is a secure, encrypted terminal-based secret manager with both CLI and TUI interfaces.

Quick shortcuts:
  envy -i file.env          Import .env file into vault
  envy -t project           Export project to .env file
  envy -s project KEY=VAL   Set a secret (alias for 'envy set')

For more options, use subcommands:
  envy set project KEY=VALUE -e prod
  envy run project -- command`,

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

		runTUI()
	},
}

func init() {
	// Root-level flags for quick operations
	RootCmd.Flags().StringVarP(&importFile, "import", "i", "", "Import .env file into vault")
	RootCmd.Flags().StringVarP(&exportProj, "export", "t", "", "Export project to .env file")
	RootCmd.Flags().BoolVar(&showVersion, "version", false, "Show version information")
}

// Execute is the main entry point for the CLI.
// It handles argument preprocessing and then runs Cobra.
func Execute() {
	args := preprocessArgs(os.Args)
	os.Args = args

	// Handle version flag early
	for _, arg := range os.Args[1:] {
		if arg == "--version" {
			fmt.Println("Envy v1.1.0") // TODO: Set actual version from build tags
			os.Exit(0)
		}
	}

	if err := RootCmd.Execute(); err != nil {
		if err.Error() != "pflag: help requested" {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

// preprocessArgs transforms shorthand flags into their subcommand equivalents.
// This allows us to use `-s project KEY=VALUE` as shorthand for `set project KEY=VALUE`.
// set command is not used with --set so cannot use normal cobra command structure.
func preprocessArgs(args []string) []string {
	if len(args) < 2 {
		return args
	}

	newArgs := []string{args[0]}
	i := 1

	for i < len(args) {
		arg := args[i]

		if arg == "-s" || arg == "--set" {
			newArgs = append(newArgs, "set")
			i++

			// Collect remaining args until we hit another flag or end
			for i < len(args) {
				nextArg := args[i]

				// Pass through other flags like -e, --env, -h, --help
				if nextArg == "-e" || nextArg == "--env" {
					newArgs = append(newArgs, nextArg)
					i++
					// Also add the value if present
					if i < len(args) && !isFlag(args[i]) {
						newArgs = append(newArgs, args[i])
						i++
					}
				} else if nextArg == "-h" || nextArg == "--help" {
					newArgs = append(newArgs, nextArg)
					i++
				} else if isFlag(nextArg) {
					// Unknown flag, pass it through
					newArgs = append(newArgs, nextArg)
					i++
				} else {
					// Positional argument (project name or KEY=VALUE)
					newArgs = append(newArgs, nextArg)
					i++
				}
			}
		} else {
			newArgs = append(newArgs, arg)
			i++
		}
	}

	return newArgs
}

// isFlag checks if an argument looks like a flag (starts with -)
func isFlag(arg string) bool {
	return len(arg) > 0 && arg[0] == '-'
}

func runTUI() {
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
}
