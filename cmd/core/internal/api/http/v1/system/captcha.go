package system

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lyonmu/quebec/cmd/core/internal/dto/response"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/pkg/code"
)

// Captcha
// @Tags      系统管理
// @Summary   获取验证码
// @Description 获取验证码,返回包括随机数id,base64,验证码长度
// @Produce      json
// @Success   200  {object}  code.Response{code=number,data=response.CaptchaResponse,message=string}  "50000,success"
// @Router    /v1/system/captcha [get]
func (b *SystemV1ApiGroup) SystemCaptcha(c *gin.Context) {
	// 获取客户端标识
	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// 构造限流key
	rateLimitKey := fmt.Sprintf("captcha_limit:%s:%s", clientIP, userAgent)

	// 获取配置参数（默认与原有逻辑兼容）
	maxRequests := global.Cfg.Core.Captcha.RateLimit
	if maxRequests <= 0 {
		maxRequests = 5 // 默认值
	}

	windowSeconds := global.Cfg.Core.Captcha.RateWindow
	if windowSeconds <= 0 {
		windowSeconds = global.Cfg.Core.Captcha.Cache // 默认与缓存时间一致
	}

	// 获取当前计数
	count, err := global.RedisCli.Get(context.Background(), rateLimitKey).Int()
	if err == nil && count >= maxRequests {
		code.CaptchaTooFrequent.Failed(c)
		return
	}

	// 原子计数增加
	incrResult, err := global.RedisCli.Incr(context.Background(), rateLimitKey).Result()
	if err != nil {
		global.Logger.Sugar().Errorw("验证码限流失败", "error", err)
		code.CaptchaGeneratorFailed.Failed(c)
		return
	}

	// 首次计数设置过期时间
	if incrResult == 1 {
		global.RedisCli.Expire(context.Background(), rateLimitKey,
			time.Duration(windowSeconds)*time.Second)
	}
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
