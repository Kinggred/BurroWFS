package rest

// PutFilesInputSchema defines the expected input for uploading a file.
type PutFilesInputSchema struct {
	Items FileBodySchema `json:"items"`
}

type FileBodySchema struct {
	RelativePath string `json:"relative_path"`
	ContentType  string `json:"content_type"`
	Data         []byte `json:"data"`
}
