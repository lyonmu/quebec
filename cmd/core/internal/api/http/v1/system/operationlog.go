package system

import (
	"github.com/gin-gonic/gin"
	"github.com/lyonmu/quebec/cmd/core/internal/dto/request"
	"github.com/lyonmu/quebec/cmd/core/internal/dto/response"
	"github.com/lyonmu/quebec/pkg/code"
)

// SystemOperationLogPage
// @Tags      系统管理
// @Summary   操作日志分页列表
// @Description 获取操作日志分页列表
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Param     data  query      request.OperationLogPageReq  true  "操作日志列表信息"
// @Produce   json
// @Success   200  {object}  code.Response{code=number,data=response.SystemOperationLogListResp,message=string}  "50000,success"
// @Router    /v1/system/operation-log/page [get]
func (b *SystemV1ApiGroup) SystemOperationLogPage(c *gin.Context) {
	var req request.OperationLogPageReq
	var _ response.SystemOperationLogListResp
	if err := c.ShouldBindQuery(&req); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	resp, err := systemsvc.OperationLogPage(c.Request.Context(), &req)
	if err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(resp, c)
}
