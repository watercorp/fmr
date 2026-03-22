package filedetails

import (
	"path/filepath"
	"strings"
)

type FileDetails struct {
	FullPath  string
	Name      string
	Directory string
	Extension string
	BaseName  string
}

// Create a new instance of FileDetails from the full file path
func New(inputPath string) FileDetails {
	// Store some items we need to reuse
	basePath := filepath.Base(inputPath)
	ext := getFullExtension(basePath)

	return FileDetails{
		FullPath:  inputPath,
		Name:      basePath,
		Directory: filepath.Dir(inputPath),
		Extension: ext,
		BaseName:  strings.TrimSuffix(basePath, ext),
	}
}

// Returns the full path for the string
func (fd FileDetails) String() string {
	return fd.FullPath
}

// Get the extension from the first dot
func getFullExtension(name string) string {
	// Locate the index of the first .
	firstIndex := strings.Index(name, ".")

	// Confirm we've found an index and it's not at the beginning of the file
	if firstIndex >= 1 {
		return name[firstIndex:]
	} else {
		return ""
	}
}
