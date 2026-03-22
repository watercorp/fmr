package frontmatter_test

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/watercorp/fmr/frontmatter"
)

func TestFrontmatter_New(t *testing.T) {
	// Define file paths
	filePaths := []struct {
		inputFilePath    string
		expectedFilePath string
		shouldError      bool
	}{
		{
			inputFilePath:    "testdata/success/1.md",
			expectedFilePath: "testdata/success/1_exp.yaml",
			shouldError:      false,
		},
		{
			inputFilePath:    "testdata/error/1.md",
			expectedFilePath: "",
			shouldError:      true,
		},
	}

	// Create a variable successful test cases
	type testCase struct {
		reader      *bufio.Reader
		shouldError bool
		expected    []byte
	}
	tests := []testCase{}

	// Open the test files and create readers
	for _, f := range filePaths {
		// Open the test file
		testFileHandle, err := os.Open(f.inputFilePath)
		if err != nil {
			t.Fatalf("error opening test file %s", err)
		}
		defer testFileHandle.Close()

		var expectedBytes []byte
		if f.expectedFilePath != "" {
			// Open the expected file
			expectedFileHandle, err := os.Open(f.expectedFilePath)
			if err != nil {
				t.Fatalf("error opening expected file %s", err)
			}
			defer expectedFileHandle.Close()

			// Read in the expected file
			expectedBytes, err = io.ReadAll(expectedFileHandle)
		}

		// Add the data to the test case
		tests = append(tests, testCase{
			reader:      bufio.NewReader(testFileHandle),
			shouldError: f.shouldError,
			expected:    expectedBytes,
		})
	}

	// Loop through each test case
	for i, tc := range tests {
		t.Run(fmt.Sprintf("new frontmatter error %t %d", tc.shouldError, i), func(t *testing.T) {
			// Call the New function
			fm, err := frontmatter.New(tc.reader)

			if tc.shouldError {
				// Confirm we received an error
				if err == nil {
					t.Errorf("expected an error but got none")
				}
			} else {
				// Confirm there is no error
				if err != nil {
					t.Errorf("expected no error, got %s", err)
				}

				// Confirm the values match as byte arrays
				if !cmp.Equal([]byte(fm), tc.expected) {
					t.Errorf("values do not match\nhave: %v\nwant: %v\n", fm, tc.expected)
				}
			}
		})
	}
}

func TestFrontmatter_Map(t *testing.T) {
	// Define successful test cases
	testsSuccess := []struct {
		input    frontmatter.Frontmatter
		expected map[string]any
	}{
		{
			input: []byte("one: value 1\ntwo: value 2"),
			expected: map[string]any{
				"one": "value 1",
				"two": "value 2",
			},
		},
	}

	// Loop through test cases
	for i, tc := range testsSuccess {
		t.Run(fmt.Sprintf("map success %d", i), func(t *testing.T) {
			// Call the Map function
			m, err := tc.input.Map()

			// Confirm there is no error
			if err != nil {
				t.Errorf("expected no error, got %s", err)
			}

			// Confirm the values match
			if diff := cmp.Diff(m, tc.expected); diff != "" {
				t.Errorf("values do not match (-want +got):\n%s", diff)
			}
		})
	}

	// Define error test cases
	testsError := []frontmatter.Frontmatter{
		[]byte("one: two:"),
	}

	// Loop through each test case
	for i, tc := range testsError {
		t.Run(fmt.Sprintf("map error %d", i), func(t *testing.T) {
			// Call the Map function
			_, err := tc.Map()

			// Confirm we received an error
			if err == nil {
				t.Errorf("expected an error but got none")
			}
		})
	}
}

func TestFrontmatter_MapWithoutEmpty(t *testing.T) {
	// Define successful test cases
	testsSuccess := []struct {
		input    frontmatter.Frontmatter
		expected map[string]any
	}{
		{
			input: []byte("one:\ntwo: value 2"),
			expected: map[string]any{
				"two": "value 2",
			},
		},
	}

	// Loop through test cases
	for i, tc := range testsSuccess {
		t.Run(fmt.Sprintf("map success %d", i), func(t *testing.T) {
			// Call the Map function
			m, err := tc.input.MapWithoutEmpty()

			// Confirm there is no error
			if err != nil {
				t.Errorf("expected no error, got %s", err)
			}

			// Confirm the values match
			if diff := cmp.Diff(m, tc.expected); diff != "" {
				t.Errorf("values do not match (-want +got):\n%s", diff)
			}
		})
	}

	// Define error test cases
	testsError := []frontmatter.Frontmatter{
		[]byte("one: two:"),
	}

	// Loop through each test case
	for i, tc := range testsError {
		t.Run(fmt.Sprintf("map error %d", i), func(t *testing.T) {
			// Call the Map function
			_, err := tc.MapWithoutEmpty()

			// Confirm we received an error
			if err == nil {
				t.Errorf("expected an error but got none")
			}
		})
	}
}

func TestFrontmatter_String(t *testing.T) {
	// Define test cases
	tests := []struct {
		input    frontmatter.Frontmatter
		expected string
	}{
		{
			input:    []byte("one: value 1\ntwo: value 2\n"),
			expected: "one: value 1\ntwo: value 2\n",
		},
	}

	// Loop through test cases
	for i, tc := range tests {
		t.Run(fmt.Sprintf("string %d", i), func(t *testing.T) {
			// Call String function
			output := tc.input.String()

			// Comfirm the values match
			if output != tc.expected {
				t.Errorf("values do not match, have: %s want: %s", output, tc.expected)
			}
		})
	}
}

func TestFrontMatter_WrappedWithSeparatorString(t *testing.T) {
	// Define test cases
	tests := []struct {
		input    frontmatter.Frontmatter
		expected string
	}{
		{
			input:    []byte("one: value 1\ntwo: value 2\n"),
			expected: "---\none: value 1\ntwo: value 2\n---\n",
		},
	}

	// Loop through each test case
	for i, tc := range tests {
		t.Run(fmt.Sprintf("wrapped with separator string %d", i), func(t *testing.T) {
			// Call the WrappedWithSeparatorString function
			output := tc.input.WrappedWithSeparatorString()

			// Confirm the values match
			if output != tc.expected {
				t.Errorf("values do not match, have: %s want: %s", output, tc.expected)
			}
		})
	}
}

func TestFrontMatter_WrappedWithSeparatorBytes(t *testing.T) {
	// Define test cases
	tests := []struct {
		input    frontmatter.Frontmatter
		expected []byte
	}{
		{
			input:    []byte("one: value 1\ntwo: value 2\n"),
			expected: []byte("---\none: value 1\ntwo: value 2\n---\n"),
		},
	}

	// Loop through each test case
	for i, tc := range tests {
		t.Run(fmt.Sprintf("wrapped with separator bytes %d", i), func(t *testing.T) {
			// Call the WrappedWithSeparatorBytes function
			output := tc.input.WrappedWithSeparatorBytes()

			// Confirm the values match
			if !cmp.Equal(output, tc.expected) {
				t.Errorf("values do not match, have: %s want: %s", output, tc.expected)
			}
		})
	}
}

func TestFrontmatter_Merge(t *testing.T) {
	// Define successful test cases
	testsSuccess := []struct {
		base     []byte
		patch    []byte
		expected frontmatter.Frontmatter
	}{
		{
			base:     []byte("one: value 1"),
			patch:    []byte("two: value 2"),
			expected: []byte("one: value 1\ntwo: value 2\n"),
		},
	}

	// Loop through each test case
	for i, tc := range testsSuccess {
		t.Run(fmt.Sprintf("merge success %d", i), func(t *testing.T) {
			// Call the Merge function
			fm, err := frontmatter.Merge(tc.base, tc.patch)

			// Confirm there is no error
			if err != nil {
				t.Errorf("no error expected, got %s", err)
			}

			// Confirm the values match
			if !cmp.Equal(fm, tc.expected) {
				t.Errorf("values do not match, have: %s want: %s", fm, tc.expected)
			}
		})
	}

	// Define error test cases
	testsError := []struct {
		base  []byte
		patch []byte
	}{
		{
			base:  []byte("one: two:"),
			patch: []byte("three: value 1"),
		},
		{
			base:  []byte("three: value 1"),
			patch: []byte("one: two:"),
		},
	}

	// Loop through each test case
	for i, tc := range testsError {
		t.Run(fmt.Sprintf("merge error %d", i), func(t *testing.T) {
			// Call the Merge function
			_, err := frontmatter.Merge(tc.base, tc.patch)

			// Confirm we received an error
			if err == nil {
				t.Errorf("expected an error but got none")
			}
		})
	}
}
