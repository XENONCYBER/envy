// Package auth handles password prompts and authentication
package auth

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func PromptPassword(prompt string) (string, error) {
	fmt.Print(prompt)

	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()

	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}

	password := strings.TrimSpace(string(bytePassword))
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	return password, nil
}

func PromptNewPassword() (string, error) {
	password, err := PromptPassword("Create master password: ")
	if err != nil {
		return "", err
	}

	if len(password) < 8 {
		return "", fmt.Errorf("password must be at least 8 characters")
	}

	confirm, err := PromptPassword("Confirm master password: ")
	if err != nil {
		return "", err
	}

	if password != confirm {
		return "", fmt.Errorf("passwords do not match")
	}

	return password, nil
}

func PromptText(prompt string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}
	return strings.TrimSpace(text), nil
}
