package bootstrap

import (
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/lyonmu/quebec/cmd/gateway/internal/global"
	"github.com/lyonmu/quebec/pkg/logger"
	"github.com/lyonmu/quebec/pkg/metrics"
	"github.com/lyonmu/quebec/pkg/mq/serializer"
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
		global.Metrics = metrics.NewPrometheusRegistry()

		// 初始化 ALS Kafka Producer
		// Key 使用 String 编码，Value 使用 Binary 编码
		if len(global.Cfg.Gateway.AlsTopic) > 0 {
			producer, err := global.Cfg.Kafka.Producer(
				global.Cfg.Gateway.AlsTopic,
				serializer.NewStringCodec[string, string](),
				serializer.NewBinaryCodec[[]byte, []byte](),
			)
			if err != nil {
				global.Logger.Sugar().Errorf("failed to create ALS Kafka producer: %v", err)
				// 不退出，ALS 功能降级为仅日志
			} else {
				global.AlsKafkaProducer = producer
				global.Logger.Sugar().Infof("ALS Kafka producer initialized, topic: %s", global.Cfg.Gateway.AlsTopic)
			}
		}

	})

	if err := InitServer(); err != nil {
		global.Logger.Sugar().Error("初始化服务失败: %v", err)
		os.Exit(1)
	}
}
