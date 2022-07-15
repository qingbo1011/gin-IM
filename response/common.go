package response

type Response struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   any    `json:"data"`
	Error  string `json:"error"`
}
