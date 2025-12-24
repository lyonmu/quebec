package request

import (
	"github.com/lyonmu/quebec/cmd/core/internal/common"
	"github.com/lyonmu/quebec/pkg/constant"
)

type SystemLoginRequest struct {
	Username  string `json:"username,omitempty" binding:"required" form:"username"`     // 用户名
	Password  string `json:"password,omitempty" binding:"required" form:"password"`     // 密码
	Captcha   string `json:"captcha,omitempty" binding:"required" form:"captcha"`       // 验证码
	CaptchaId string `json:"captcha_id,omitempty" binding:"required" form:"captcha_id"` // 验证码id
}

type SystemOnlineUserListReq struct {
	UserID    string `json:"user_id,omitempty"  form:"user_id"`                                                                                // 用户ID
	AccessIP  string `json:"access_ip,omitempty"  form:"access_ip"`                                                                            // 访问IP
	StartTime int64  `json:"start_time,omitempty" form:"start_time"`                                                                           // 开始时间
	EndTime   int64  `json:"end_time,omitempty" form:"end_time"`                                                                               // 结束时间
	Page      int    `json:"page,omitempty" binding:"required,min=1" form:"page" minimum:"1" default:"1"`                                      // 页码
	PageSize  int    `json:"page_size,omitempty" binding:"required,min=10,max=1000" form:"page_size" minimum:"10" maximum:"1000" default:"10"` // 每页条数
}

type SystemRolePageReq struct {
	Name     string           `json:"name,omitempty" form:"name"`                                                                                       // 角色名称
	Status   constant.YesOrNo `json:"status,omitempty" form:"status"`                                                                                   // 角色状态 [1: 启用, 2: 禁用]
	Page     int              `json:"page,omitempty" binding:"required,min=1" form:"page" minimum:"1" default:"1"`                                      // 页码
	PageSize int              `json:"page_size,omitempty" binding:"required,min=10,max=1000" form:"page_size" minimum:"10" maximum:"1000" default:"10"` // 每页条数
}

type SystemRoleListReq struct {
	Name   string           `json:"name,omitempty" form:"name"`     // 角色名称
	Status constant.YesOrNo `json:"status,omitempty" form:"status"` // 角色状态 [1: 启用, 2: 禁用]
}

type SystemRoleAddReq struct {
	Name   string  `json:"name,omitempty" binding:"required" form:"name"` // 角色名称
	Remark *string `json:"remark,omitempty" form:"remark"`                // 角色备注
}

type SystemRoleUpdateReq struct {
	Name   *string `json:"name,omitempty" form:"name"`     // 角色名称
	Remark *string `json:"remark,omitempty" form:"remark"` // 角色备注
}

type SystemUserPageReq struct {
	Username string           `json:"username,omitempty" form:"username"`                                                                               // 用户名名称
	Nickname string           `json:"nickname,omitempty" form:"nickname"`                                                                               // 用户昵称
	RoleID   string           `json:"role_id,omitempty" form:"role_id"`                                                                                 // 角色ID
	Email    string           `json:"email,omitempty" form:"email"`                                                                                     // 用户邮箱
	Status   constant.YesOrNo `json:"status,omitempty" form:"status"`                                                                                   // 用户状态 [1: 启用, 2: 禁用]
	Page     int              `json:"page,omitempty" binding:"required,min=1" form:"page" minimum:"1" default:"1"`                                      // 页码
	PageSize int              `json:"page_size,omitempty" binding:"required,min=10,max=1000" form:"page_size" minimum:"10" maximum:"1000" default:"10"` // 每页条数
}

type SystemUserAddReq struct {
	Username string            `json:"username,omitempty" binding:"required" form:"username"` // 用户名
	Nickname string            `json:"nickname,omitempty" binding:"required" form:"nickname"` // 用户昵称
	Password string            `json:"password,omitempty" binding:"required" form:"password"` // 密码
	RoleID   string            `json:"role_id,omitempty" binding:"required" form:"role_id"`   // 角色ID
	Email    *string           `json:"email,omitempty" binding:"required" form:"email"`       // 用户邮箱
	Remark   *string           `json:"remark,omitempty" form:"remark"`                        // 用户备注
	Status   *constant.YesOrNo `json:"status,omitempty" form:"status"`                        // 用户状态 [1: 启用, 2: 禁用]
}

type SystemUserUpdateReq struct {
	Username *string `json:"username,omitempty" form:"username"` // 用户名
	Nickname *string `json:"nickname,omitempty" form:"nickname"` // 用户昵称
	RoleID   *string `json:"role_id,omitempty" form:"role_id"`   // 角色ID
	Email    *string `json:"email,omitempty" form:"email"`       // 用户邮箱
	Remark   *string `json:"remark,omitempty" form:"remark"`     // 用户备注
}

type SystemUserListReq struct {
	Username string           `json:"username,omitempty" form:"username"` // 用户名
	Nickname string           `json:"nickname,omitempty" form:"nickname"` // 用户昵称
	RoleID   string           `json:"role_id,omitempty" form:"role_id"`   // 角色ID
	Email    string           `json:"email,omitempty" form:"email"`       // 用户邮箱
	Status   constant.YesOrNo `json:"status,omitempty" form:"status"`     // 用户状态 [1: 启用, 2: 禁用]
}

type SystemUserEditPasswordReq struct {
	PrePassword     string `json:"pre_password,omitempty" binding:"required" form:"pre_password"`         // 旧密码
	NewPassword     string `json:"new_password,omitempty" binding:"required" form:"new_password"`         // 新密码
	ConfirmPassword string `json:"confirm_password,omitempty" binding:"required" form:"confirm_password"` // 确认密码
}

type OperationLogReq struct {
	ID                   string               `json:"id,omitempty"`                     // ID
	AccessIP             string               `json:"access_ip,omitempty"`              // 访问IP
	OperationTime        int64                `json:"operation_time,omitempty"`         // 操作时间
	OperationType        common.OperationType `json:"operation_type,omitempty"`         // 操作类型
	Os                   string               `json:"os,omitempty"`                     // 操作系统
	Platform             string               `json:"platform,omitempty"`               // 操作平台
	BrowserName          string               `json:"browser_name,omitempty"`           // 浏览器名称
	BrowserVersion       string               `json:"browser_version,omitempty"`        // 浏览器版本
	BrowserEngineName    string               `json:"browser_engine_name,omitempty"`    // 浏览器引擎名称
	BrowserEngineVersion string               `json:"browser_engine_version,omitempty"` // 浏览器引擎版本
}

type OperationLogPageReq struct {
	ID            string               `json:"id,omitempty"`                                                                                                     // ID
	AccessIP      string               `json:"access_ip,omitempty"`                                                                                              // 访问IP
	StartTime     int64                `json:"start_time,omitempty" form:"start_time"`                                                                           // 开始时间
	EndTime       int64                `json:"end_time,omitempty" form:"end_time"`                                                                               // 结束时间
	OperationType common.OperationType `json:"operation_type,omitempty"`                                                                                         // 操作类型
	Page          int                  `json:"page,omitempty" binding:"required,min=1" form:"page" minimum:"1" default:"1"`                                      // 页码
	PageSize      int                  `json:"page_size,omitempty" binding:"required,min=10,max=1000" form:"page_size" minimum:"10" maximum:"1000" default:"10"` // 每页条数
}
