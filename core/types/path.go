package types

import (
	"net/url"
	"strings"
)

// Path represents a parsed file or directory path
type Path struct {
	Raw        string
	Uri        string
	Name       string
	Clean      string
	ParentPath string
}

func (p *Path) IsFolder() bool {
	if p.Clean == "" {
		return false
	}
	// Treat root and any path ending with "/" as folders
	if p.Clean == "/" || strings.HasSuffix(p.Raw, "/") {
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

func (p *Path) SplitAndEncode(raw string, sep string) []string {
	parts := strings.Split(raw, sep)
	for i, part := range parts {
		parts[i] = url.PathEscape(part)
	}
	return parts
}
