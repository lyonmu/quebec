package resp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpResponse struct {
	Code    int         `json:"code,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func (r HttpResponse) Error() string {
	return r.Message
}

func (r HttpResponse) Response(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, r)
	return
}

func (r HttpResponse) ResponseWithData(ctx *gin.Context, data interface{}) {
	r.Data = data
	ctx.JSON(http.StatusOK, r)
	return
}
