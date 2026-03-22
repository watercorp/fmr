package filedetails_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/watercorp/fmr/filedetails"
)

func TestFileDetails_New(t *testing.T) {
	// Define test file paths and expected results
	tests := []struct {
		filePath string
		expected filedetails.FileDetails
	}{
		{
			filePath: "/Users/username/Documents/doc.md",
			expected: filedetails.FileDetails{
				FullPath:  "/Users/username/Documents/doc.md",
				Name:      "doc.md",
				Directory: "/Users/username/Documents",
				Extension: ".md",
				BaseName:  "doc",
			},
		},
		{
			filePath: "/a.a",
			expected: filedetails.FileDetails{
				FullPath:  "/a.a",
				Name:      "a.a",
				Directory: "/",
				Extension: ".a",
				BaseName:  "a",
			},
		},
		{
			filePath: "/.config",
			expected: filedetails.FileDetails{
				FullPath:  "/.config",
				Name:      ".config",
				Directory: "/",
				Extension: "",
				BaseName:  ".config",
			},
		},
		{
			filePath: "/test.",
			expected: filedetails.FileDetails{
				FullPath:  "/test.",
				Name:      "test.",
				Directory: "/",
				Extension: ".",
				BaseName:  "test",
			},
		},
	}

	// Loop through each test case
	for i, tc := range tests {
		t.Run(fmt.Sprintf("new file details %d", i), func(t *testing.T) {
			// Create the new file details
			fd := filedetails.New(tc.filePath)

			// Compare against the expectation
			if diff := cmp.Diff(tc.expected, fd); diff != "" {
				t.Errorf("new file details do not match (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFileDetails_String(t *testing.T) {
	// Define the file paths
	tests := []string{
		"/Users/username/Documents/doc.md",
		"/test.md",
		".config",
		"test.",
	}

	// Loop through each test case
	for i, tc := range tests {
		t.Run(fmt.Sprintf("file details stringer %d", i), func(t *testing.T) {
			// Create the new file details
			fd := filedetails.New(tc)

			// Compare the input to the String function
			if tc != fd.String() {
				t.Errorf("input does not match string\nhave: %s\nwant:%s\n", fd.String(), tc)
			}
		})
	}
}
