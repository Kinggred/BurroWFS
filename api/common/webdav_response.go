package common

import (
	"burrowfs/api/schemas"
	"encoding/xml"
	"net/http"
)

func WebDAVResponse(w http.ResponseWriter, code, message string, data any, status int) {
	resp := schemas.WebDAVResponse{
		Xmlns:   "DAV:",
		Code:    code,
		Message: message,
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	_ = xml.NewEncoder(w).Encode(resp)
}
