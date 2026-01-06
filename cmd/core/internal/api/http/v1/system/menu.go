package system

import (
	"github.com/gin-gonic/gin"
	"github.com/lyonmu/quebec/cmd/core/internal/dto/request"
	"github.com/lyonmu/quebec/pkg/code"
)

// SystemMenuPage
// @Tags      系统管理
// @Summary   菜单分页列表
// @Description 获取菜单分页列表
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Param     data  query      request.SystemMenuPageReq      true  "菜单列表信息"
// @Produce   json
// @Success   200  {object}  code.Response{code=number,data=response.SystemMenuListResp,message=string}  "50000,success"
// @Router    /v1/system/menu/page [get]
func (b *SystemV1ApiGroup) SystemMenuPage(c *gin.Context) {

	var req request.SystemMenuPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	resp, err := systemsvc.MenuPage(c.Request.Context(), &req)
	if err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(resp, c)
}

// SystemMenuList
// @Tags      系统管理
// @Summary   全部菜单列表
// @Description 获取全部菜单列表
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Param     data  query      request.SystemMenuListReq      true  "菜单列表信息"
// @Produce   json
// @Success   200  {object}  code.Response{code=number,data=[]response.SystemMenuResp,message=string}  "50000,success"
// @Router    /v1/system/menu/list [get]
func (b *SystemV1ApiGroup) SystemMenuList(c *gin.Context) {

	var req request.SystemMenuListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	resp, err := systemsvc.MenuList(c.Request.Context(), &req)
	if err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(resp, c)
}

// SystemMenuTree
// @Tags      系统管理
// @Summary   菜单树形列表
// @Description 获取菜单树形列表
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Produce   json
// @Success   200  {object}  code.Response{code=number,data=[]response.SystemMenuTreeResp,message=string}  "50000,success"
// @Router    /v1/system/menu/tree [get]
func (b *SystemV1ApiGroup) SystemMenuTree(c *gin.Context) {

	resp, err := systemsvc.MenuTree(c.Request.Context())
	if err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(resp, c)
}

// SystemMenuLabel
// @Tags      系统管理
// @Summary   菜单标签
// @Description 获取菜单标签
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Produce   json
// @Success   200  {object}  code.Response{code=number,data=[]response.Options,message=string}  "50000,success"
// @Router    /v1/system/menu/label [get]
func (b *SystemV1ApiGroup) SystemMenuLabel(c *gin.Context) {

	resp, err := systemsvc.MenuLabel(c.Request.Context())
	if err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(resp, c)
}

// SystemMenuGetById
// @Tags      系统管理
// @Summary   获取菜单详情
// @Description 获取菜单详情
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Param     data  path      request.IdReq      true  "菜单ID"
// @Produce   json
// @Success   200  {object}  code.Response{code=number,data=response.SystemMenuResp,message=string}  "50000,success"
// @Router    /v1/system/menu/{id} [get]
func (b *SystemV1ApiGroup) SystemMenuGetById(c *gin.Context) {

	var id request.IdReq
	if err := c.ShouldBindUri(&id); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	resp, err := systemsvc.MenuGetById(c.Request.Context(), id.ID)
	if err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(resp, c)
}

// SystemMenuEnable
// @Tags      系统管理
// @Summary   启停菜单
// @Description 启停菜单状态
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Param     data  path      request.IdReq      true  "菜单ID"
// @Param     data  body      request.EnableReq      true  "状态 [1: 启用, 2: 禁用]"
// @Produce   json
// @Success   200  {object}  code.Response{code=number,message=string}  "50000,success"
// @Router    /v1/system/menu/enable/{id} [put]
func (b *SystemV1ApiGroup) SystemMenuEnable(c *gin.Context) {

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

	if err := systemsvc.MenuEnable(c.Request.Context(), id.ID, &req); err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(nil, c)
}

// SystemRoleMenus
// @Tags      系统管理
// @Summary   获取角色菜单
// @Description 获取角色绑定的菜单列表
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Param     data  path      request.IdReq      true  "角色ID"
// @Produce   json
// @Success   200  {object}  code.Response{code=number,data=[]response.SystemRoleMenuResp,message=string}  "50000,success"
// @Router    /v1/system/role/{id}/menus [get]
func (b *SystemV1ApiGroup) SystemRoleMenus(c *gin.Context) {

	var id request.IdReq
	if err := c.ShouldBindUri(&id); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	resp, err := systemsvc.GetRoleMenus(c.Request.Context(), id.ID)
	if err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(resp, c)
}

// SystemRoleBindMenus
// @Tags      系统管理
// @Summary   绑定角色菜单
// @Description 为角色绑定菜单
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Param     data  path      request.IdReq      true  "角色ID"
// @Param     data  body      request.SystemRoleMenuBindReq      true  "菜单ID列表"
// @Produce   json
// @Success   200  {object}  code.Response{code=number,message=string}  "50000,success"
// @Router    /v1/system/role/{id}/menus [put]
func (b *SystemV1ApiGroup) SystemRoleBindMenus(c *gin.Context) {

	var id request.IdReq
	if err := c.ShouldBindUri(&id); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	var req request.SystemRoleMenuBindReq
	if err := c.ShouldBindJSON(&req); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	if err := systemsvc.BindRoleMenus(c.Request.Context(), id.ID, &req); err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(nil, c)
}

// SystemRoleAddMenu
// @Tags      系统管理
// @Summary   添加菜单到角色
// @Description 为角色添加单个菜单
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Param     data  path      request.IdReq      true  "角色ID"
// @Param     data  path      request.MenuIdReq      true  "菜单ID"
// @Produce   json
// @Success   200  {object}  code.Response{code=number,message=string}  "50000,success"
// @Router    /v1/system/role/{role_id}/menu/{menu_id} [post]
func (b *SystemV1ApiGroup) SystemRoleAddMenu(c *gin.Context) {

	var req request.MenuIdReq
	if err := c.ShouldBindUri(&req); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	if err := systemsvc.AddRoleMenu(c.Request.Context(), req.RoleID, req.MenuID); err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(nil, c)
}

// SystemRoleRemoveMenu
// @Tags      系统管理
// @Summary   从角色移除菜单
// @Description 从角色移除单个菜单
// @securityDefinitions.apikey ApiKeyAuth
// @In        header
// @Name      x-quebec-token
// @Param     data  path      request.IdReq      true  "角色ID"
// @Param     data  path      request.MenuIdReq      true  "菜单ID"
// @Produce   json
// @Success   200  {object}  code.Response{code=number,message=string}  "50000,success"
// @Router    /v1/system/role/{role_id}/menu/{menu_id} [delete]
func (b *SystemV1ApiGroup) SystemRoleRemoveMenu(c *gin.Context) {

	var req request.MenuIdReq
	if err := c.ShouldBindUri(&req); err != nil {
		code.InvalidParams.Failed(c)
		return
	}

	if err := systemsvc.RemoveRoleMenu(c.Request.Context(), req.RoleID, req.MenuID); err != nil {
		err.(*code.Response).Failed(c)
		return
	}

	code.Success.Success(nil, c)
}
