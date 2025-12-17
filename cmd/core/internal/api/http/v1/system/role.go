package system

import (
	"github.com/gin-gonic/gin"
	"github.com/lyonmu/quebec/cmd/core/internal/dto/request"
	"github.com/lyonmu/quebec/cmd/core/internal/dto/response"
	"github.com/lyonmu/quebec/pkg/code"
)

// SystemRolePage
// @Tags      系统管理
// @Summary   角色分页列表
// @Description 获取角色分页列表
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Param     data  query      request.SystemRolePageReq      true  "角色列表信息"
// @Produce   json
// @Success   200  {object}  code.Response{code=number,data=response.SystemRoleListResp,message=string}  "50000,success"
// @Router    /v1/system/role/page [get]
func (b *SystemV1ApiGroup) SystemRolePage(c *gin.Context) {

	var req request.SystemRolePageReq
	var _ response.SystemRoleListResp
	if err := c.ShouldBindQuery(&req); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	resp, err := systemsvc.RolePage(c.Request.Context(), &req)
	if err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(resp, c)
}

// SystemRoleList
// @Tags      系统管理
// @Summary   全部角色列表
// @Description 获取全部角色列表
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Param     data  query      request.SystemRoleListReq      true  "角色列表信息"
// @Produce   json
// @Success   200  {object}  code.Response{code=number,data=response.SystemRoleListResp,message=string}  "50000,success"
// @Router    /v1/system/role/list [get]
func (b *SystemV1ApiGroup) SystemRoleList(c *gin.Context) {

	var req request.SystemRoleListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	resp, err := systemsvc.RoleList(c.Request.Context(), &req)
	if err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(resp, c)
}

// SystemOnlineUserLabel
// @Tags      系统管理
// @Summary   角色标签
// @Description 获取角色标签
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Produce   json
// @Success   200  {object}  code.Response{code=number,data=[]response.Options,message=string}  "50000,success"
// @Router    /v1/system/role/label [get]
func (b *SystemV1ApiGroup) SystemRoleLabel(c *gin.Context) {

	resp, err := systemsvc.RoleLabel(c.Request.Context())
	if err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(resp, c)
}

// SystemRoleDelete
// @Tags      系统管理
// @Summary   删除角色
// @Description 删除角色
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Param     data  path      request.IdReq      true  "角色ID"
// @Produce   json
// @Success   200  {object}  code.Response{code=number,message=string}  "50000,success"
// @Router    /v1/system/role/{id} [delete]
func (b *SystemV1ApiGroup) SystemRoleDelete(c *gin.Context) {

	var req request.IdReq
	if err := c.ShouldBindUri(&req); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	if err := systemsvc.RoleDelete(c.Request.Context(), req.ID); err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(nil, c)
}

// SystemRoleAdd
// @Tags      系统管理
// @Summary   添加角色信息
// @Description 添加角色信息
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Param     data  body      request.SystemRoleAddReq      true  "角色信息"
// @Produce   json
// @Success   200  {object}  code.Response{code=number,message=string}  "50000,success"
// @Router    /v1/system/role [post]
func (b *SystemV1ApiGroup) SystemRoleAdd(c *gin.Context) {

	var req request.SystemRoleAddReq
	if err := c.ShouldBindJSON(&req); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	if err := systemsvc.RoleAdd(c.Request.Context(), &req); err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(nil, c)
}

// SystemRoleEdit
// @Tags      系统管理
// @Summary   编辑角色信息
// @Description 编辑角色信息
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Param     data  path      request.IdReq      true  "角色ID"
// @Param     data  body      request.SystemRoleAddReq      true  "角色信息"
// @Produce   json
// @Success   200  {object}  code.Response{code=number,message=string}  "50000,success"
// @Router    /v1/system/role/{id} [put]
func (b *SystemV1ApiGroup) SystemRoleEdit(c *gin.Context) {

	var id request.IdReq
	if err := c.ShouldBindUri(&id); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	var req request.SystemRoleAddReq
	if err := c.ShouldBindJSON(&req); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	if err := systemsvc.RoleUpdate(c.Request.Context(), id.ID, &req); err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(nil, c)
}

// SystemRoleGetById
// @Tags      系统管理
// @Summary   获取角色详情
// @Description 获取角色详情
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Param     data  path      request.IdReq      true  "角色ID"
// @Produce   json
// @Success   200  {object}  code.Response{code=number,data=response.SystemRoleResp,message=string}  "50000,success"
// @Router    /v1/system/role/{id} [get]
func (b *SystemV1ApiGroup) SystemRoleGetById(c *gin.Context) {

	var id request.IdReq
	if err := c.ShouldBindUri(&id); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	resp, err := systemsvc.RoleGetById(c.Request.Context(), id.ID)
	if err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(resp, c)
}

// SystemRoleEnable
// @Tags      系统管理
// @Summary   启停角色
// @Description 启停角色状态
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Param     data  path      request.IdReq      true  "角色ID"
// @Param     data  body      request.EnableReq      true  "状态 [1: 启用, 2: 禁用]"
// @Produce   json
// @Success   200  {object}  code.Response{code=number,message=string}  "50000,success"
// @Router    /v1/system/role/enable/{id} [put]
func (b *SystemV1ApiGroup) SystemRoleEnable(c *gin.Context) {

	var id request.IdReq
	if err := c.ShouldBindUri(&id); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	var req request.EnableReq
	if err := c.ShouldBindJSON(&req); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	if err := systemsvc.RoleEnable(c.Request.Context(), id.ID, &req); err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(nil, c)
}
