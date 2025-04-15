package fileops

import (
	"os"
	"testing"
)

func TestReadIni(t *testing.T) {
	// Test cases struct
	tests := []struct {
		name           string
		content        string
		expectedMap    map[string]string
		expectError    bool
		expectedIssues int
	}{
		{
			name: "valid ini file",
			content: `[section1]
key1 = value1
key2 = value2

[section2]
key3 = value3`,
			expectedMap: map[string]string{
				"section1::key1": "value1",
				"section1::key2": "value2",
				"section2::key3": "value3",
			},
			expectError:    false,
			expectedIssues: 0,
		},
		{
			name: "with comments and quoted values",
			content: `[section1]
# This is a comment
key1 = "quoted value"
key2 = 'single quoted'
key3 = value with # comment

[tricky_section]
key1='value1"more' # comment
key2="value' some more "
key3="#value some more'"
`,
			expectedMap: map[string]string{
				"section1::key1":       "quoted value",
				"section1::key2":       "single quoted",
				"section1::key3":       "value with",
				"tricky_section::key1": `value1"more`,
				"tricky_section::key2": `value' some more `,
				"tricky_section::key3": `#value some more'`,
			},
			expectError:    false,
			expectedIssues: 0,
		},
		{
			name: "with empty values",
			content: `[section1]
key1 = 
key2 = value2`,
			expectedMap: map[string]string{
				"section1::key2": "value2",
			},
			expectError:    false,
			expectedIssues: 0,
		},
		{
			name: "with initial blank line",
			content: `
[section1]
key1 = 
key2 = value2`,
			expectedMap: map[string]string{
				"section1::key2": "value2",
			},
			expectError:    false,
			expectedIssues: 0,
		},
		{
			name: "missing section",
			content: `key1 = value1
key2 = value2`,
			expectedMap:    map[string]string{},
			expectError:    true,
			expectedIssues: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpFile, err := os.CreateTemp("", "test_*.ini")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			// Write content to file
			if _, err := tmpFile.WriteString(tt.content); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			tmpFile.Close()

			// Test ReadIni function
			results, issues, err := ReadIni(tmpFile.Name())

			// Check error
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check issues length
			if len(issues) != tt.expectedIssues {
				t.Errorf("Expected %d issues, got %d", tt.expectedIssues, len(issues))
			}

			// Check results
			if len(results) != len(tt.expectedMap) {
				t.Errorf("Expected %d results, got %d", len(tt.expectedMap), len(results))
			}

			for k, v := range tt.expectedMap {
				if results[k] != v {
					t.Errorf("Expected %q=%q, got %q=%q", k, v, k, results[k])
				}
			}
		})
	}
}
