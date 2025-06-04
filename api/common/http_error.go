package common

import (
	"burrowfs/api/schemas"
	"encoding/json"
	"net/http"
)

func HttpError(w http.ResponseWriter, code int, detail string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := schemas.ErrorSchema{
		Status: code,
		Detail: detail,
	}

	schemaStr, err := json.Marshal(response)
	if err != nil {
		HttpError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	_, err = w.Write(schemaStr)
	if err != nil {
		HttpError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
}
