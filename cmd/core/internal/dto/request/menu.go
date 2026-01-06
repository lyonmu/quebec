package request

import (
	corecommon "github.com/lyonmu/quebec/cmd/core/internal/common"
	"github.com/lyonmu/quebec/pkg/constant"
)

type SystemMenuPageReq struct {
	Name           string              `json:"name,omitempty" form:"name"`                                                                                       // 菜单名称
	MenuCode       string              `json:"menu_code,omitempty" form:"menu_code"`                                                                             // 菜单编码
	MenuType       corecommon.MenuType `json:"menu_type,omitempty" form:"menu_type"`                                                                             // 菜单类型 [1: 目录, 2: 菜单, 3: 按钮]
	Status         constant.YesOrNo    `json:"status,omitempty" form:"status"`                                                                                   // 菜单状态 [1: 启用, 2: 禁用]
	ParentMenuCode string              `json:"parent_menu_code,omitempty" form:"parent_menu_code"`                                                               // 父菜单编码
	Page           int                 `json:"page,omitempty" binding:"required,min=1" form:"page" minimum:"1" default:"1"`                                      // 页码
	PageSize       int                 `json:"page_size,omitempty" binding:"required,min=10,max=1000" form:"page_size" minimum:"10" maximum:"1000" default:"10"` // 每页条数
}

type SystemMenuListReq struct {
	Name           string              `json:"name,omitempty" form:"name"`                         // 菜单名称
	MenuCode       string              `json:"menu_code,omitempty" form:"menu_code"`               // 菜单编码
	MenuType       corecommon.MenuType `json:"menu_type,omitempty" form:"menu_type"`               // 菜单类型 [1: 目录, 2: 菜单, 3: 按钮]
	Status         constant.YesOrNo    `json:"status,omitempty" form:"status"`                     // 菜单状态 [1: 启用, 2: 禁用]
	ParentMenuCode string              `json:"parent_menu_code,omitempty" form:"parent_menu_code"` // 父菜单编码
}

type SystemRoleMenuBindReq struct {
	MenuIDs []string `json:"menu_ids" binding:"required"` // 菜单ID列表
}
