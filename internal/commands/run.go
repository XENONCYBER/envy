package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"envy/internal/auth"
	"envy/internal/config"
	"envy/internal/domain"
	"envy/internal/storage"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:                "run [project] -- [command]",
	Short:              "Run a command with secrets injected as environment variables",
	DisableFlagParsing: true, // disabled to handle flag manually
	Long: `Run a command with project secrets loaded into the environment.

This command fetches all secrets from the specified project and executes
the given command with those secrets available as environment variables.
The secrets are isolated to the subprocess and do not affect your shell.

Examples:
  envy run myproject -- npm start
  envy run production -- python app.py
  envy run dev -- make build
  envy run staging -- docker-compose up

The secrets are only available to the child process and are cleaned up
when the process exits.`,
	RunE: runWithSecrets,
}

func init() {
	RootCmd.AddCommand(runCmd)
}

func runWithSecrets(cmd *cobra.Command, args []string) error {
	// Flag parsing is disabled, so help command is handled manually
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			cmd.Help()
			return nil
		}
	}

	separatorIndex := -1
	for i, arg := range args {
		if arg == "--" {
			separatorIndex = i
			break
		}
	}

	if separatorIndex == -1 {
		return fmt.Errorf("missing '--' separator\n\nUsage: envy run [project] -- [command]\n\nExamples:\n  envy run myproject -- npm start\n  envy run production -- python app.py")
	}

	if separatorIndex == 0 {
		return fmt.Errorf("missing project name before '--'\n\nUsage: envy run [project] -- [command]")
	}

	if separatorIndex >= len(args)-1 {
		return fmt.Errorf("missing command after '--'\n\nUsage: envy run [project] -- [command]")
	}

	projectName := args[0]
	commandArgs := args[separatorIndex+1:]

	appConfig := config.LoadAppConfig()
	storage.SetConfig(appConfig.Backend)

	if err := config.EnsureDirectories(); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	firstRun, err := storage.IsFirstRun()
	if err != nil {
		return fmt.Errorf("failed to check vault status: %w", err)
	}

	if firstRun {
		return fmt.Errorf("no vault found. Please run 'envy' to create a vault first")
	}

	password, err := auth.PromptPassword("Enter master password: ")
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}

	projects, _, err := storage.Load(password)
	if err != nil {
		return fmt.Errorf("failed to load vault: %w", err)
	}

	// Project search is case-insensitive and first first match wins
	var project *domain.Project
	for i := range projects {
		if strings.EqualFold(projects[i].Name, projectName) {
			project = &projects[i]
			break
		}
	}

	if project == nil {
		return fmt.Errorf("project '%s' not found", projectName)
	}

	// Build environment map
	env := os.Environ()

	secretCount := 0
	for _, key := range project.Keys {
		envVar := fmt.Sprintf("%s=%s", key.Key, key.Current.Value)
		env = append(env, envVar)
		secretCount++
	}

	fmt.Fprintf(os.Stderr, "Loaded %d secrets from '%s' (%s)\n", secretCount, project.Name, project.Environment)
	fmt.Fprintf(os.Stderr, "Running: %s\n\n", formatCommand(commandArgs))

	return executeCommand(commandArgs, env)
}

func executeCommand(args []string, env []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command specified")
	}

	cmdName := args[0]
	cmdArgs := args[1:]

	command := exec.Command(cmdName, cmdArgs...)
	command.Env = env
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		return fmt.Errorf("command failed: %w", err)
	}

	return nil
}

// Helper function to format the command
func formatCommand(args []string) string {
	quoted := make([]string, len(args))
	for i, arg := range args {
		// Quote arguments with spaces
		if strings.Contains(arg, " ") {
			quoted[i] = fmt.Sprintf(`"%s"`, arg)
		} else {
			quoted[i] = arg
		}
	}
	return strings.Join(quoted, " ")
}
