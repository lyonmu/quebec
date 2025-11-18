package bootstrap

import (
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/cmd/core/internal/initialize"
	"github.com/lyonmu/quebec/pkg/logger"
	"github.com/lyonmu/quebec/pkg/tools"
	"gopkg.in/yaml.v3"
)

var once sync.Once

func Start() {

	once.Do(func() {
		global.Logger = logger.NewZapLogger(global.Cfg.Log)
		global.Logger.Sugar().Info("服务初始化开始")

		bytes, _ := yaml.Marshal(global.Cfg)
		global.Logger.Sugar().Infof("==================== launching quebec core with config ====================\n%s", string(bytes))

		if strings.ToLower(global.Cfg.Log.Level) != "debug" {
			gin.SetMode(gin.ReleaseMode)
		}

		idGenerator, err := tools.NewSonySnowFlake(global.Cfg.MachineID)
		if err != nil {
			global.Logger.Sugar().Error("初始化雪花算法失败: %v", err)
			os.Exit(1)
		}
		global.Id = idGenerator

		if err := InitMySQL(global.Cfg.MySQL); err != nil {
			global.Logger.Sugar().Error("初始化mysql失败: %v", err)
			os.Exit(1)
		}

		global.Redis = global.Cfg.Redis.Client(global.Cfg.Log.Module)

		global.CaptchaGenerator = global.Cfg.Core.Captcha.WithRedis(global.Redis)

		if err := initialize.Init(global.EntClient); err != nil {
			global.Logger.Sugar().Error("初始化数据失败: %v", err)
			os.Exit(1)
		}

		if err := InitServer(); err != nil {
			global.Logger.Sugar().Error("初始化服务失败: %v", err)
			os.Exit(1)
		}
	})

}
