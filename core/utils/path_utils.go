package utils

import (
	"burrowfs/core/types"
	"strings"
)

// RetrievePath parses the given raw path and returns a Path struct.
// If isUri is false, it treats the raw path as a full URL and extracts the path component.
// If isUri is true, it treats the raw path as a direct path.
// If the path is a directory, it should end with a "/".
func RetrievePath(raw string, isUri bool) *types.Path {
	path := types.Path{}
	path.Uri = raw
	path.Raw = raw

	if !isUri {
		path.Uri = "/" + strings.Join(path.SplitAndEncode(raw, "/")[3:], "/")
	}
	if path.Uri == "/" {
		path.Clean = "/"
	} else {
		path.Clean = strings.Join(path.SplitAndEncode(strings.TrimSuffix(path.Uri, "/"), "/"), "/")
	}

	splitPath := strings.Split(path.Clean, "/")
	if len(splitPath) == 2 {
		path.ParentPath = "/"
	} else {
		path.ParentPath = strings.Join(splitPath[:len(splitPath)-1], "/")
	}

	path.Name = splitPath[len(splitPath)-1]
	return &path

}
