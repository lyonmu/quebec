package system

import (
	"github.com/gin-gonic/gin"
	"github.com/lyonmu/quebec/cmd/core/internal/dto/request"
	"github.com/lyonmu/quebec/cmd/core/internal/dto/response"
	"github.com/lyonmu/quebec/pkg/code"
	"github.com/lyonmu/quebec/pkg/tools"
)

// Login
// @Tags      系统管理
// @Summary   登陆
// @Description 使用用户名密码登陆系统
// @Param     data  body      request.SystemLoginRequest      true  "用户登陆信息"
// @Produce   json
// @Success   200  {object}  code.Response{code=number,data=response.SystemInfoResponse,message=string}  "50000,success"
// @Router    /v1/system/login [post]
func (b *SystemV1ApiGroup) SystemLogin(c *gin.Context) {

	var req request.SystemLoginRequest
	var _ response.SystemInfoResponse
	if err := c.ShouldBindJSON(&req); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	ua := tools.ParseUserAgent(c.Request.UserAgent())
	access_ip := c.ClientIP()
	resp, err := systemsvc.Login(&req, ua, access_ip, c.Request.Context())
	if err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(resp, c)
}
