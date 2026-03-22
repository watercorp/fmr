package frontmatter

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"
)

const SEP = "---"

type Frontmatter []byte

func New(reader *bufio.Reader) (Frontmatter, error) {
	// Create a variables to store the parsed frontmatter and buffer
	var (
		buf      strings.Builder
		fm       Frontmatter
		inside   bool
		captured bool
		err      error
	)

	// Loop through the file line by line
	for {
		// Read the line in
		line, _, err := reader.ReadLine()
		if err != nil && err != io.EOF {
			return nil, err
		}

		// Break if we've reached the end of the file somehow
		if err == io.EOF {
			err = nil
			break
		}

		// Check if we have the frontmatter block
		if strings.HasPrefix(string(line), SEP) && !captured {
			// If we're not inside yet, flag that we are
			if !inside {
				inside = true
			} else {
				// We have finished capturing the frontmatter, unmarshal the buffer
				inside = false
				captured = true
				fm = []byte(buf.String())

				// Leave the reader where it is for further file processing
				break
			}
		} else if inside {
			// Write the frontmatter to the buffer for later parsing
			fmt.Fprintf(&buf, "%s\n", line)
		}
	}

	// If the frontmatter is not captured, return an error
	if !captured {
		return nil, fmt.Errorf("unable to capture frontmatter")
	}

	return fm, err
}

// Return the frontmatter as a map
func (fm *Frontmatter) Map() (map[string]any, error) {
	var newMap map[string]any
	err := yaml.Unmarshal(*fm, &newMap)

	return newMap, err
}

// Return the frontmatter as a map without empty values
func (fm *Frontmatter) MapWithoutEmpty() (map[string]any, error) {
	// Convert the frontmatter to a map first
	newMap, err := fm.Map()
	if err != nil {
		return nil, err
	}

	// Clean the map
	return cleanMap(newMap), nil
}

func (fm Frontmatter) String() string {
	return string(fm)
}

func (fm Frontmatter) WrappedWithSeparatorString() string {
	return fmt.Sprintf("%s\n%s%s\n", SEP, fm.String(), SEP)
}
func (fm Frontmatter) WrappedWithSeparatorBytes() []byte {
	return fmt.Appendf([]byte{}, "%s\n%s%s\n", SEP, fm.String(), SEP)
}

func Merge(base []byte, patch []byte) (Frontmatter, error) {
	// Parse the bytes
	baseYaml, err := parser.ParseBytes(base, 0)
	if err != nil {
		return nil, fmt.Errorf("error parsing base yaml: %w", err)
	}
	patchYaml, err := parser.ParseBytes(patch, 0)
	if err != nil {
		return nil, fmt.Errorf("error parsing patch yaml: %w", err)
	}

	// Create a new path at the root
	rootPath, err := yaml.PathString("$")
	if err != nil {
		return nil, err
	}

	// Merge the yaml
	err = rootPath.MergeFromReader(baseYaml, patchYaml)
	if err != nil {
		return nil, err
	}

	return []byte(baseYaml.String()), nil
}

func cleanMap(m map[string]any) map[string]any {
	for k, v := range m {
		if isEmpty(v) {
			delete(m, k)
			continue
		}

		if nestedMap, ok := v.(map[string]any); ok {
			cleanNested := cleanMap(nestedMap)
			if len(cleanNested) == 0 {
				delete(m, k)
			} else {
				m[k] = cleanNested
			}
		}
	}

	return m
}

// Function to check if a value is empty
func isEmpty(v any) bool {
	// Check for nil
	if v == nil {
		return true
	}

	// Get the reflect value
	rv := reflect.ValueOf(v)

	// Check the kind to determine the method of checking empty
	switch rv.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		return rv.Len() == 0
	case reflect.Pointer, reflect.Interface:
		return rv.IsNil()
	}

	return false
}
