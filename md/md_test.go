package md_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/watercorp/fmr/md"
)

func TestMd_NewTaskListItem(t *testing.T) {
	// Define test cases
	tests := []struct {
		whitespace string
		text       string
		checked    bool
		expected   *md.TaskListItem
	}{
		{
			whitespace: "",
			text:       "Step 1",
			checked:    false,
			expected: &md.TaskListItem{
				Whitespace: "",
				Text:       "Step 1",
				Checked:    false,
			},
		},
		{
			whitespace: "  ",
			text:       "Step 2",
			checked:    false,
			expected: &md.TaskListItem{
				Whitespace: "  ",
				Text:       "Step 2",
				Checked:    false,
			},
		},
		{
			whitespace: "",
			text:       "Step 1",
			checked:    true,
			expected: &md.TaskListItem{
				Whitespace: "",
				Text:       "Step 1",
				Checked:    true,
			},
		},
	}

	// Loop through each test case
	for i, tc := range tests {
		t.Run(fmt.Sprintf("new task list item %d", i), func(t *testing.T) {
			// Call the NewTaskListItem function
			ti := md.NewTaskListItem(tc.whitespace, tc.text, tc.checked)

			// Confirm the values match
			if diff := cmp.Diff(ti, tc.expected); diff != "" {
				t.Errorf("values do not match (-want -got):\n%s", diff)
			}
		})
	}
}

func TestMd_String(t *testing.T) {
	// Define test cases
	tests := []struct {
		input    md.TaskListItem
		expected string
	}{
		{
			input: md.TaskListItem{
				Whitespace: "",
				Text:       "Step 1",
				Checked:    false,
			},
			expected: "- [ ] Step 1",
		},
		{
			input: md.TaskListItem{
				Whitespace: "",
				Text:       "Step 1",
				Checked:    true,
			},
			expected: "- [X] Step 1",
		},
		{
			input: md.TaskListItem{
				Whitespace: "  ",
				Text:       "Step 1",
				Checked:    false,
			},
			expected: "  - [ ] Step 1",
		},
		{
			input: md.TaskListItem{
				Whitespace: "  ",
				Text:       "This is a new step",
				Checked:    true,
			},
			expected: "  - [X] This is a new step",
		},
	}

	// Loop through each test case
	for i, tc := range tests {
		t.Run(fmt.Sprintf("string %d", i), func(t *testing.T) {
			// Call the String function
			output := tc.input.String()

			// Confirm the values match
			if output != tc.expected {
				t.Errorf("values do not match, have: %s want: %s", output, tc.expected)
			}
		})
	}
}

func TestMd_Equal(t *testing.T) {
	// Define test cases
	tests := []struct {
		firstItem    md.TaskListItem
		secondItem   md.TaskListItem
		matchChecked bool
		expected     bool
	}{
		{
			firstItem: md.TaskListItem{
				Whitespace: "",
				Text:       "Step 1",
				Checked:    false,
			},
			secondItem: md.TaskListItem{
				Whitespace: "",
				Text:       "Step 1",
				Checked:    false,
			},
			matchChecked: true,
			expected:     true,
		},
		{
			firstItem: md.TaskListItem{
				Whitespace: "  ",
				Text:       "Step 1",
				Checked:    false,
			},
			secondItem: md.TaskListItem{
				Whitespace: "  ",
				Text:       "Step 1",
				Checked:    false,
			},
			matchChecked: true,
			expected:     true,
		},
		{
			firstItem: md.TaskListItem{
				Whitespace: "",
				Text:       "Step 1",
				Checked:    false,
			},
			secondItem: md.TaskListItem{
				Whitespace: "",
				Text:       "Step 1",
				Checked:    true,
			},
			matchChecked: false,
			expected:     true,
		},
		{
			firstItem: md.TaskListItem{
				Whitespace: "",
				Text:       "Step 1",
				Checked:    false,
			},
			secondItem: md.TaskListItem{
				Whitespace: "",
				Text:       "Step 1",
				Checked:    true,
			},
			matchChecked: true,
			expected:     false,
		},
		{
			firstItem: md.TaskListItem{
				Whitespace: "  ",
				Text:       "Step 1",
				Checked:    true,
			},
			secondItem: md.TaskListItem{
				Whitespace: "",
				Text:       "Step 1",
				Checked:    true,
			},
			matchChecked: true,
			expected:     false,
		},
		{
			firstItem: md.TaskListItem{
				Whitespace: "",
				Text:       "Step 1 {{.title}}",
				Checked:    true,
			},
			secondItem: md.TaskListItem{
				Whitespace: "",
				Text:       "Step 1 this is a title",
				Checked:    true,
			},
			matchChecked: true,
			expected:     true,
		},
		{
			firstItem: md.TaskListItem{
				Whitespace: "",
				Text:       "Step 1 {{.title}}",
				Checked:    true,
			},
			secondItem: md.TaskListItem{
				Whitespace: "",
				Text:       "Step 2",
				Checked:    true,
			},
			matchChecked: true,
			expected:     false,
		},
	}

	// Loop through test cases
	for i, tc := range tests {
		t.Run(fmt.Sprintf("equal %d", i), func(t *testing.T) {
			// Call the Equal function
			output := tc.firstItem.Equal(tc.secondItem, tc.matchChecked)

			// Confirm the values match
			if output != tc.expected {
				t.Errorf("values do not match, have: %t want: %t", output, tc.expected)
			}
		})
	}
}

func TestMd_MatchTaskListItem(t *testing.T) {
	// Define successful test cases
	testsSuccess := []struct {
		input    string
		expected *md.TaskListItem
	}{
		{
			input: "- [ ] Step 1",
			expected: &md.TaskListItem{
				Whitespace: "",
				Text:       "Step 1",
				Checked:    false,
			},
		},
		{
			input: "- [X] Step 1",
			expected: &md.TaskListItem{
				Whitespace: "",
				Text:       "Step 1",
				Checked:    true,
			},
		},
		{
			input: "  - [ ] Step 1",
			expected: &md.TaskListItem{
				Whitespace: "  ",
				Text:       "Step 1",
				Checked:    false,
			},
		},
	}

	// Loop through each test case
	for i, tc := range testsSuccess {
		t.Run(fmt.Sprintf("match task list item success %d", i), func(t *testing.T) {
			// Call the MatchTaskListItem function
			ti, err := md.MatchTaskListItem(tc.input)

			// Confirm there is no error
			if err != nil {
				t.Errorf("expected no error, got %s", err)
			}

			// Confirm the values match
			if diff := cmp.Diff(ti, tc.expected); diff != "" {
				t.Errorf("values do not match, (-want +got):\n%s", diff)
			}
		})
	}

	// Define error test cases
	testsError := []string{
		"Step 1",
		"- [] Step 1",
		"- Step 1",
		"",
	}

	// Loop through each test case
	for i, tc := range testsError {
		t.Run(fmt.Sprintf("match task list item error %d", i), func(t *testing.T) {
			// Call the MatchTaskListItem function
			_, err := md.MatchTaskListItem(tc)

			// Confirm we have an error
			if err == nil {
				t.Errorf("expected an error but got none")
			}
		})
	}
}
