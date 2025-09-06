package schemas

import "encoding/xml"

// WebDAVResponse is a standardized XML response for WebDAV operations.
type WebDAVResponse struct {
	XMLName xml.Name `xml:"d:response"`
	Xmlns   string   `xml:"xmlns:d,attr"`
	Code    string   `xml:"d:code"`
	Message string   `xml:"d:message"`
	Data    any      `xml:"d:data,omitempty"`
}
