package process

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/watercorp/fmr/md"
)

func TestProcess_internal_recordTaskItems(t *testing.T) {
	// Open a valid test file
	testFile, err := os.Open("testdata/valid.md")
	if err != nil {
		t.Fatalf("failed to open test file: %s", err)
	}
	defer testFile.Close()

	// Create a reader
	reader := bufio.NewReader(testFile)

	// Define the test cases
	tests := []struct {
		matchCheckedState bool
		checked           bool
		expected          []md.TaskListItem
	}{
		{
			matchCheckedState: true,
			checked:           true,
			expected: []md.TaskListItem{
				{
					Whitespace: "",
					Text:       "Step 1",
					Checked:    true,
				},
				{
					Whitespace: "",
					Text:       "Step 2",
					Checked:    true,
				},
				{
					Whitespace: "  ",
					Text:       "Sub Step 1",
					Checked:    true,
				},
				{
					Whitespace: "",
					Text:       "Step 3 value 1",
					Checked:    true,
				},
			},
		},
		{
			matchCheckedState: true,
			checked:           false,
			expected: []md.TaskListItem{
				{
					Whitespace: "  ",
					Text:       "Sub Step 2",
					Checked:    false,
				},
				{
					Whitespace: "",
					Text:       "Step 3",
					Checked:    false,
				},
				{
					Whitespace: "",
					Text:       "Step 4",
					Checked:    false,
				},
			},
		},
		{
			matchCheckedState: false,
			checked:           true,
			expected: []md.TaskListItem{
				{
					Whitespace: "",
					Text:       "Step 1",
					Checked:    true,
				},
				{
					Whitespace: "",
					Text:       "Step 2",
					Checked:    true,
				},
				{
					Whitespace: "  ",
					Text:       "Sub Step 1",
					Checked:    true,
				},
				{
					Whitespace: "  ",
					Text:       "Sub Step 2",
					Checked:    false,
				},
				{
					Whitespace: "",
					Text:       "Step 3 value 1",
					Checked:    true,
				},
				{
					Whitespace: "",
					Text:       "Step 3",
					Checked:    false,
				},
				{
					Whitespace: "",
					Text:       "Step 4",
					Checked:    false,
				},
			},
		},
	}

	// Loop through each test case
	for i, tc := range tests {
		t.Run(fmt.Sprintf("record task items %d", i), func(t *testing.T) {
			// Call the recordTaskItems function
			ti, err := recordTaskItems(reader, tc.matchCheckedState, tc.checked)

			// Confirm there are no errors
			if err != nil {
				t.Errorf("expected no error, got %s", err)
			}

			// Confirm the values match
			if diff := cmp.Diff(ti, tc.expected); diff != "" {
				t.Errorf("values do not match (-want +got):\n%s", diff)
			}
		})

		// Reset the reader for the next test
		_, err := testFile.Seek(0, io.SeekStart)
		if err != nil {
			t.Fatal("unable to reset reader")
		}
	}
}

func TestProcess_internal_replaceTaskListItems(t *testing.T) {
	// Define task list items
	taskListItems := []md.TaskListItem{
		{
			Whitespace: "",
			Text:       "Step 1",
			Checked:    true,
		},
		{
			Whitespace: "",
			Text:       "Step 2",
			Checked:    true,
		},
		{
			Whitespace: "  ",
			Text:       "Sub Step 1",
			Checked:    true,
		},
		{
			Whitespace: "",
			Text:       "Step 5 value 1",
			Checked:    true,
		},
	}

	// Open a valid template file
	templateFile, err := os.Open("testdata/task-list-items-template.md")
	if err != nil {
		t.Fatalf("failed to open template file: %s", err)
	}
	defer templateFile.Close()

	// Create a reader for the template
	templateReader := bufio.NewReader(templateFile)

	// Open a valid test file
	testFile, err := os.Open("testdata/task-list-items.md")
	if err != nil {
		t.Fatalf("failed to open test file: %s", err)
	}
	defer testFile.Close()

	// Create a reader
	reader := bufio.NewReader(testFile)

	// Read in the entire test file for comparison
	testFileContents, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal("failed to read test file contents")
	}

	// Create a bytes buffer
	destination := bytes.Buffer{}
	writer := bufio.NewWriter(&destination)

	// Run the test
	t.Run("replace task list items", func(t *testing.T) {
		// Call the replaceTaskListItems functions
		err := replaceTaskListItems(taskListItems, templateReader, writer)

		// Confirm we have no error
		if err != nil {
			t.Errorf("expected no error, got %s", err)
		}

		// Flush the buffer to the destination
		writer.Flush()

		// Confirm the contents of the test file match the writer
		if !cmp.Equal(destination.Bytes(), testFileContents) {
			t.Errorf("values do not match, have:\n%s\nwant:\n%s", destination.Bytes(), testFileContents)
		}
	})

}
