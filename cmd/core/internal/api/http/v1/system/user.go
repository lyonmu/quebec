package system

import (
	"github.com/gin-gonic/gin"
	"github.com/lyonmu/quebec/cmd/core/internal/dto/request"
	"github.com/lyonmu/quebec/cmd/core/internal/dto/response"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/pkg/code"
	"github.com/lyonmu/quebec/pkg/constant"
)

// SystemUserPage
// @Tags      系统管理
// @Summary   用户分页列表
// @Description 获取用户分页列表
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Param     data  query      request.SystemUserPageReq      true  "用户列表信息"
// @Produce   json
// @Success   200  {object}  code.Response{code=number,data=response.SystemUserListResp,message=string}  "50000,success"
// @Router    /v1/system/user/page [get]
func (b *SystemV1ApiGroup) SystemUserPage(c *gin.Context) {

	var req request.SystemUserPageReq
	var _ response.SystemUserListResp
	if err := c.ShouldBindQuery(&req); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	resp, err := systemsvc.UserPage(c.Request.Context(), &req)
	if err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(resp, c)
}

// SystemRoleList
// @Tags      系统管理
// @Summary   全部用户列表
// @Description 获取全部用户列表
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Param     data  query      request.SystemUserListReq      true  "用户列表信息"
// @Produce   json
// @Success   200  {object}  code.Response{code=number,data=response.SystemUserListResp,message=string}  "50000,success"
// @Router    /v1/system/user/list [get]
func (b *SystemV1ApiGroup) SystemUserList(c *gin.Context) {

	var req request.SystemUserListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	resp, err := systemsvc.UserList(c.Request.Context(), &req)
	if err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(resp, c)
}

// SystemOnlineUserLabel
// @Tags      系统管理
// @Summary   用户标签
// @Description 获取用户标签
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Produce   json
// @Success   200  {object}  code.Response{code=number,data=[]response.Options,message=string}  "50000,success"
// @Router    /v1/system/user/label [get]
func (b *SystemV1ApiGroup) SystemUserLabel(c *gin.Context) {

	resp, err := systemsvc.UserLabel(c.Request.Context())
	if err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(resp, c)
}

// SystemRoleDelete
// @Tags      系统管理
// @Summary   删除用户
// @Description 删除用户
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Param     data  path      request.IdReq      true  "用户ID"
// @Produce   json
// @Success   200  {object}  code.Response{code=number,message=string}  "50000,success"
// @Router    /v1/system/user/{id} [delete]
func (b *SystemV1ApiGroup) SystemUserDelete(c *gin.Context) {

	var req request.IdReq
	if err := c.ShouldBindUri(&req); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	if err := systemsvc.UserDelete(c.Request.Context(), req.ID); err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(nil, c)
}

// SystemRoleAdd
// @Tags      系统管理
// @Summary   添加用户信息
// @Description 添加用户信息
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Param     data  body      request.SystemUserAddReq      true  "用户信息"
// @Produce   json
// @Success   200  {object}  code.Response{code=number,message=string}  "50000,success"
// @Router    /v1/system/user [post]
func (b *SystemV1ApiGroup) SystemUserAdd(c *gin.Context) {

	var req request.SystemUserAddReq
	if err := c.ShouldBindJSON(&req); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	if err := systemsvc.UserAdd(c.Request.Context(), &req); err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(nil, c)
}

// SystemRoleEdit
// @Tags      系统管理
// @Summary   编辑用户信息
// @Description 编辑用户信息
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Param     data  path      request.IdReq      true  "用户ID"
// @Param     data  body      request.SystemUserUpdateReq      true  "用户信息"
// @Produce   json
// @Success   200  {object}  code.Response{code=number,message=string}  "50000,success"
// @Router    /v1/system/user/{id} [put]
func (b *SystemV1ApiGroup) SystemUserEdit(c *gin.Context) {

	var id request.IdReq
	if err := c.ShouldBindUri(&id); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	var req request.SystemUserUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	if err := systemsvc.UserUpdate(c.Request.Context(), id.ID, &req); err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(nil, c)
}

// SystemRoleGetById
// @Tags      系统管理
// @Summary   获取用户详情
// @Description 获取用户详情
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Param     data  path      request.IdReq      true  "用户ID"
// @Produce   json
// @Success   200  {object}  code.Response{code=number,data=response.SystemUserResp,message=string}  "50000,success"
// @Router    /v1/system/user/{id} [get]
func (b *SystemV1ApiGroup) SystemUserGetById(c *gin.Context) {

	var id request.IdReq
	if err := c.ShouldBindUri(&id); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	resp, err := systemsvc.UserGetById(c.Request.Context(), id.ID)
	if err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(resp, c)
}

// SystemRoleEnable
// @Tags      系统管理
// @Summary   启停用户
// @Description 启停用户状态
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Param     data  path      request.IdReq      true  "用户ID"
// @Param     data  body      request.EnableReq      true  "状态 [1: 启用, 2: 禁用]"
// @Produce   json
// @Success   200  {object}  code.Response{code=number,message=string}  "50000,success"
// @Router    /v1/system/user/enable/{id} [put]
func (b *SystemV1ApiGroup) SystemUserEnable(c *gin.Context) {

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

	if err := systemsvc.UserEnable(c.Request.Context(), id.ID, &req); err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(nil, c)
}

// SystemRoleEditPassword
// @Tags      系统管理
// @Summary   修改用户密码
// @Description 修改用户密码
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Param     data  path      request.IdReq      true  "用户ID"
// @Param     data  body      request.SystemUserEditPasswordReq      true  "用户密码信息"
// @Produce   json
// @Success   200  {object}  code.Response{code=number,message=string}  "50000,success"
// @Router    /v1/system/user/password/{id} [put]
func (b *SystemV1ApiGroup) SystemUserEditPassword(c *gin.Context) {

	var id request.IdReq
	if err := c.ShouldBindUri(&id); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	var req request.SystemUserEditPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		code.PasswordIncorrect.Failed(c)
		return
	}

	if err := systemsvc.UserEditPassword(c.Request.Context(), id.ID, &req); err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(nil, c)
}

// SystemRoleEditPassword
// @Tags      系统管理
// @Summary   修改自己的密码
// @Description 修改自己的密码
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Param     data  body      request.SystemUserEditPasswordReq      true  "用户密码信息"
// @Produce   json
// @Success   200  {object}  code.Response{code=number,message=string}  "50000,success"
// @Router    /v1/system/user/password/self [put]
func (b *SystemV1ApiGroup) SystemUserEditPasswordSelf(c *gin.Context) {

	var req request.SystemUserEditPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		code.PasswordIncorrect.Failed(c)
		return
	}

	claims, err := global.JwtToolEntity.ParseToken(c.GetHeader(constant.ApiTokenName), global.Cfg.Core.Jwt.Sign)
	if err != nil {
		code.Unauthorized.Unauthorized(c)
		return
	}

	if err := systemsvc.UserEditPassword(c.Request.Context(), claims.UserId, &req); err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(nil, c)
}
