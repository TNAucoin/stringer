package auth

import (
	"os"
	"os/exec"
	"testing"
)

func TestResolveGithubToken(t *testing.T) {
	tests := []struct {
		name          string
		cliToken      string
		envToken      string
		setupMockGh   bool
		mockGhOutput  string
		mockGhError   bool
		expectedToken string
		expectError   bool
	}{
		{
			name:          "CLI token provided",
			cliToken:      "cli-token-123",
			envToken:      "",
			setupMockGh:   false,
			expectedToken: "cli-token-123",
			expectError:   false,
		},
		{
			name:          "Environment variable token",
			cliToken:      "",
			envToken:      "env-token-456",
			setupMockGh:   false,
			expectedToken: "env-token-456",
			expectError:   false,
		},
		{
			name:          "GitHub CLI token",
			cliToken:      "",
			envToken:      "",
			setupMockGh:   true,
			mockGhOutput:  "gh-token-789\n",
			mockGhError:   false,
			expectedToken: "gh-token-789",
			expectError:   false,
		},
		{
			name:          "GitHub CLI not found",
			cliToken:      "",
			envToken:      "",
			setupMockGh:   false, // This will cause LookPath to fail
			expectedToken: "",
			expectError:   true,
		},
		{
			name:          "GitHub CLI error",
			cliToken:      "",
			envToken:      "",
			setupMockGh:   true,
			mockGhOutput:  "",
			mockGhError:   true,
			expectedToken: "",
			expectError:   true,
		},
		{
			name:          "No token available",
			cliToken:      "",
			envToken:      "",
			setupMockGh:   true,
			mockGhOutput:  "", // Empty output
			mockGhError:   false,
			expectedToken: "",
			expectError:   true,
		},
	}

	// Save original functions to restore later
	originalLookPath := lookupGHPath
	originalCommand := execGHCommand
	// Restore original functions when test completes
	defer func() {
		lookupGHPath = originalLookPath
		execGHCommand = originalCommand
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original env var to restore later
			originalEnv, envVarExists := os.LookupEnv("GITHUB_TOKEN")
			defer func() {
				if envVarExists {
					os.Setenv("GITHUB_TOKEN", originalEnv)
				} else {
					os.Unsetenv("GITHUB_TOKEN")
				}
			}()

			// Set up environment variable if needed
			if tt.envToken != "" {
				os.Setenv("GITHUB_TOKEN", tt.envToken)
			} else {
				os.Unsetenv("GITHUB_TOKEN")
			}

			// Mock exec.LookPath
			lookupGHPath = func(file string) (string, error) {
				if file == "gh" && tt.setupMockGh {
					return "/usr/local/bin/gh", nil
				}
				return "", exec.ErrNotFound
			}

			// Mock exec.Command
			execGHCommand = func(command string, args ...string) *exec.Cmd {
				if command == "gh" && args[0] == "auth" && args[1] == "token" {
					if tt.mockGhError {
						return exec.Command("false") // This will cause an error when Output() is called
					}
					return exec.Command("echo", tt.mockGhOutput) // This will echo the mock output
				}
				return exec.Command("echo", "unexpected command")
			}

			// Call the function being tested
			token, err := ResolveGithubToken(tt.cliToken)

			// Check results
			if tt.expectError && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
			if token != tt.expectedToken {
				t.Errorf("expected token %q but got %q", tt.expectedToken, token)
			}
		})
	}
}
