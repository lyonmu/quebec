package config

import (
	"context"
	"fmt"
	"time"

	"github.com/lyonmu/quebec/cmd/core/internal/common"
	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
)

type Captcha struct {
	Length int `name:"length" env:"CAPTCHA_LENGTH" default:"4" help:"验证码长度" mapstructure:"length" yaml:"length" json:"length"`
	Width  int `name:"width" env:"CAPTCHA_WIDTH" default:"240" help:"验证码宽度" mapstructure:"width" yaml:"width" json:"width"`
	Height int `name:"height" env:"CAPTCHA_HEIGHT" default:"80" help:"验证码高度" mapstructure:"height" yaml:"height" json:"height"`
	Cache  int `name:"cache" env:"CAPTCHA_CACHE" default:"300" help:"验证码缓存时间" mapstructure:"cache" yaml:"cache" json:"cache"`
}

type CaptchaWithRedis struct {
	Captcha
	Redis redis.UniversalClient
}

func (s *CaptchaWithRedis) Set(id string, value string) error {
	return s.Redis.Set(context.Background(), fmt.Sprintf(common.CaptchaCache, id), value, time.Duration(s.Cache)*time.Second).Err()
}

func (s *CaptchaWithRedis) Verify(id, answer string, clear bool) bool {
	value := s.Redis.Get(context.Background(), fmt.Sprintf(common.CaptchaCache, id)).Val()
	s.Redis.Del(context.Background(), fmt.Sprintf(common.CaptchaCache, id))
	return value == answer
}

func (s *CaptchaWithRedis) Get(id string, clear bool) (value string) {
	value = s.Redis.Get(context.Background(), fmt.Sprintf(common.CaptchaCache, id)).Val()
	if clear {
		s.Redis.Del(context.Background(), fmt.Sprintf(common.CaptchaCache, id))
	}
	return
}

func (c Captcha) WithRedis(redis redis.UniversalClient) *base64Captcha.Captcha {
	return base64Captcha.NewCaptcha(
		base64Captcha.NewDriverDigit(c.Height, c.Width, c.Length, 0, 1),
		&CaptchaWithRedis{Captcha: c, Redis: redis},
	)
}
