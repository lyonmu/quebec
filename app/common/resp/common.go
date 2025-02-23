package resp

var (
	Success = HttpResponse{Code: 100000, Message: "操作成功"}
	Failed  = HttpResponse{Code: 999999, Message: "操作失败"}
)
