package schemas

import "net/http"

func MoveWebDavResponse(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}
