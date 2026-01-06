package response

import (
	corecommon "github.com/lyonmu/quebec/cmd/core/internal/common"
	"github.com/lyonmu/quebec/cmd/core/internal/ent"
	"github.com/lyonmu/quebec/pkg/constant"
)

type SystemMenuResp struct {
	ID             string              `json:"id,omitempty"`               // 菜单ID
	Name           string              `json:"name,omitempty"`             // 菜单名称
	MenuCode       string              `json:"menu_code,omitempty"`        // 菜单编码
	MenuType       corecommon.MenuType `json:"menu_type,omitempty"`        // 菜单类型 [1: 目录, 2: 菜单, 3: 按钮]
	ApiPath        string              `json:"api_path,omitempty"`         // 菜单API路径
	ApiPathMethod  string              `json:"api_path_method,omitempty"`  // 菜单API方法
	Order          int8                `json:"order,omitempty"`            // 菜单排序
	ParentMenuCode string              `json:"parent_menu_code,omitempty"` // 父菜单编码
	ParentName     string              `json:"parent_name,omitempty"`      // 父菜单名称
	Status         constant.YesOrNo    `json:"status,omitempty"`           // 菜单状态 [1: 启用, 2: 禁用]
}

func (r *SystemMenuResp) LoadDb(e *ent.CoreMenu) {
	r.ID = e.ID
	r.Name = e.Name
	r.MenuCode = e.MenuCode
	r.MenuType = e.MenuType
	r.ApiPath = e.APIPath
	r.ApiPathMethod = e.APIPathMethod
	r.Order = e.Order
	r.ParentMenuCode = e.ParentMenuCode
	r.Status = e.Status

	if e.Edges.MenuFromParent != nil {
		r.ParentName = e.Edges.MenuFromParent.Name
	}
}

type SystemMenuListResp struct {
	Total    int               `json:"total,omitempty"`     // 总条数
	Items    []*SystemMenuResp `json:"items,omitempty"`     // 菜单列表
	Page     int               `json:"page,omitempty"`      // 页码
	PageSize int               `json:"page_size,omitempty"` // 每页条数
}

type SystemMenuTreeResp struct {
	ID            string                `json:"id"`                        // 菜单ID
	Name          string                `json:"name"`                      // 菜单名称
	MenuCode      string                `json:"menu_code,omitempty"`       // 菜单编码
	MenuType      corecommon.MenuType   `json:"menu_type"`                 // 菜单类型
	ApiPath       string                `json:"api_path,omitempty"`        // 菜单API路径
	ApiPathMethod string                `json:"api_path_method,omitempty"` // 菜单API方法
	Order         int8                  `json:"order"`                     // 菜单排序
	Status        constant.YesOrNo      `json:"status"`                    // 菜单状态
	Children      []*SystemMenuTreeResp `json:"children,omitempty"`        // 子菜单
}

type SystemRoleMenuResp struct {
	MenuID   string `json:"menu_id"`   // 菜单ID
	MenuName string `json:"menu_name"` // 菜单名称
}
