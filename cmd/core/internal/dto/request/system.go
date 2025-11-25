package request

type SystemLoginRequest struct {
	Username  string `json:"username" binding:"required" form:"username"`     // 用户名
	Password  string `json:"password" binding:"required" form:"password"`     // 密码
	Captcha   string `json:"captcha" binding:"required" form:"captcha"`       // 验证码
	CaptchaId string `json:"captcha_id" binding:"required" form:"captcha_id"` // 验证码id
}

type SystemOnlineUserListReq struct {
	UserID    string `json:"user_id"  form:"user_id"`                                                                              // 用户ID
	AccessIP  string `json:"access_ip"  form:"access_ip"`                                                                          // 访问IP
	StartTime int64  `json:"start_time" form:"start_time"`                                                                         // 开始时间
	EndTime   int64  `json:"end_time" form:"end_time"`                                                                             // 结束时间
	Page      int    `json:"page" binding:"required,min=1" form:"page" minimum:"1" default:"1"`                                    // 页码
	PageSize  int    `json:"page_size" binding:"required,min=10,max=100" form:"page_size" minimum:"10" maximum:"100" default:"10"` // 每页条数
}
