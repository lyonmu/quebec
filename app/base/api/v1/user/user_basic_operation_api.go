package user

import (
	"github.com/gin-gonic/gin"
	commonresp "github.com/lyonmu/quebec/app/common/resp"
)

type BasicUserOperationApiV1 struct{}

func (BasicUserOperationApiV1) Login(ctx *gin.Context) {

	if res := userService.Login(ctx); res != nil {
		res.Response(ctx)
		return
	}
	commonresp.Success.Response(ctx)
}

func (BasicUserOperationApiV1) Signup(ctx *gin.Context) {
	if res := userService.Signup(ctx); res != nil {
		res.Response(ctx)
		return
	}
	commonresp.Success.Response(ctx)

}

func (BasicUserOperationApiV1) Captcha(ctx *gin.Context) {
	if res := userService.Captcha(ctx); res != nil {
		res.Response(ctx)
		return
	}
	commonresp.Success.Response(ctx)
}
