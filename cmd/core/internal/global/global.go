package global

import (
	"github.com/lyonmu/quebec/cmd/core/internal/config"
	"github.com/lyonmu/quebec/cmd/core/internal/ent"
	"github.com/lyonmu/quebec/pkg/tools"
	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	Cfg              config.Config
	Logger           *zap.Logger
	EntClient        *ent.Client
	Redis            redis.UniversalClient
	Id               tools.IDGenerator
	CaptchaGenerator *base64Captcha.Captcha
)
