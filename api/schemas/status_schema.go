package schemas

type StatusSchema struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Time     string `json:"time"`
	Debug    bool   `json:"debug"`
}
