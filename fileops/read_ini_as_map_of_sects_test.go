package fileops

import (
	"os"
	"reflect"
	"testing"
)

func TestReadIniAsMapOfSections(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		content        string
		expectedMap    map[string]map[string]string
		expectError    bool
		expectedIssues int
	}{
		{
			name: "valid ini file with multiple sections",
			content: `[section1]
key1 = value1
key2 = value2

[section2]
key3 = value3
key4 = value4`,
			expectedMap: map[string]map[string]string{
				"section1": {
					"key1": "value1",
					"key2": "value2",
				},
				"section2": {
					"key3": "value3",
					"key4": "value4",
				},
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
			expectedMap: map[string]map[string]string{
				"section1": {
					"key1": "quoted value",
					"key2": "single quoted",
					"key3": "value with",
				},
				"tricky_section": {
					"key1": `value1"more`,
					"key2": `value' some more `,
					"key3": `#value some more'`,
				},
			},
			expectError:    false,
			expectedIssues: 0,
		},
		{
			name: "with empty values",
			content: `[section1]
key1 = 
key2 = value2`,
			expectedMap: map[string]map[string]string{
				"section1": {
					"key2": "value2",
				},
			},
			expectError:    false,
			expectedIssues: 1,
		},
		{
			name: "missing section",
			content: `key1 = value1
key2 = value2`,
			expectedMap:    map[string]map[string]string{},
			expectError:    true,
			expectedIssues: 0,
		},
		{
			name: "multiple sections with same key names",
			content: `[section1]
common = value1
unique1 = unique1

[section2]
common = value2
unique2 = unique2`,
			expectedMap: map[string]map[string]string{
				"section1": {
					"common":  "value1",
					"unique1": "unique1",
				},
				"section2": {
					"common":  "value2",
					"unique2": "unique2",
				},
			},
			expectError:    false,
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
			defer func() {
				_ = os.Remove(tmpFile.Name())
			}()

			// Write content to file
			if _, err := tmpFile.WriteString(tt.content); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			_ = tmpFile.Close() // close it so we can read it with our function

			// Test ReadIniAsMapOfSections function
			results, issues, err := ReadIniAsMapOfSections(tmpFile.Name())

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
			if !reflect.DeepEqual(results, tt.expectedMap) {
				t.Errorf("Results don't match expected map\nExpected: %v\nGot: %v", tt.expectedMap, results)
			}
		})
	}
}
