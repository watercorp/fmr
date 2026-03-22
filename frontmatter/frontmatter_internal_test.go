package frontmatter

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFrontmatter_internal_cleanMap(t *testing.T) {
	// Define test cases
	tests := []struct {
		input    map[string]any
		expected map[string]any
	}{
		{
			input: map[string]any{
				"one": "value 1",
				"two": "value 2",
			},
			expected: map[string]any{
				"one": "value 1",
				"two": "value 2",
			},
		},
		{
			input: map[string]any{
				"one": nil,
				"two": "value 2",
			},
			expected: map[string]any{
				"two": "value 2",
			},
		},
		{
			input: map[string]any{
				"one": map[string]any{
					"one": "value 1",
					"two": nil,
				},
				"two": "value 2",
			},
			expected: map[string]any{
				"one": map[string]any{
					"one": "value 1",
				},
				"two": "value 2",
			},
		},
		{
			input: map[string]any{
				"one": map[string]any{
					"one": nil,
				},
				"two": "value 2",
			},
			expected: map[string]any{
				"two": "value 2",
			},
		},
	}

	// Loop through each test case
	for i, tc := range tests {
		t.Run(fmt.Sprintf("clean map %d", i), func(t *testing.T) {
			// Call the cleanMap function
			output := cleanMap(tc.input)

			// Confirm the values match
			if diff := cmp.Diff(output, tc.expected); diff != "" {
				t.Errorf("values do not match (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFrontmatter_internal_isEmpty(t *testing.T) {
	// Create a test pointer
	var strPointer *string

	// Define test cases
	tests := []struct {
		input    any
		expected bool
	}{
		{
			input:    []string{"one", "two", "three"},
			expected: false,
		},
		{
			input:    []string{},
			expected: true,
		},
		{
			input:    []string{"one", "two", "three"}[:0],
			expected: true,
		},
		{
			input: map[string]any{
				"one": "value 1",
			},
			expected: false,
		},
		{
			input:    map[string]any{},
			expected: true,
		},
		{
			input:    strPointer,
			expected: true,
		},
		{
			input:    10,
			expected: false,
		},
	}

	// Loop through test cases
	for i, tc := range tests {
		t.Run(fmt.Sprintf("is empty %d", i), func(t *testing.T) {
			// Call isEmpty function
			output := isEmpty(tc.input)

			// Confirm values match
			if output != tc.expected {
				t.Errorf("values do not match, have: %t want: %t", output, tc.expected)
			}
		})
	}
}
