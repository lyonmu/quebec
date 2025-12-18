package system

import (
	"github.com/gin-gonic/gin"
	"github.com/lyonmu/quebec/cmd/core/internal/dto/request"
	"github.com/lyonmu/quebec/cmd/core/internal/dto/response"
	"github.com/lyonmu/quebec/pkg/code"
)

// SystemOnlineUserList
// @Tags      系统管理
// @Summary   在线用户列表
// @Description 获取在线用户列表
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Param     data  query      request.SystemOnlineUserListReq      true  "在线用户列表信息"
// @Produce   json
// @Success   200  {object}  code.Response{code=number,data=response.SystemOnlineUserListResp,message=string}  "50000,success"
// @Router    /v1/system/onlineuser/list [get]
func (b *SystemV1ApiGroup) SystemOnlineUserList(c *gin.Context) {

	var req request.SystemOnlineUserListReq
	var _ response.SystemOnlineUserListResp
	if err := c.ShouldBindQuery(&req); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	resp, err := systemsvc.ListOnlineUser(&req, c.Request.Context())
	if err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(resp, c)
}

// SystemOnlineUserLabel
// @Tags      系统管理
// @Summary   在线用户标签
// @Description 获取在线用户标签
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Produce   json
// @Success   200  {object}  code.Response{code=number,data=[]response.Options,message=string}  "50000,success"
// @Router    /v1/system/onlineuser/label [get]
func (b *SystemV1ApiGroup) SystemOnlineUserLabel(c *gin.Context) {

	resp, err := systemsvc.OnlineUserLabel(c.Request.Context())
	if err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(resp, c)
}

// SystemOnlineUserClearance
// @Tags      系统管理
// @Summary   清除在线用户
// @Description 清除在线用户
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Param     data  path      request.IdReq      true  "在线用户ID"
// @Produce   json
// @Success   200  {object}  code.Response{code=number,message=string}  "50000,success"
// @Router    /v1/system/onlineuser/clearance/{id} [delete]
func (b *SystemV1ApiGroup) SystemOnlineUserClearance(c *gin.Context) {

	var req request.IdReq
	if err := c.ShouldBindUri(&req); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	err := systemsvc.Clearance(&req, c.Request.Context())
	if err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(nil, c)
}
