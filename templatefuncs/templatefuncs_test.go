package templatefuncs_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/watercorp/fmr/templatefuncs"
)

func TestTemplateFuncs_DistinguishedName(t *testing.T) {
	// Define successful test cases
	testsSuccess := []struct {
		input    []string
		expected string
	}{
		{
			input:    []string{"one.example.com"},
			expected: "DC=one,DC=example,DC=com",
		},
		{
			input:    []string{"com"},
			expected: "DC=com",
		},
		{
			input:    []string{"example.com"},
			expected: "DC=example,DC=com",
		},
		{
			input:    []string{"one.example.com", "main.users"},
			expected: "OU=main,OU=users,DC=one,DC=example,DC=com",
		},
		{
			input:    []string{"one.example.com", "main.users", "username"},
			expected: "CN=username,OU=main,OU=users,DC=one,DC=example,DC=com",
		},
	}

	// Loop through each test case
	for i, tc := range testsSuccess {
		t.Run(fmt.Sprintf("distinguished name success %d", i), func(t *testing.T) {
			// Generate the distinguished name
			dn, err := templatefuncs.DistinguishedName(tc.input...)

			// Confirm there is no error
			if err != nil {
				t.Errorf("expected no error, got %s", err)
			}

			// Check if the values match
			if dn != tc.expected {
				t.Errorf("values do not match. have: %s want: %s", dn, tc.expected)
			}
		})
	}

	// Define error test cases
	testsError := [][]string{
		{},
		{"one", "two", "three", "four"},
	}

	// Loop through each test case
	for i, tc := range testsError {
		t.Run(fmt.Sprintf("distinguished name error %d", i), func(t *testing.T) {
			// Generate the distinguished name
			_, err := templatefuncs.DistinguishedName(tc...)

			// Confirm we have an error
			if err == nil {
				t.Error("expected error, got none")
			}
		})
	}
}

func TestTemplateFuncs_FileHashes(t *testing.T) {
	// Define the successful test cases
	tests := []struct {
		filePath       string
		expectedMd5    string
		expectedSha1   string
		expectedSha256 string
	}{
		{
			filePath:       "testdata/hash-file1.txt",
			expectedMd5:    "0f980cfb2feae586c7e64a9fbd9d6c83",
			expectedSha1:   "9b7105c2dc32cedf5762f3f82570ba72ffbad4e5",
			expectedSha256: "b3c085878c28adfb0338399271db8a54fe4f385e28a07cc9241b3fe994432e34",
		},
		{
			filePath:       "testdata/hash-file2.txt",
			expectedMd5:    "732bde0bf2c5d4657970c02bdd449b01",
			expectedSha1:   "817971c2cc10ebfd7aa8e8bbf30fbf79b19b167e",
			expectedSha256: "ebd64bb69e221418b7691f13cbe2823c9a5cca8cae32f60ffde253ed99665a97",
		},
		{
			filePath:       "testdata/hash-file3.txt",
			expectedMd5:    "0063bb8a0936c12455aaf9f7d73d5d07",
			expectedSha1:   "d58906e53c072c45f3ffb3701f1aee71acf4d4f0",
			expectedSha256: "199f23a200071729bf94f38973a842e2834abd30b865af1b5e147c139efe5b5a",
		},
	}

	// Loop through each test case
	for i, tc := range tests {
		t.Run(fmt.Sprintf("md5 file hash %d", i), func(t *testing.T) {
			// Call the Md5File function
			computedHash, err := templatefuncs.Md5File(tc.filePath)

			// Confirm there is no error
			if err != nil {
				t.Errorf("expected no error, got %s", err)
			}

			// Confirm the values match
			if computedHash != tc.expectedMd5 {
				t.Errorf("values do not match, have %s want %s", computedHash, tc.expectedMd5)
			}
		})

		t.Run(fmt.Sprintf("sha1 file hash %d", i), func(t *testing.T) {
			// Call the Sha1File function
			computedHash, err := templatefuncs.Sha1File(tc.filePath)

			// Confirm there is no error
			if err != nil {
				t.Errorf("expected no error, got %s", err)
			}

			// Confirm the values match
			if computedHash != tc.expectedSha1 {
				t.Errorf("values do not match, have %s want %s", computedHash, tc.expectedSha1)
			}
		})

		t.Run(fmt.Sprintf("sha256 file hash %d", i), func(t *testing.T) {
			// Call the Sha256File function
			computedHash, err := templatefuncs.Sha256File(tc.filePath)

			// Confirm there is no error
			if err != nil {
				t.Errorf("expected no error, got %s", err)
			}

			// Confirm the values match
			if computedHash != tc.expectedSha256 {
				t.Errorf("values do not match, have %s want %s", computedHash, tc.expectedSha256)
			}
		})
	}

	testsError := []string{
		"testdata/not-a-real-file.txt",
	}

	// Loop through error test cases
	for i, tc := range testsError {
		t.Run(fmt.Sprintf("md5 file hash error %d", i), func(t *testing.T) {
			// Call the Md5File function
			_, err := templatefuncs.Md5File(tc)

			// Confirm we have an error
			if err == nil {
				t.Error("expected error, got none")
			}
		})

		t.Run(fmt.Sprintf("sha256 file hash error %d", i), func(t *testing.T) {
			// Call the Sha256 function
			_, err := templatefuncs.Sha256File(tc)

			// Confirm we have an error
			if err == nil {
				t.Error("expected error, got none")
			}
		})
	}
}

func TestTemplateFuncs_JoinStrings(t *testing.T) {
	// Define the successful test cases
	tests := []struct {
		delim    string
		parts    []string
		expected string
	}{
		{
			delim:    ".",
			parts:    []string{"one", "two"},
			expected: "one.two",
		},
		{
			delim:    ".",
			parts:    []string{"one"},
			expected: "one",
		},
		{
			delim:    " ",
			parts:    []string{"one", "two", "three"},
			expected: "one two three",
		},
		{
			delim:    "...",
			parts:    []string{"one", "two", "three"},
			expected: "one...two...three",
		},
	}

	// Loop through each test case
	for i, tc := range tests {
		t.Run(fmt.Sprintf("join strings %d", i), func(t *testing.T) {
			// Call the JoinStrings function
			joined := templatefuncs.JoinStrings(tc.delim, tc.parts...)

			// Confirm the values match
			if joined != tc.expected {
				t.Errorf("values do not match, have %s want %s", joined, tc.expected)
			}
		})
	}
}

func TestTemplateFuncs_ReplaceString(t *testing.T) {
	// Define successful test cases
	testsSuccess := []struct {
		pattern     string
		source      string
		replacement string
		expected    string
	}{
		{
			pattern:     `-`,
			source:      "One-Two",
			replacement: "",
			expected:    "OneTwo",
		},
	}

	// Loop through each test case
	for i, tc := range testsSuccess {
		t.Run(fmt.Sprintf("replace string success %d", i), func(t *testing.T) {
			// Call the ReplaceString function
			output, err := templatefuncs.ReplaceString(tc.pattern, tc.source, tc.replacement)

			// Confirm there was no error
			if err != nil {
				t.Errorf("no error expected, got %s", err)
			}

			// Confirm the values match
			if output != tc.expected {
				t.Errorf("values do not match, have: %s want: %s", output, tc.expected)
			}
		})
	}

	// Define error test cases
	testsError := []struct {
		pattern     string
		source      string
		replacement string
	}{
		{
			pattern:     `\`,
			source:      "One-Two",
			replacement: "",
		},
	}

	// Loop through each test case
	for i, tc := range testsError {
		t.Run(fmt.Sprintf("replace string error %d", i), func(t *testing.T) {
			// Call the ReplaceString function
			_, err := templatefuncs.ReplaceString(tc.pattern, tc.source, tc.replacement)

			// Confirm there was an error
			if err == nil {
				t.Error("expected an error but got none")
			}
		})
	}
}

func TestTemplateFuncs_ShortFqdn(t *testing.T) {
	// Define test cases
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "one.example.com",
			expected: "one",
		},
		{
			input:    "one",
			expected: "one",
		},
		{
			input:    "one,example,com",
			expected: "one,example,com",
		},
		{
			input:    "",
			expected: "",
		},
	}

	// Loop through test cases
	for i, tc := range tests {
		t.Run(fmt.Sprintf("short fqdn %d", i), func(t *testing.T) {
			// Call the ShortFqdn function
			output, err := templatefuncs.ShortFqdn(tc.input)

			// Confirm there is no error
			if err != nil {
				t.Errorf("expected no error, got %s", err)
			}

			// Confirm the values match
			if output != tc.expected {
				t.Errorf("values do not match, want: %s have %s", output, tc.expected)
			}
		})
	}
}

func TestTemplateFuncs_StringPart(t *testing.T) {
	// Define successful test cases
	testsSuccess := []struct {
		input    string
		index    int
		delim    string
		expected string
	}{
		{
			input:    "one.example.com",
			index:    0,
			delim:    ".",
			expected: "one",
		},
		{
			input:    "one.example.com",
			index:    2,
			delim:    ".",
			expected: "com",
		},
		{
			input:    "one,two,three",
			index:    1,
			delim:    ",",
			expected: "two",
		},
		{
			input:    "one",
			index:    0,
			delim:    ",",
			expected: "one",
		},
		{
			input:    "",
			index:    0,
			delim:    ",",
			expected: "",
		},
	}

	// Loop through each test case
	for i, tc := range testsSuccess {
		t.Run(fmt.Sprintf("string part success %d", i), func(t *testing.T) {
			// Call the StringPart function
			output, err := templatefuncs.StringPart(tc.input, tc.index, tc.delim)

			// Confirm there is no error
			if err != nil {
				t.Errorf("expected no error, got %s", err)
			}

			// Confirm the values match
			if output != tc.expected {
				t.Errorf("values do not match have: %s want: %s", output, tc.expected)
			}
		})
	}

	// Define error test cases
	testsError := []struct {
		input string
		index int
		delim string
	}{
		{
			input: "one.example.com",
			index: -1,
			delim: ".",
		},
		{
			input: "one.example.com",
			index: 3,
			delim: ".",
		},
	}

	// Loop through each test case
	for i, tc := range testsError {
		t.Run(fmt.Sprintf("string part error %d", i), func(t *testing.T) {
			// Call the StringPart function
			_, err := templatefuncs.StringPart(tc.input, tc.index, tc.delim)

			// Confirm there was an error
			if err == nil {
				t.Error("expected an error but got none")
			}
		})
	}
}

func TestTemplateFuncs_TitleCase(t *testing.T) {
	// Define test cases
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "title",
			expected: "Title",
		},
	}

	// Loop through test cases
	for i, tc := range tests {
		t.Run(fmt.Sprintf("title case %d", i), func(t *testing.T) {
			// Call TitleCase function
			output := templatefuncs.TitleCase(tc.input)

			// Confirm the values match
			if output != tc.expected {
				t.Errorf("values do not match, have: %s want: %s", output, tc.expected)
			}
		})
	}
}

func TestTemplateFuncs_TrimPrefix(t *testing.T) {
	// Define test cases
	tests := []struct {
		prefix   string
		input    string
		expected string
	}{
		{
			prefix:   ".",
			input:    ".name",
			expected: "name",
		},
	}

	// Loop through test cases
	for i, tc := range tests {
		t.Run(fmt.Sprintf("trime prefix %d", i), func(t *testing.T) {
			// Call TrimPrefix function
			output := templatefuncs.TrimPrefix(tc.prefix, tc.input)

			// Confirm the values match
			if output != tc.expected {
				t.Errorf("values do not match, have: %s want: %s", output, tc.expected)
			}
		})
	}
}

func TestTemplateFuncs_TrimSuffix(t *testing.T) {
	// Define test cases
	tests := []struct {
		suffix   string
		input    string
		expected string
	}{
		{
			suffix:   ".",
			input:    "name.",
			expected: "name",
		},
	}

	// Loop through test cases
	for i, tc := range tests {
		t.Run(fmt.Sprintf("trime suffix %d", i), func(t *testing.T) {
			// Call TrimSuffix function
			output := templatefuncs.TrimSuffix(tc.suffix, tc.input)

			// Confirm the values match
			if output != tc.expected {
				t.Errorf("values do not match, have: %s want: %s", output, tc.expected)
			}
		})
	}
}

func TestTemplateFuncs_ToJson(t *testing.T) {
	// Define successful test cases
	testsSuccess := []struct {
		data     any
		expected string
	}{
		{
			data: map[string]any{
				"one":     "value 1",
				"two":     "value 2",
				"listone": []string{"one", "two"},
			},
			expected: `{"listone":["one","two"],"one":"value 1","two":"value 2"}`,
		},
		{
			data:     []string{"one", "two"},
			expected: `["one","two"]`,
		},
		{
			data:     "one",
			expected: `"one"`,
		},
		{
			data: map[string]any{
				"one": map[string]string{
					"two": "three",
				},
			},
			expected: `{"one":{"two":"three"}}`,
		},
	}

	// Loop through each test case
	for i, tc := range testsSuccess {
		t.Run(fmt.Sprintf("to json success %d", i), func(t *testing.T) {
			// Call the ToJson function
			j, err := templatefuncs.ToJson(tc.data)

			// Confirm there is no error
			if err != nil {
				t.Errorf("expected no error, got %s", err)
			}

			// Confirm the output matched
			if j != tc.expected {
				t.Errorf("json does not match, have: %s want: %s", j, tc.expected)
			}
		})
	}

	// Define error test cases
	testsError := []any{
		math.NaN(),
	}

	// Loop through each test base
	for i, tc := range testsError {
		t.Run(fmt.Sprintf("to json error %d", i), func(t *testing.T) {
			// Call the ToJson function
			_, err := templatefuncs.ToJson(tc)

			// Confirm we received an error
			if err == nil {
				t.Error("expected error and got none")
			}
		})
	}
}
