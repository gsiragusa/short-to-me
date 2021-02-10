package api

type ResponseApi struct {
	Status    string `json:"status"`
	Operation string `json:"operation"`
	Url       string `json:"url"`
}

type ResponseCount struct {
	Status    string `json:"status"`
	Operation string `json:"operation"`
	Count     int64  `json:"count"`
}
