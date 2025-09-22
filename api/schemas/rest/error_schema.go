package rest

type ErrorSchema struct {
	Status int    `json:"status"`
	Detail string `json:"error"`
}
