package auth

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ResolveGithubToken(cliToken string) (string, error) {
	if cliToken != "" {
		return cliToken, nil
	}
	if envToken := os.Getenv("GITHUB_TOKEN"); envToken != "" {
		return envToken, nil
	}
	ghToken, err := getGHAuthToken()
	if err == nil && ghToken != "" {
		return ghToken, nil
	}
	return "", fmt.Errorf("no GitHub token found")
}

func getGHAuthToken() (string, error) {
	_, err := exec.LookPath("gh")
	if err != nil {
		return "", fmt.Errorf("GitHub CLI (gh) not found in PATH")
	}

	cmd := exec.Command("gh", "auth", "token")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get GitHub token via gh CLI: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}
