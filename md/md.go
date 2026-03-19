package md

import (
	"fmt"
	"regexp"
)

// Types
type TaskListItem struct {
	Whitespace string
	Text       string
	Checked    bool
}

var templateMatch = regexp.MustCompile(`\\{\\{.+?\\}\\}`)

func NewTaskListItem(whitespace string, text string, checked bool) *TaskListItem {
	return &TaskListItem{
		Whitespace: whitespace,
		Text:       text,
		Checked:    checked,
	}
}

func (ti TaskListItem) String() string {
	// Determine what value to use as a checkmark
	check := " "
	if ti.Checked {
		check = "X"
	}

	// Return the string
	return fmt.Sprintf("%s- [%s] %s", ti.Whitespace, check, ti.Text)
}

func (ti TaskListItem) Equal(ti2 TaskListItem, matchChecked bool) bool {
	// If we come across a template item, replace it with some generic regex
	matchText := templateMatch.ReplaceAllString(regexp.QuoteMeta(ti.Text), ".+?")

	// Check if we found a full match against the text
	checkMatchText, err := regexp.MatchString(fmt.Sprintf("^%s$", matchText), ti2.Text)
	if err != nil {
		return false
	}

	// Confirm both the Whitespace is equal and we had a text match
	if matchChecked {
		if ti.Whitespace == ti2.Whitespace && checkMatchText && ti.Checked == ti2.Checked {
			return true
		}
	} else {
		if ti.Whitespace == ti2.Whitespace && checkMatchText {
			return true
		}
	}

	return false

}

// Define task list item match
var taskListItemMatch = regexp.MustCompile(`(\s*)- \[([ |X])\] (.+)`)

func MatchTaskListItem(line string) (*TaskListItem, error) {
	// Initialize the checked varible
	var checked bool

	found := taskListItemMatch.FindStringSubmatch(line)

	// Confirm we have the correct amount of captured items
	if len(found) != 4 {
		return nil, fmt.Errorf("not matched")
	}

	// Determine if the item is checked
	if found[2] == "X" {
		checked = true
	}

	return &TaskListItem{
		Whitespace: found[1],
		Text:       found[3],
		Checked:    checked,
	}, nil
}
