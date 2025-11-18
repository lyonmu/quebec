package system

import (
	"github.com/gin-gonic/gin"
	"github.com/lyonmu/quebec/cmd/core/internal/dto/response"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/pkg/code"
)

// Captcha
// @Tags      系统管理
// @Summary   获取验证码
// @Description 获取验证码,返回包括随机数id,base64,验证码长度,是否开启验证码
// @Produce      json
// @Success   200  {object}  code.Response{code=number,data=response.CaptchaResponse,message=string}  "50000,success"
// @Router    /system/captcha [get]
func (b *SystemV1ApiGroup) SystemCaptcha(c *gin.Context) {

	id, b64s, _, err := global.CaptchaGenerator.Generate()
	if err != nil {
		global.Logger.Sugar().Error("验证码获取失败!", err)
		code.CaptchaGeneratorFailed.Failed(c)
		return
	}
	captchaResponse := response.CaptchaResponse{
		Id:       id,
		Pictures: b64s,
		Length:   global.Cfg.Core.Captcha.Length,
	}
	code.Success.Success(captchaResponse, c)
}
