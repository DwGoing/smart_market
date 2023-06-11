package app

type Request struct {
	Id     string `json:"id"`
	Method string `json:"method"`
}

type Response struct {
	Id      string `json:"id"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}
