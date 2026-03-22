package filedetails

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestFileDetails_internal_getFullExtension(t *testing.T) {
	// Define the file paths and the expected extensions
	tests := []struct {
		inputPath string
		expected  string
	}{
		{
			inputPath: "/Users/username/Documents/doc.md",
			expected:  ".md",
		},
		{
			inputPath: "/test.json",
			expected:  ".json",
		},
		{
			inputPath: "extract.tar.gz",
			expected:  ".tar.gz",
		},
		{
			inputPath: "/Users/username/.config",
			expected:  "",
		},
	}

	// Loop through each test case
	for i, tc := range tests {
		t.Run(fmt.Sprintf("file details get full extension %d", i), func(t *testing.T) {
			// Get the base path
			basePath := filepath.Base(tc.inputPath)

			// Get the full extension
			ext := getFullExtension(basePath)

			// Check that the extension matches the expected value
			if ext != tc.expected {
				t.Errorf("extensions do not match - have: %s want: %s\n", ext, tc.expected)
			}
		})
	}
}
