package schemas

import (
	"burrowfs/api/schemas/webdav"
	"burrowfs/core/logging"
	"encoding/xml"
	"net/http"
)

func LockWebDavResponse(w http.ResponseWriter, status int, resp webdav.LockResponse) {
	logger := logging.Get("xmlResponse/lock")
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(status)

	enc := xml.NewEncoder(w)
	enc.Indent(" ", " ")
	if err := enc.Encode(resp); err != nil {
		logger.Error(err.Error())
		http.Error(w, "Failed to write XML response", http.StatusInternalServerError)
		return
	}
	_ = enc.Flush()
}
