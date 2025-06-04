package schemas

type ErrorSchema struct {
	Status int    `json:"status"`
	Detail string `json:"error"`
}
