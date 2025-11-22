package code

var (
	CaptchaGeneratorFailed = Response{Code: 51001, Message: "验证码生成失败"}
	UserNotExists          = Response{Code: 51001, Message: "用户不存在"}
	UserPasswordError      = Response{Code: 51002, Message: "密码错误"}
	InvalidCaptcha         = Response{Code: 51003, Message: "无效的验证码"}
)
