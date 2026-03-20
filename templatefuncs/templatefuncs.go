package templatefuncs

import (
	"fmt"
	"html/template"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Define reuseable items
var titleCaser = cases.Title(language.Und)

// Define the func map
var FuncMap = template.FuncMap{
	"dn":         DistinguishedName,
	"join":       strings.Join,
	"joinstr":    JoinStrings,
	"lower":      strings.ToLower,
	"part":       StringPart,
	"replace":    ReplaceString,
	"shortFqdn":  ShortFqdn,
	"title":      TitleCase,
	"trimprefix": TrimPrefix,
	"trimsuffix": TrimSuffix,
	"upper":      strings.ToUpper,
}

// Builds a distinguished name from the supplied values. Domain, OUs (dot formatted), CN
func DistinguishedName(values ...string) (string, error) {
	// Check how many strings we were passed
	if len(values) > 3 {
		return "", fmt.Errorf("too many values passed")
	} else if len(values) < 1 {
		return "", fmt.Errorf("no values passed. one required")
	}

	dnParts := []string{}
	domainParts := strings.Split(values[0], ".")

	// Handle CN if specified
	if len(values) == 3 {
		dnParts = append(dnParts, fmt.Sprintf("CN=%s", values[2]))
	}

	// Handle OUs if specified
	if len(values) >= 2 {
		ouParts := strings.Split(values[1], ".")
		for _, o := range ouParts {
			dnParts = append(dnParts, fmt.Sprintf("OU=%s", o))
		}
	}

	// Handle domain
	for _, d := range domainParts {
		dnParts = append(dnParts, fmt.Sprintf("DC=%s", d))
	}

	return strings.Join(dnParts, ","), nil
}

func JoinStrings(d string, parts ...string) string {
	return strings.Join(parts, d)
}

// Takes the string and returns the specified part after being split by the delimiter
func StringPart(input string, i int, d string) (string, error) {
	parts := strings.Split(input, d)

	if i < 0 || i >= len(parts) {
		return "", fmt.Errorf("Out of bounds")
	}

	return parts[i], nil
}

func ReplaceString(pattern string, source string, replacement string) (string, error) {
	// Compile the regex
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}

	// Return the replacement
	return re.ReplaceAllString(source, replacement), nil
}

// Returns the first part of a FQDN split by "."
func ShortFqdn(input string) (string, error) {
	return StringPart(input, 0, ".")
}

// Title cases the input string
func TitleCase(input string) string {
	return titleCaser.String(input)
}

// Fix some function parameter orders
func TrimPrefix(prefix string, s string) string {
	return strings.TrimPrefix(s, prefix)
}

func TrimSuffix(suffix string, s string) string {
	return strings.TrimSuffix(s, suffix)
}
