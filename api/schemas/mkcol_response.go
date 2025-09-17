package schemas

import "net/http"

func MkcolWebDavResponse(w http.ResponseWriter, status int, location string) {
	w.WriteHeader(status)
	if status == 201 {
		w.Header().Set("Location", location)
	}
	w.Header().Set("Content-Length", "0")
}
