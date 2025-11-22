package request

type SystemLoginRequest struct {
	Username  string `json:"username" binding:"required" form:"username"`     // 用户名
	Password  string `json:"password" binding:"required" form:"password"`     // 密码
	Captcha   string `json:"captcha" binding:"required" form:"captcha"`       // 验证码
	CaptchaId string `json:"captcha_id" binding:"required" form:"captcha_id"` // 验证码id
}
