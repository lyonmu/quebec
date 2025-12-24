package request

import "github.com/lyonmu/quebec/pkg/constant"

type IdReq struct {
	ID string `json:"id" binding:"required" form:"id" uri:"id"` // ID
}

type MenuIdReq struct {
	RoleID string `json:"role_id" binding:"required" uri:"role_id"`     // 角色ID
	MenuID string `json:"menu_id" binding:"required" uri:"menu_id"`     // 菜单ID
}

type EnableReq struct {
	Status constant.YesOrNo `json:"status" binding:"required,min=1,max=2" minimum:"1" maximum:"2" form:"status"` // 状态 [1: 启用, 2: 禁用]
}
