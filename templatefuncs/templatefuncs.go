package templatefuncs

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"html/template"
	"io"
	"os"
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
	"json":       ToJson,
	"lower":      strings.ToLower,
	"md5File":    Md5File,
	"part":       StringPart,
	"replace":    ReplaceString,
	"sha1File":   Sha1File,
	"sha256File": Sha256File,
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

// Create a shared function to handle file hasing
func hashFile(filePath string, hashAlg string) (string, error) {

	// Open the file to read the contents
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Initialize the hasher
	var hasher hash.Hash = nil

	// Create the hasher
	switch hashAlg {
	case "md5":
		hasher = md5.New()
	case "sha1":
		hasher = sha1.New()
	case "sha256":
		hasher = sha256.New()
	default:
		return "", fmt.Errorf("invalid hash specified")
	}

	// Stream the file to the hasher
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	// Compute the bytes
	hashBytes := hasher.Sum(nil)

	// Return the string representation
	return hex.EncodeToString(hashBytes), nil
}

// Joins all passed strings using the delimiter
func JoinStrings(d string, parts ...string) string {
	return strings.Join(parts, d)
}

// Generates the MD5 hash of the supplied file
func Md5File(filePath string) (string, error) {
	return hashFile(filePath, "md5")
}

// Replaces a regex pattern in a string with the replacement
func ReplaceString(pattern string, source string, replacement string) (string, error) {
	// Compile the regex
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}

	// Return the replacement
	return re.ReplaceAllString(source, replacement), nil
}

// Generates the SHA1 hash of the supplied file
func Sha1File(filePath string) (string, error) {
	return hashFile(filePath, "sha1")
}

// Generates the SHA256 hash of the supplied file
func Sha256File(filePath string) (string, error) {
	return hashFile(filePath, "sha256")
}

// Returns the first part of a FQDN split by "."
func ShortFqdn(input string) (string, error) {
	return StringPart(input, 0, ".")
}

// Takes the string and returns the specified part after being split by the delimiter
func StringPart(input string, i int, d string) (string, error) {
	parts := strings.Split(input, d)

	if i < 0 || i >= len(parts) {
		return "", fmt.Errorf("Out of bounds")
	}

	return parts[i], nil
}

// Title cases the input string
func TitleCase(input string) string {
	return titleCaser.String(input)
}

// Trims the prefix from a string
func TrimPrefix(prefix string, s string) string {
	return strings.TrimPrefix(s, prefix)
}

// Trims the suffix from a string
func TrimSuffix(suffix string, s string) string {
	return strings.TrimSuffix(s, suffix)
}

// Converts the data to a JSON string
func ToJson(data any) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
