package common

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
)

// ParseJSON parses the JSON body of an HTTP request into the provided interface.
func ParseJSON(r *http.Request, v interface{}, validate bool) error {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(r.Body)

	r.Body = http.MaxBytesReader(nil, r.Body, 1<<20)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return err
	}

	if validate {
		err = ValidateStruct(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func ValidateStruct(v interface{}) error {
	if validatable, ok := v.(Validatable); ok {
		return validatable.Validate()
	}
	return nil
}

func ParseXML(r *http.Request, v interface{}) error {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(r.Body)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(body, v)
	return err
}
