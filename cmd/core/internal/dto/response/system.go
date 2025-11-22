package response

type CaptchaResponse struct {
	Id       string `json:"id" example:"1234567890"`   // 验证码id
	Pictures string `json:"pictures" example:"base64"` // 验证码图片
	Length   int    `json:"length" example:"4"`        // 验证码长度
}

type SystemInfoResponse struct {
	Username string `json:"username"` // 用户名
	Token    string `json:"token"`    // token
}
