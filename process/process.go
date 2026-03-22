package process

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"text/template"

	"github.com/watercorp/fmr/filedetails"
	"github.com/watercorp/fmr/frontmatter"
	"github.com/watercorp/fmr/md"
	"github.com/watercorp/fmr/templatefuncs"
)

// Set the writer for validation
var ValidateWriter io.Writer = os.Stdout

func ReplaceSourceFromTemplate(sourceFilePath string, templateFilePath string, retainTaskListItems bool) error {
	// Generate file details for the source and template paths to simpilify processing
	sourceFd := filedetails.New(sourceFilePath)
	templateFd := filedetails.New(templateFilePath)

	// Open source
	sourceFile, err := os.Open(sourceFd.FullPath)
	if err != nil {
		return fmt.Errorf("cannot open source file: %w", err)
	}
	defer sourceFile.Close()

	// Open Template
	templateFile, err := os.Open(templateFd.FullPath)
	if err != nil {
		return fmt.Errorf("cannot open template file: %w", err)
	}
	defer templateFile.Close()

	// Create a temp file for the output
	outputTempFile, err := os.CreateTemp(sourceFd.Directory, fmt.Sprintf("*-%s.tmp", sourceFd.Name))
	if err != nil {
		return fmt.Errorf("cannot create temporary output file: %w", err)
	}
	defer outputTempFile.Close()
	defer os.Remove(outputTempFile.Name())

	// Generate file details for the temporary output file
	outputTempFd := filedetails.New(outputTempFile.Name())

	// Create a buffer for the template that we're processing
	newTemplateBuffer := new(bytes.Buffer)

	// Create Readers and Writers
	sourceReader := bufio.NewReader(sourceFile)
	templateReader := bufio.NewReader(templateFile)
	outputTempWriter := bufio.NewWriter(outputTempFile)
	newTemplateWriter := bufio.NewWriter(newTemplateBuffer)

	// Gather the frontmater from both the source and the template
	sourceFrontmatter, err := frontmatter.New(sourceReader)
	if err != nil {
		return err
	}
	templateFrontmatter, err := frontmatter.New(templateReader)
	if err != nil {
		return err
	}

	// Merge the frontmatter
	mergedFrontmatter, err := frontmatter.Merge(templateFrontmatter, sourceFrontmatter)
	if err != nil {
		return err
	}

	// Convert the frontmatter to a map for processing
	mergedMap, err := mergedFrontmatter.Map()
	if err != nil {
		return err
	}

	// Write the merged frontmatter to the temporary output file
	newTemplateBuffer.Write(mergedFrontmatter.WrappedWithSeparatorBytes())

	// Only handle checked task items if we need to
	if retainTaskListItems {

		// Continue reading the source file to record the checked task list items
		checkedItems, err := recordTaskItems(sourceReader, true, true)
		if err != nil {
			return err
		}

		// Continue reading the template and replace the checks where necessary inside the output file
		replaceTaskListItems(checkedItems, templateReader, newTemplateWriter)
	} else {
		// Read in the rest of the template to the buffer
		for {
			line, _, err := templateReader.ReadLine()
			if err != nil && err != io.EOF {
				return err
			}
			if err == io.EOF {
				break
			}
			fmt.Fprintf(newTemplateWriter, "%s\n", line)
		}
	}

	// Flush the writer to the buffer
	newTemplateWriter.Flush()

	// Parse the template
	tmpl, err := template.New("main").Funcs(templatefuncs.FuncMap).Parse(newTemplateBuffer.String())
	if err != nil {
		return err
	}

	// Execute the template
	tmpl.Execute(outputTempWriter, mergedMap)

	// Flush the output temp writer
	outputTempWriter.Flush()

	// Ensure the source file is closed
	sourceFile.Close()

	// Rename the temporary output to the source file
	err = os.Rename(outputTempFd.FullPath, sourceFd.FullPath)
	if err != nil {
		return err
	}

	return nil
}

func ReplaceSourceOnly(sourceFilePath string) error {
	// Generate file details for the source path to simpilify processing
	sourceFd := filedetails.New(sourceFilePath)

	// Open source
	sourceFile, err := os.Open(sourceFd.FullPath)
	if err != nil {
		return fmt.Errorf("cannot open source file: %w", err)
	}
	defer sourceFile.Close()

	// Create a temp file for the output
	outputTempFile, err := os.CreateTemp(sourceFd.Directory, fmt.Sprintf("*-%s.tmp", sourceFd.Name))
	if err != nil {
		return fmt.Errorf("cannot create temporary output file: %w", err)
	}
	defer outputTempFile.Close()
	defer os.Remove(outputTempFile.Name())

	// Generate file details for the temporary output file
	outputTempFd := filedetails.New(outputTempFile.Name())

	// Create Readers and Writers
	sourceReader := bufio.NewReader(sourceFile)
	outputTempWriter := bufio.NewWriter(outputTempFile)

	// Gather the frontmater from the source
	sourceFrontmatter, err := frontmatter.New(sourceReader)
	if err != nil {
		return err
	}

	// Convert the frontmatter to a map for processing
	sourceFrontmatterMap, err := sourceFrontmatter.Map()
	if err != nil {
		return err
	}

	// Read the whole file
	readSource, err := io.ReadAll(sourceReader)
	if err != nil {
		return err
	}

	// Combine the source frontmatter and the rest of the read source
	fullSource := append(sourceFrontmatter.WrappedWithSeparatorBytes(), readSource...)

	// Parse the template
	tmpl, err := template.New("main").Funcs(templatefuncs.FuncMap).Parse(string(fullSource))
	if err != nil {
		return err
	}

	// Execute the template
	tmpl.Execute(outputTempWriter, sourceFrontmatterMap)

	// Flush the output temp writer
	outputTempWriter.Flush()

	// Ensure the source file is closed
	sourceFile.Close()

	// Rename the temporary output to the source file
	err = os.Rename(outputTempFd.FullPath, sourceFd.FullPath)
	if err != nil {
		return err
	}

	return nil
}

func ReplaceOther(sourceFilePath string, destinationFilePath string, templateFilePath *string) error {
	// Generate file details for the source and template paths to simpilify processing
	sourceFd := filedetails.New(sourceFilePath)
	destinationFd := filedetails.New(destinationFilePath)

	// Open source
	sourceFile, err := os.Open(sourceFd.FullPath)
	if err != nil {
		return fmt.Errorf("cannot open source file: %w", err)
	}
	defer sourceFile.Close()

	// Create a temp file for the destination
	destinationTempFile, err := os.CreateTemp(destinationFd.Directory, fmt.Sprintf("*-%s.tmp", destinationFd.Name))
	if err != nil {
		return fmt.Errorf("cannot create temporary destination file: %w", err)
	}
	defer destinationTempFile.Close()
	defer os.Remove(destinationTempFile.Name())

	// Generate file details for the temporary destination file
	destinationTempFd := filedetails.New(destinationTempFile.Name())

	// Create Readers and Writers
	sourceReader := bufio.NewReader(sourceFile)
	destinationTempWriter := bufio.NewWriter(destinationTempFile)

	// Gather the frontmater from the source
	sourceFrontmatter, err := frontmatter.New(sourceReader)
	if err != nil {
		return err
	}

	// Convert the frontmatter to a map for processing
	sourceFrontmatterMap, err := sourceFrontmatter.Map()
	if err != nil {
		return err
	}

	// Determine if we are using a different template
	var finalTemplateFd filedetails.FileDetails
	if templateFilePath != nil {
		finalTemplateFd = filedetails.New(*templateFilePath)
	} else {
		finalTemplateFd = destinationFd
	}

	// Create the template
	tmpl := template.New(finalTemplateFd.Name).Funcs(templatefuncs.FuncMap)

	// Check if we need to adjust the delimiters
	if slices.Contains([]string{".json", ".jsonc"}, destinationFd.Extension) {
		// Switch the delimiters to ";;" for JSON files
		tmpl = tmpl.Delims("<<", ">>")
	}

	// Parse the template
	tmpl, err = tmpl.ParseFiles(finalTemplateFd.FullPath)
	if err != nil {
		return err
	}

	// Execute the template
	tmpl.Execute(destinationTempWriter, sourceFrontmatterMap)

	// Flush the destination temp writer
	destinationTempWriter.Flush()

	// Rename the temporary destination to the destination file
	err = os.Rename(destinationTempFd.FullPath, destinationFd.FullPath)
	if err != nil {
		return err
	}

	return nil
}

func ValidateTaskListItems(sourceFilePath string, templateFilePath string) error {
	// Initialize variables
	var sourceTaskListItems, templateTaskListItems []md.TaskListItem
	var sourceUniqueTaskListItems = map[string]map[string]bool{}
	var templateUniqueTaskListItems = map[string]map[string]bool{}
	var errorLists = struct {
		TemplateNonUnique          []md.TaskListItem
		SourceNonUnique            []md.TaskListItem
		TemplateMissingCheckedTask []md.TaskListItem
	}{
		TemplateNonUnique:          []md.TaskListItem{},
		SourceNonUnique:            []md.TaskListItem{},
		TemplateMissingCheckedTask: []md.TaskListItem{},
	}

	// Generate file details for the source and template paths to simpilify processing
	sourceFd := filedetails.New(sourceFilePath)
	templateFd := filedetails.New(templateFilePath)

	// Open source
	sourceFile, err := os.Open(sourceFd.FullPath)
	if err != nil {
		return fmt.Errorf("cannot open source file: %w", err)
	}
	defer sourceFile.Close()

	// Open Template
	templateFile, err := os.Open(templateFd.FullPath)
	if err != nil {
		return fmt.Errorf("cannot open template file: %w", err)
	}
	defer templateFile.Close()

	// Create Readers
	sourceReader := bufio.NewReader(sourceFile)
	templateReader := bufio.NewReader(templateFile)

	// Search the template line by line
	for {
		line, _, err := templateReader.ReadLine()
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			break
		}

		// Match the task list item
		match, err := md.MatchTaskListItem(string(line))

		// Only proceed if we have a match
		if err == nil && match != nil {
			// Add the task to the source list items
			templateTaskListItems = append(templateTaskListItems, *match)

			// Check if our map contains a key for the whitespace
			if w, ok := templateUniqueTaskListItems[match.Whitespace]; ok {
				// Check if the task text is already in the map
				if _, ok := w[match.Text]; ok {
					//fmt.Printf("Template - Non-Unique Task Item Found:\n%s\n", match)
					errorLists.TemplateNonUnique = append(errorLists.TemplateNonUnique, *match)
				} else {
					// Add the task to the map
					w[match.Text] = true
				}
			} else {
				// Add the whitespace entry to the map
				templateUniqueTaskListItems[match.Whitespace] = map[string]bool{
					match.Text: true,
				}
			}
		}
	}

	// Search the source line by line
	for {
		line, _, err := sourceReader.ReadLine()
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			break
		}

		// Match the task list item
		match, err := md.MatchTaskListItem(string(line))

		// Only proceed if we have a match
		if err == nil && match != nil {
			// Add the task to the source list items
			sourceTaskListItems = append(sourceTaskListItems, *match)

			// If the item is checked, confirm we have a matching version in the template
			if match.Checked {
				found := slices.ContainsFunc(templateTaskListItems, func(ti md.TaskListItem) bool {
					return ti.Equal(*match, false)
				})

				if !found {
					//fmt.Printf("Template - Missing Task for Checked Item:\n%s\n", match)
					errorLists.TemplateMissingCheckedTask = append(errorLists.TemplateMissingCheckedTask, *match)
				}
			}

			// Check if our map contains a key for the whitespace
			if w, ok := sourceUniqueTaskListItems[match.Whitespace]; ok {
				// Check if the task text is already in the map
				if _, ok := w[match.Text]; ok {
					//fmt.Printf("Source - Non-Unique Task Item Found:\n%s\n", match)
					errorLists.SourceNonUnique = append(errorLists.SourceNonUnique, *match)
				} else {
					// Add the task to the map
					w[match.Text] = true
				}
			} else {
				// Add the whitespace entry to the map
				sourceUniqueTaskListItems[match.Whitespace] = map[string]bool{
					match.Text: true,
				}
			}
		}
	}

	// Check if we have any found issues
	if len(errorLists.SourceNonUnique) > 0 || len(errorLists.TemplateNonUnique) > 0 || len(errorLists.TemplateMissingCheckedTask) > 0 {
		if len(errorLists.TemplateNonUnique) > 0 {
			fmt.Fprintln(ValidateWriter, "Template - Non-Unique Task Items Found:")
			for _, v := range errorLists.TemplateNonUnique {
				fmt.Fprintln(ValidateWriter, v)
			}
			fmt.Fprintln(ValidateWriter)
		}

		if len(errorLists.SourceNonUnique) > 0 {
			fmt.Fprintln(ValidateWriter, "Source - Non-Unique Task Items Found:")
			for _, v := range errorLists.SourceNonUnique {
				fmt.Fprintln(ValidateWriter, v)
			}
			fmt.Fprintln(ValidateWriter)
		}

		if len(errorLists.TemplateMissingCheckedTask) > 0 {
			fmt.Fprintln(ValidateWriter, "Template - Missing Tasks for Checked Item:")
			for _, v := range errorLists.TemplateMissingCheckedTask {
				fmt.Fprintln(ValidateWriter, v)
			}
		}
	} else {
		fmt.Fprintln(ValidateWriter, "No issues found with Source or Template!")
	}

	return nil
}

func recordTaskItems(reader *bufio.Reader, matchCheckedState bool, checked bool) ([]md.TaskListItem, error) {
	// Create a variable to hold the found task list items
	taskListItems := []md.TaskListItem{}

	// Search by each line
	for {
		line, _, err := reader.ReadLine()
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF {
			break
		}

		// Match the task list item
		match, err := md.MatchTaskListItem(string(line))

		// Only add the item if it equals the checked flag
		if err == nil {
			if matchCheckedState && match.Checked == checked {
				taskListItems = append(taskListItems, *match)
			} else if !matchCheckedState {
				taskListItems = append(taskListItems, *match)
			}
		}
	}

	return taskListItems, nil
}

func replaceTaskListItems(items []md.TaskListItem, reader *bufio.Reader, writer *bufio.Writer) error {
	// Initialize the check index
	taskListItemIndex := 0

	// Regex for matching template items
	templateMatch := regexp.MustCompile(`\\{\\{.+?\\}\\}`)

	// Search line by line
	for {
		line, _, err := reader.ReadLine()
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			break
		}

		// Check if we have matched an unchecked item
		match, err := md.MatchTaskListItem(string(line))
		if err == nil && match.Checked == false {
			// Confirm we have more recorded checks we can process
			if taskListItemIndex < len(items) {
				// If we come across a template item, replace it with some generic regex
				matchText := templateMatch.ReplaceAllString(regexp.QuoteMeta(match.Text), ".+?")

				// Check if we found a full match against the text
				checkMatchText, err := regexp.MatchString(fmt.Sprintf("^%s$", matchText), items[taskListItemIndex].Text)
				if err != nil {
					return err
				}

				// Confirm both the Whitespace is equal and we had a text match
				if match.Whitespace == items[taskListItemIndex].Whitespace && checkMatchText {
					// Write the line with the recorded check
					fmt.Fprintf(writer, "%s\n", items[taskListItemIndex])
					taskListItemIndex++
				} else {
					// If we did not have a match, write the line as it was
					fmt.Fprintf(writer, "%s\n", line)
				}
			} else {
				// If there are not more checks to process, write the line as it was
				fmt.Fprintf(writer, "%s\n", line)
			}
		} else {
			// If this isn't a task list item, write the line as it was
			fmt.Fprintf(writer, "%s\n", line)
		}
	}

	return nil
}
