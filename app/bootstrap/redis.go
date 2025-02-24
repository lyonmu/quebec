package bootstrap

import (
	"context"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	commonconfig "github.com/lyonmu/quebec/app/common/config"
)

func InitRedis(cfg commonconfig.RedisConfig) redis.UniversalClient {
	options := &redis.UniversalOptions{
		Addrs:    cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: 100,
	}

	rdb := redis.NewUniversalClient(options)

	// 开启 tracing instrumentation.
	if err := redisotel.InstrumentTracing(rdb); err != nil {
		logrus.Error("start tracing instrumentation failed,err is ", err)
	}

	// 开启 metrics instrumentation.
	if err := redisotel.InstrumentMetrics(rdb); err != nil {
		logrus.Error("start metrics instrumentation failed,err is ", err)
	}

	// 检查连接
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		logrus.Error("redis connect failed: ", err)
	} else {
		logrus.Info("redis connect successfully")
	}

	return rdb
}
