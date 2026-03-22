package process_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/watercorp/fmr/process"
)

func TestProcess_ReplaceSourceFromTemplate(t *testing.T) {
	// Create a new temp file that will become the source
	tempSourceFile, err := os.CreateTemp("testdata", "valid-temp-*.md")
	if err != nil {
		t.Fatal("failed to create temp source file")
	}
	defer os.Remove(tempSourceFile.Name())
	defer tempSourceFile.Close()

	// Open the source file
	sourceFile, err := os.Open("testdata/valid.md")
	if err != nil {
		t.Fatal("failed to open source file")
	}
	defer sourceFile.Close()

	// Read in the entire source file for later comparison
	sourceFileContentsBuffer := bytes.Buffer{}

	// Create a tee to read the source file into a buffer and to the temp file
	tee := io.TeeReader(sourceFile, io.MultiWriter(tempSourceFile, &sourceFileContentsBuffer))

	// Copy the source file to both the temp file and buffer
	_, err = io.Copy(io.Discard, tee)
	if err != nil {
		t.Fatal("failed to copy source file tee")
	}

	// Ensure the data is synced
	err = tempSourceFile.Sync()
	if err != nil {
		t.Fatal("failed to sync temp source file")
	}

	// Seek back to the beginning so we can read from the temp file
	_, err = tempSourceFile.Seek(0, io.SeekStart)
	if err != nil {
		t.Fatal("unable to reset offset in temp file")
	}

	// Run the successful test
	t.Run("replace source from template success", func(t *testing.T) {
		// Call the ReplaceSourceFromTemplate function
		err := process.ReplaceSourceFromTemplate(tempSourceFile.Name(), "testdata/valid-template.md", true)
		if err != nil {
			t.Errorf("no error expected, got %s", err)
		}

		// Read the contents of the temp source file
		tempSourceContents, err := io.ReadAll(tempSourceFile)
		if err != nil {
			t.Fatal("failed to read temp source file")
		}

		// Confirm the temp source and the source match
		if !cmp.Equal(tempSourceContents, sourceFileContentsBuffer.Bytes()) {
			t.Errorf("files do not match, have:\n%s\nwant:\n%s", tempSourceContents, sourceFileContentsBuffer.Bytes())
		}
	})

	// Run tests with bad paths for error checking
	testsError := []struct {
		sourceFilePath   string
		templateFilePath string
	}{
		{
			sourceFilePath:   "testdata/invalidsourcepath.md",
			templateFilePath: "testdata/valid-template.md",
		},
		{
			sourceFilePath:   "testdata/valid.md",
			templateFilePath: "testdata/invalidtemplatepath.md",
		},
		{
			sourceFilePath:   "testdata/invalid-frontmatter.md",
			templateFilePath: "testdata/valid-template.md",
		},
		{
			sourceFilePath:   "testdata/valid.md",
			templateFilePath: "testdata/invalid-frontmatter.md",
		},
	}

	// Loop through each test case
	for i, tc := range testsError {
		t.Run(fmt.Sprintf("replace source from template error %d", i), func(t *testing.T) {
			// Call the ReplaceSourceFromTemplate function
			err := process.ReplaceSourceFromTemplate(tc.sourceFilePath, tc.templateFilePath, true)

			// Confirm we received an error
			if err == nil {
				t.Errorf("expected an error but got none")
			}
		})
	}
}

func TestProcess_ReplaceSourceFromTemplate_NoChecks(t *testing.T) {
	// Create a new temp file that will become the source
	tempSourceFile, err := os.CreateTemp("testdata", "valid-no-checks-temp-*.md")
	if err != nil {
		t.Fatal("failed to create temp source file")
	}
	defer os.Remove(tempSourceFile.Name())
	defer tempSourceFile.Close()

	// Open the source file
	sourceFile, err := os.Open("testdata/valid-no-checks.md")
	if err != nil {
		t.Fatal("failed to open source file")
	}
	defer sourceFile.Close()

	// Read in the entire source file for later comparison
	sourceFileContentsBuffer := bytes.Buffer{}

	// Create a tee to read the source file into a buffer and to the temp file
	tee := io.TeeReader(sourceFile, io.MultiWriter(tempSourceFile, &sourceFileContentsBuffer))

	// Copy the source file to both the temp file and buffer
	_, err = io.Copy(io.Discard, tee)
	if err != nil {
		t.Fatal("failed to copy source file tee")
	}

	// Ensure the data is synced
	err = tempSourceFile.Sync()
	if err != nil {
		t.Fatal("failed to sync temp source file")
	}

	// Seek back to the beginning so we can read from the temp file
	_, err = tempSourceFile.Seek(0, io.SeekStart)
	if err != nil {
		t.Fatal("unable to reset offset in temp file")
	}

	// Run the test
	t.Run("replace source from template", func(t *testing.T) {
		// Call the ReplaceSourceFromTemplate function
		err := process.ReplaceSourceFromTemplate(tempSourceFile.Name(), "testdata/valid-template.md", false)
		if err != nil {
			t.Errorf("no error expected, got %s", err)
		}

		// Read the contents of the temp source file
		tempSourceContents, err := io.ReadAll(tempSourceFile)
		if err != nil {
			t.Fatal("failed to read temp source file")
		}

		// Confirm the temp source and the source match
		if !cmp.Equal(tempSourceContents, sourceFileContentsBuffer.Bytes()) {
			t.Errorf("files do not match, have:\n%s\nwant:\n%s", tempSourceContents, sourceFileContentsBuffer.Bytes())
		}
	})

}

func TestProcess_ReplaceSourceOnly(t *testing.T) {
	// Create a new temp file that will become the source
	tempSourceFile, err := os.CreateTemp("testdata", "source-replace-temp-*.md")
	if err != nil {
		t.Fatal("failed to create temp source file")
	}
	defer os.Remove(tempSourceFile.Name())
	defer tempSourceFile.Close()

	// Open the source file
	sourceFile, err := os.Open("testdata/source-replace.md")
	if err != nil {
		t.Fatal("failed to open source file")
	}
	defer sourceFile.Close()

	// Read in the entire source file for later comparison
	sourceFileContentsBuffer := bytes.Buffer{}

	// Create a tee to read the source file into a buffer and to the temp file
	tee := io.TeeReader(sourceFile, io.MultiWriter(tempSourceFile, &sourceFileContentsBuffer))

	// Copy the source file to both the temp file and buffer
	_, err = io.Copy(io.Discard, tee)
	if err != nil {
		t.Fatal("failed to copy source file tee")
	}

	// Ensure the data is synced
	err = tempSourceFile.Sync()
	if err != nil {
		t.Fatal("failed to sync temp source file")
	}

	// Seek back to the beginning so we can read from the temp file
	_, err = tempSourceFile.Seek(0, io.SeekStart)
	if err != nil {
		t.Fatal("unable to reset offset in temp file")
	}

	// Run the test
	t.Run("replace source only success", func(t *testing.T) {
		// Call the ReplaceSourceOnly function
		err := process.ReplaceSourceOnly(tempSourceFile.Name())
		if err != nil {
			t.Errorf("no error expected, got %s", err)
		}

		// Read the contents of the temp source file
		tempSourceContents, err := io.ReadAll(tempSourceFile)
		if err != nil {
			t.Fatal("failed to read temp source file")
		}

		// Confirm the temp source and the source match
		if !cmp.Equal(tempSourceContents, sourceFileContentsBuffer.Bytes()) {
			t.Errorf("files do not match, have:\n%s\nwant:\n%s", tempSourceContents, sourceFileContentsBuffer.Bytes())
		}
	})

	// Run tests with bad paths for error checking
	testsError := []string{
		"testdata/invalidsourcepath.md",
		"testdata/invalid-frontmatter.md",
	}

	// Loop through each test case
	for i, tc := range testsError {
		t.Run(fmt.Sprintf("replace source only error %d", i), func(t *testing.T) {
			// Call the ReplaceSourceOnly function
			err := process.ReplaceSourceOnly(tc)

			// Confirm we received an error
			if err == nil {
				t.Errorf("expected an error but got none")
			}
		})
	}
}

func TestProcess_ReplaceOther(t *testing.T) {
	// Create a new temp file that will become the destination
	tempDestinationFile, err := os.CreateTemp("testdata", "valid-temp-*.json")
	if err != nil {
		t.Fatal("failed to create temp destination file")
	}
	defer os.Remove(tempDestinationFile.Name())
	defer tempDestinationFile.Close()

	// Open the destination file
	destinationFile, err := os.Open("testdata/valid.json")
	if err != nil {
		t.Fatal("failed to open destination file")
	}
	defer destinationFile.Close()

	// Read in the entire destination file for later comparison
	destinationFileContentsBuffer := bytes.Buffer{}

	// Create a tee to read the destination file into a buffer and to the temp file
	tee := io.TeeReader(destinationFile, io.MultiWriter(tempDestinationFile, &destinationFileContentsBuffer))

	// Copy the destination file to both the temp file and buffer
	_, err = io.Copy(io.Discard, tee)
	if err != nil {
		t.Fatal("failed to copy destination file tee")
	}

	// Ensure the data is synced
	err = tempDestinationFile.Sync()
	if err != nil {
		t.Fatal("failed to sync temp destination file")
	}

	// Seek back to the beginning so we can read from the temp file
	_, err = tempDestinationFile.Seek(0, io.SeekStart)
	if err != nil {
		t.Fatal("unable to reset offset in temp file")
	}

	// Run the test
	t.Run("replace other success", func(t *testing.T) {
		// Create template path for pointer
		templatePath := "testdata/valid-template.json"

		// Call the ReplaceOther function
		err := process.ReplaceOther("testdata/valid.md", tempDestinationFile.Name(), &templatePath)
		if err != nil {
			t.Errorf("no error expected, got %s", err)
		}

		// Read the contents of the temp destination file
		tempDestinationContents, err := io.ReadAll(tempDestinationFile)
		if err != nil {
			t.Fatal("failed to read temp destination file")
		}

		// Confirm the temp destination and the destination match
		if !cmp.Equal(tempDestinationContents, destinationFileContentsBuffer.Bytes()) {
			t.Errorf("files do not match, have:\n%s\nwant:\n%s", tempDestinationContents, destinationFileContentsBuffer.Bytes())
		}
	})

	// Run tests with bad paths for error checking
	testsError := []struct {
		sourceFilePath      string
		destinationFilePath string
		templateFilePath    *string
	}{
		{
			sourceFilePath:      "testdata/invalidsourcepath.md",
			destinationFilePath: "testdata/valid-destination.json",
			templateFilePath:    new("testdata/valid-template.json"),
		},
		{
			sourceFilePath:      "testdata/valid.md",
			destinationFilePath: "testdata/valid-destination.json",
			templateFilePath:    new("testdata/invalidtemplatepath.json"),
		},
		{
			sourceFilePath:      "testdata/invalid-frontmatter.md",
			destinationFilePath: "testdata/valid-destination.json",
			templateFilePath:    new("testdata/valid-template.json"),
		},
		{
			sourceFilePath:      "testdata/valid.md",
			destinationFilePath: "testdata/invaliddir/invaliddestination.json",
			templateFilePath:    new("testdata/valid-template.json"),
		},
	}

	// Loop through each test case
	for i, tc := range testsError {
		t.Run(fmt.Sprintf("replace other error %d", i), func(t *testing.T) {
			// Call the ReplaceOther function
			err := process.ReplaceOther(tc.sourceFilePath, tc.destinationFilePath, tc.templateFilePath)

			// Confirm we received an error
			if err == nil {
				t.Errorf("expected an error but got none")
			}
		})
	}
}

func TestProcess_ReplaceOther_NoTemplate(t *testing.T) {
	// Create a new temp file that will become the destination
	tempDestinationFile, err := os.CreateTemp("testdata", "valid-temp-*.json")
	if err != nil {
		t.Fatal("failed to create temp destination file")
	}
	defer os.Remove(tempDestinationFile.Name())
	defer tempDestinationFile.Close()

	// Open the destination file
	destinationFile, err := os.Open("testdata/valid-template.json")
	if err != nil {
		t.Fatal("failed to open destination file")
	}
	defer destinationFile.Close()

	// // Read in the entire destination file for later comparison
	// destinationFileContentsBuffer := bytes.Buffer{}

	// // Create a tee to read the destination file into a buffer and to the temp file
	// tee := io.TeeReader(destinationFile, io.MultiWriter(tempDestinationFile, &destinationFileContentsBuffer))

	// Copy the destination file to the temp file
	_, err = io.Copy(tempDestinationFile, destinationFile)
	if err != nil {
		t.Fatal("failed to copy destination file tee")
	}

	// Open the final file
	finalFile, err := os.Open("testdata/valid.json")
	if err != nil {
		t.Fatal("failed to open final file")
	}
	defer finalFile.Close()

	// Read the contents of the final file for later comparison
	finalFileContents, err := io.ReadAll(finalFile)
	if err != nil {
		t.Fatal("unable to read contents of final file")
	}

	// // Ensure the data is synced
	// err = tempDestinationFile.Sync()
	// if err != nil {
	// 	t.Fatal("failed to sync temp destination file")
	// }

	// Seek back to the beginning so we can read from the temp file
	_, err = tempDestinationFile.Seek(0, io.SeekStart)
	if err != nil {
		t.Fatal("unable to reset offset in temp file")
	}

	// Run the test
	t.Run("replace other", func(t *testing.T) {
		// Call the ReplaceOther function
		err := process.ReplaceOther("testdata/valid.md", tempDestinationFile.Name(), nil)
		if err != nil {
			t.Errorf("no error expected, got %s", err)
		}

		// Reopen the temp destination file after the file has been renamed
		reopenTempDestinationFile, err := os.Open(tempDestinationFile.Name())
		if err != nil {
			t.Fatal("unable to reopen temp destination file")
		}
		defer reopenTempDestinationFile.Close()

		// Read the contents of the temp destination file
		tempDestinationContents, err := io.ReadAll(reopenTempDestinationFile)
		if err != nil {
			t.Fatal("failed to read temp destination file")
		}

		// Confirm the temp destination and the destination match
		if !cmp.Equal(tempDestinationContents, finalFileContents) {
			t.Errorf("files do not match, have:\n%s\nwant:\n%s", tempDestinationContents, finalFileContents)
		}
	})
}

func TestProcess_ValidateTaskListItems(t *testing.T) {
	// Create a new bytes buffer
	outputBuf := &bytes.Buffer{}

	// Change the validate writer
	process.ValidateWriter = outputBuf

	t.Run("validate task list items no issues", func(t *testing.T) {
		// Call ValidateTaskListItems
		err := process.ValidateTaskListItems("testdata/valid.md", "testdata/valid-template.md")

		// Confirm there is no error
		if err != nil {
			t.Errorf("no error expected, got %s", err)
		}

		// Set the expected message
		expectedMessage := "No issues found with Source or Template!\n"

		// Check the output and confirm it's what we expect
		if !cmp.Equal(outputBuf.Bytes(), []byte(expectedMessage)) {
			t.Errorf("values do not match, have:\n%s\nwant:\n %s", outputBuf.Bytes(), expectedMessage)
		}
	})

	// Reset the buffer
	outputBuf.Reset()

	t.Run("validate task list items issues found", func(t *testing.T) {
		// Call ValidateTaskListItems
		err := process.ValidateTaskListItems("testdata/validate-issues.md", "testdata/validate-issues-template.md")

		// Confirm there is no error
		if err != nil {
			t.Errorf("no error expected, got %s", err)
		}

		// Set the expected message
		expectedMessage := `Template - Non-Unique Task Items Found:
- [ ] Step 1

Source - Non-Unique Task Items Found:
- [ ] Step 2

Template - Missing Tasks for Checked Item:
- [X] Step 4
`

		// Check the output and confirm it's what we expect
		if !cmp.Equal(outputBuf.Bytes(), []byte(expectedMessage)) {
			t.Errorf("values do not match, have:\n%s\nwant:\n%s", outputBuf.Bytes(), expectedMessage)
		}
	})

	// Run tests with bad paths for error checking
	testsError := []struct {
		sourceFilePath   string
		templateFilePath string
	}{
		{
			sourceFilePath:   "testdata/invalidsourcepath.md",
			templateFilePath: "testdata/valid-template.md",
		},
		{
			sourceFilePath:   "testdata/valid.md",
			templateFilePath: "testdata/invalidtemplatepath.md",
		},
	}

	// Loop through each test case
	for i, tc := range testsError {
		t.Run(fmt.Sprintf("validate task list items error %d", i), func(t *testing.T) {
			// Call the ValidateTaskListItems function
			err := process.ValidateTaskListItems(tc.sourceFilePath, tc.templateFilePath)

			// Confirm we received an error
			if err == nil {
				t.Errorf("expected an error but got none")
			}
		})
	}
}
