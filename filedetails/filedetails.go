package filedetails

import (
	"path"
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
func New(filePath string) FileDetails {
	// Store some items we need to reuse
	basePath := path.Base(filePath)
	ext := getFullExtension(basePath)

	return FileDetails{
		FullPath:  filePath,
		Name:      basePath,
		Directory: path.Dir(filePath),
		Extension: ext,
		BaseName:  basePath[:len(basePath)-len(ext)],
	}
}

// Returns the full path for the string
func (fd FileDetails) String() string {
	return fd.FullPath
}

// Get the extension from the first dot
func getFullExtension(baseName string) string {
	// Locate the index of the first .
	firstIndex := strings.Index(baseName, ".")

	// Confirm we've found an index and it's not at the beginning of the file
	if firstIndex >= 1 {
		return baseName[firstIndex+1:]
	} else {
		return ""
	}
}
