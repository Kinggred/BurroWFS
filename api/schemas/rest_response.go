package schemas

import (
	"burrowfs/api/common"
	"encoding/json"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	schemaStr, err := json.Marshal(data)
	if err != nil {
		common.HttpError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	_, err = w.Write(schemaStr)
	if err != nil {
		common.HttpError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

}
