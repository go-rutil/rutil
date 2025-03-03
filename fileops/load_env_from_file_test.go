package fileops

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnvFromFile(t *testing.T) {
	tests := []struct {
		name           string
		envContent     string
		expectedEnv    map[string]string
		expectedIssues int
		expectError    bool
	}{
		{
			name: "basic key-value pairs",
			envContent: `
				KEY1=value1
				KEY2=value2
			`,
			expectedEnv: map[string]string{
				"KEY1": "value1",
				"KEY2": "value2",
			},
			expectedIssues: 0,
			expectError:    false,
		},
		{
			name: "with comments and empty lines",
			envContent: `
				# This is a comment
				KEY1=value1

				# Another comment
				KEY2=value2
			`,
			expectedEnv: map[string]string{
				"KEY1": "value1",
				"KEY2": "value2",
			},
			expectedIssues: 0,
			expectError:    false,
		},
		{
			name: "with quoted values",
			envContent: `
				KEY1="quoted value"
				KEY2='single quoted'
			`,
			expectedEnv: map[string]string{
				"KEY1": "quoted value",
				"KEY2": "single quoted",
			},
			expectedIssues: 0,
			expectError:    false,
		},
		{
			name: "with inline comments",
			envContent: `
				KEY1=value1 # comment
				KEY2=value2#comment
			`,
			expectedEnv: map[string]string{
				"KEY1": "value1",
				"KEY2": "value2",
			},
			expectedIssues: 0,
			expectError:    false,
		},
		{
			name: "with significant chars in values",
			envContent: `
				KEY1='value1"more' # comment
				KEY2="value' some more "
				KEY3="#value some more'"
			`,
			expectedEnv: map[string]string{
				"KEY1": `value1"more`,
				"KEY2": `value' some more `,
				"KEY3": `#value some more'`,
			},
			expectedIssues: 0,
			expectError:    false,
		},
		{
			name: "with empty values",
			envContent: `
				KEY1=
				KEY2=value2
			`,
			expectedEnv: map[string]string{
				"KEY2": "value2",
			},
			expectedIssues: 1,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary env file
			tmpDir := t.TempDir()
			envFile := filepath.Join(tmpDir, "test.env")
			err := os.WriteFile(envFile, []byte(tt.envContent), 0644)
			if err != nil {
				t.Fatalf("failed to create test env file: %v", err)
			}

			// Clear any existing env vars
			for k := range tt.expectedEnv {
				os.Unsetenv(k)
			}

			// Run EnvFromFile
			issues, err := EnvFromFile(envFile)

			// Check error expectation
			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Check number of issues
			if len(issues) != tt.expectedIssues {
				t.Errorf("expected %d issues, got %d", tt.expectedIssues, len(issues))
			}

			// Check environment variables
			for key, expectedValue := range tt.expectedEnv {
				if value := os.Getenv(key); value != expectedValue {
					t.Errorf("environment variable %s = %s, want %s", key, value, expectedValue)
				}
			}
		})
	}

	// Test with non-existent file
	t.Run("non-existent file", func(t *testing.T) {
		_, err := EnvFromFile("non-existent-file.env")
		if err == nil {
			t.Error("expected error for non-existent file, got nil")
		}
	})
}
