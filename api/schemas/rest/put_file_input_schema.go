package rest

type PutFileInputSchema struct {
	ContentType string `json:"content_type"`
	Data        []byte `json:"data"` // base64 encoded data

}
