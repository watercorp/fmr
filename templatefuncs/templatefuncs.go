package templatefuncs

import (
	"fmt"
	"html/template"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Define reuseable items
var titleCaser = cases.Title(language.Und)

// Define the func map
var FuncMap = template.FuncMap{
	"title":      TitleCase,
	"upper":      strings.ToUpper,
	"lower":      strings.ToLower,
	"trimprefix": trimPrefix,
	"trimsuffix": trimSuffix,
	"join":       strings.Join,
	"joinstr":    JoinStrings,
	"part":       StringPart,
	"shortFqdn":  ShortFqdn,
	"dn":         DistinguishedName,
}

// Fix some function parameter orders
func trimPrefix(prefix string, s string) string {
	return strings.TrimPrefix(s, prefix)
}

func trimSuffix(suffix string, s string) string {
	return strings.TrimSuffix(s, suffix)
}

func JoinStrings(d string, parts ...string) string {
	// Create a new slice from the parts
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

// Returns the first part of a FQDN split by "."
func ShortFqdn(input string) (string, error) {
	return StringPart(input, 0, ".")
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

// Title cases the input string
func TitleCase(input string) string {
	return titleCaser.String(input)
}
