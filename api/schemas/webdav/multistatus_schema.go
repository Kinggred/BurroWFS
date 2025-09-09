package webdav

import (
	"encoding/xml"
)

type MultiStatus struct {
	XMLName   xml.Name
	Xmlns     string
	Responses []FileResponse
}
