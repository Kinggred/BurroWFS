package common

import (
	"encoding/json"
	"io"
	"net/http"
)

func ParseBody(r *http.Request, v interface{}) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(r.Body)

	r.Body = http.MaxBytesReader(nil, r.Body, 1<<20)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
		// Error Handling 422
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		panic(err)
		// Error handling 422
	}
}
