package response

import (
	"github.com/lyonmu/quebec/cmd/core/internal/common"
	"github.com/lyonmu/quebec/cmd/core/internal/ent"
)

type CaptchaResponse struct {
	Id       string `json:"id" example:"1234567890"`   // 验证码id
	Pictures string `json:"pictures" example:"base64"` // 验证码图片
	Length   int    `json:"length" example:"4"`        // 验证码长度
}

type SystemInfoResponse struct {
	Username string `json:"username"`  // 用户名
	Token    string `json:"token"`     // token
	RoleName string `json:"role_name"` // 角色名称
}

type SystemOnlineUserResp struct {
	ID                   string               `json:"id"`                               // ID
	Nickname             string               `json:"nickname"`                         // 用户昵称
	AccessIP             string               `json:"access_ip"`                        // 访问IP
	LastOperationTime    int64                `json:"last_operation_time"`              // 最后操作时间
	OperationType        common.OperationType `json:"operation_type"`                   // 操作类型
	Os                   string               `json:"os,omitempty"`                     // 操作系统
	Platform             string               `json:"platform,omitempty"`               // 操作平台
	BrowserName          string               `json:"browser_name,omitempty"`           // 浏览器名称
	BrowserVersion       string               `json:"browser_version,omitempty"`        // 浏览器版本
	BrowserEngineName    string               `json:"browser_engine_name,omitempty"`    // 浏览器引擎名称
	BrowserEngineVersion string               `json:"browser_engine_version,omitempty"` // 浏览器引擎版本
}

func (r *SystemOnlineUserResp) LoadDb(e *ent.CoreOnLineUser) {
	r.ID = e.UserID
	r.AccessIP = e.AccessIP
	r.LastOperationTime = e.LastOperationTime
	r.OperationType = e.OperationType
	r.Os = e.Os
	r.Platform = e.Platform
	r.BrowserName = e.BrowserName
	r.BrowserVersion = e.BrowserVersion
	r.BrowserEngineName = e.BrowserEngineName
	r.BrowserEngineVersion = e.BrowserEngineVersion

	if e.Edges.OnLineFromUser != nil {
		r.Nickname = e.Edges.OnLineFromUser.Nickname
	}
}

type SystemOnlineUserListResp struct {
	Total    int                     `json:"total,omitempty"`     // 总条数
	Items    []*SystemOnlineUserResp `json:"items,omitempty"`     // 在线用户列表
	Page     int                     `json:"page,omitempty"`      // 页码
	PageSize int                     `json:"page_size,omitempty"` // 每页条数
}
