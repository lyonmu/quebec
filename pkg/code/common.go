package code

var (
	Success       = Response{Code: 50000, Message: "success"}
	Failed        = Response{Code: 59999, Message: "failed"}
	InvalidParams = Response{Code: 50001, Message: "invalid parameters"}
	Unauthorized  = Response{Code: 50401, Message: "unauthorized"}
	Forbidden     = Response{Code: 50403, Message: "forbidden"}
)
