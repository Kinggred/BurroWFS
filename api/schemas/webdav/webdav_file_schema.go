package webdav

import "encoding/xml"

type Prop struct {
	DisplayName   string
	ContentLength int64
	ContentType   string
	LastModified  string
	ETag          string
	ResourceType  bool // true if collection
	CreationDate  string
	Permissions   string
	LockInfo      string
}

// PropStat represents the propstat element in WebDAV responses
type PropStat struct {
	Prop   Prop
	Status string
}

// FileResponse represents a single file or folder response in WebDAV PROPFIND
// This matches the WebDAV XML schema for a resource response
type FileResponse struct {
	XMLName  xml.Name
	Href     string
	PropStat PropStat
}

// ResponseCollection represents a collection of file responses in a WebDAV multistatus response
type ResponseCollection struct {
	XMLName   xml.Name
	Responses []FileResponse
}
