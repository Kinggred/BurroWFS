package utils

import (
	"strings"
)

// Path represents a parsed file or directory path
type Path struct {
	raw        string
	uri        string
	Name       string
	Clean      string
	ParentPath string
}

// RetrievePath parses the given raw path and returns a Path struct.
// If isUri is false, it treats the raw path as a full URL and extracts the path component.
// If isUri is true, it treats the raw path as a direct path.
// If the path is a directory, it should end with a "/".
func RetrievePath(raw string, isUri bool) *Path {
	path := Path{}
	path.uri = raw
	path.raw = raw

	if !isUri {
		path.uri = "/" + strings.Join(strings.Split(raw, "/")[3:], "/")
	}
	path.Clean = strings.TrimSuffix(path.uri, "/")

	splitPath := strings.Split(path.Clean, "/")
	if len(splitPath) == 2 {
		path.ParentPath = "/"
	} else {
		path.ParentPath = strings.Join(splitPath[:len(splitPath)-1], "/")
	}

	path.Name = splitPath[len(splitPath)-1]
	return &path

}

func (p *Path) IsFolder() bool {
	if p.Clean == "" {
		return false
	}
	// Treat root and any path ending with "/" as folders
	if p.Clean == "/" || strings.HasSuffix(p.raw, "/") {
		return true
	}
	return false
}

func (p *Path) IsParentRoot() bool {
	if p.ParentPath == "/" {
		return true
	}
	return false
}
