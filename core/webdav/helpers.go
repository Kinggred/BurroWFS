package webdav

import "time"

func ExtractFileName(path string) string {
	if path == "" {
		return ""
	}
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	lastSlash := -1
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			lastSlash = i
			break
		}
	}
	if lastSlash == -1 {
		return path
	}
	return path[lastSlash+1:]
}

func ExtractParentPath(path string) string {
	if path == "" || path == "/" {
		return "/"
	}
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	lastSlash := -1
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			lastSlash = i
			break
		}
	}
	if lastSlash == -1 {
		return "/"
	}
	if lastSlash == 0 {
		return "/"
	}
	return path[:lastSlash]
}

func IsSubPath(parent, child string) bool {
	if len(parent) == 0 || len(child) == 0 {
		return false
	}
	if parent[len(parent)-1] != '/' {
		parent += "/"
	}
	return len(child) > len(parent) && child[:len(parent)] == parent
}

func FormatModifiedDate(t time.Time) string {
	if t.IsZero() {
		return "Mon, 01 Jan 1970 00:00:00 GMT"
	}
	return t.Format(time.RFC1123)
}

func FormatCreationDate(t time.Time) string {
	if t.IsZero() {
		return "1970-01-01T00:00:00Z"
	}
	return t.UTC().Format(time.RFC3339)
}
