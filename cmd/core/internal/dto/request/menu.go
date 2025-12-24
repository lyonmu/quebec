package request

import corecommon "github.com/lyonmu/quebec/cmd/core/internal/common"
import "github.com/lyonmu/quebec/pkg/constant"

type SystemMenuPageReq struct {
	Name       string             `json:"name,omitempty" form:"name"`                           // 菜单名称
	MenuType   corecommon.MenuType `json:"menu_type,omitempty" form:"menu_type"`               // 菜单类型 [1: 目录, 2: 菜单, 3: 按钮]
	Status     constant.YesOrNo    `json:"status,omitempty" form:"status"`                       // 菜单状态 [1: 启用, 2: 禁用]
	ParentID   string             `json:"parent_id,omitempty" form:"parent_id"`                 // 父菜单ID
	Page       int                `json:"page,omitempty" binding:"required,min=1" form:"page" minimum:"1" default:"1"`                                      // 页码
	PageSize   int                `json:"page_size,omitempty" binding:"required,min=10,max=1000" form:"page_size" minimum:"10" maximum:"1000" default:"10"` // 每页条数
}

type SystemMenuListReq struct {
	Name     string             `json:"name,omitempty" form:"name"`         // 菜单名称
	MenuType corecommon.MenuType `json:"menu_type,omitempty" form:"menu_type"` // 菜单类型
	Status   constant.YesOrNo    `json:"status,omitempty" form:"status"`     // 菜单状态
}

type SystemMenuAddReq struct {
	Name          string             `json:"name,omitempty" binding:"required" form:"name"`                       // 菜单名称
	MenuType      corecommon.MenuType `json:"menu_type,omitempty" form:"menu_type"`                               // 菜单类型 [1: 目录, 2: 菜单, 3: 按钮]
	ApiPath       *string            `json:"api_path,omitempty" form:"api_path"`                                   // 菜单API路径
	ApiPathMethod *string            `json:"api_path_method,omitempty" form:"api_path_method"`                     // 菜单API方法
	Order         *int8              `json:"order,omitempty" form:"order"`                                         // 菜单排序
	ParentID      *string            `json:"parent_id,omitempty" form:"parent_id"`                                 // 父菜单ID
	Component     *string            `json:"component,omitempty" form:"component"`                                 // 菜单组件
	Status        *constant.YesOrNo   `json:"status,omitempty" form:"status"`                                      // 菜单状态 [1: 启用, 2: 禁用]
	Remark        *string            `json:"remark,omitempty" form:"remark"`                                       // 菜单备注
}

type SystemMenuUpdateReq struct {
	Name          *string              `json:"name,omitempty" form:"name"`             // 菜单名称
	MenuType      *corecommon.MenuType  `json:"menu_type,omitempty" form:"menu_type"` // 菜单类型
	ApiPath       *string              `json:"api_path,omitempty" form:"api_path"`     // 菜单API路径
	ApiPathMethod *string              `json:"api_path_method,omitempty" form:"api_path_method"` // 菜单API方法
	Order         *int8                `json:"order,omitempty" form:"order"`           // 菜单排序
	ParentID      *string              `json:"parent_id,omitempty" form:"parent_id"`   // 父菜单ID
	Component     *string              `json:"component,omitempty" form:"component"`   // 菜单组件
	Status        *constant.YesOrNo    `json:"status,omitempty" form:"status"`         // 菜单状态
	Remark        *string              `json:"remark,omitempty" form:"remark"`         // 菜单备注
}

type SystemRoleMenuBindReq struct {
	MenuIDs []string `json:"menu_ids" binding:"required"` // 菜单ID列表
}
