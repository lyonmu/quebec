package request

import "github.com/lyonmu/quebec/pkg/constant"

type IdReq struct {
	ID string `json:"id" binding:"required" form:"id" uri:"id"` // ID
}

type EnableReq struct {
	Status constant.YesOrNo `json:"status" binding:"required,min=1,max=2" minimum:"1" maximum:"2" form:"status"` // 状态 [1: 启用, 2: 禁用]
}
