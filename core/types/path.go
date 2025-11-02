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

func (p *Path) safePathMerge(a string, b string) string {
	return strings.TrimSuffix(a, "/") + "/" + strings.TrimPrefix(b, "/")
}

// MergePaths appends two paths together
// Path which calls this method is treated as the base path
// its values will be overwritten by the resulting merged path
// returns self
func (p *Path) MergePaths(addon Path) *Path {
	p.Clean = p.safePathMerge(p.Clean, addon.Clean)
	p.Raw = p.safePathMerge(p.Raw, addon.Raw)
	p.Uri = p.safePathMerge(p.Uri, addon.Uri)

	splitPath := strings.Split(p.Clean, "/")
	p.Name = splitPath[len(splitPath)-1]
	if len(splitPath) == 2 {
		p.ParentPath = "/"
	} else {
		p.ParentPath = strings.Join(splitPath[:len(splitPath)-1], "/")
	}
	return p
}
