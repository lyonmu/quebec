package main

import (
	"os"

	"github.com/alecthomas/kingpin/v2"
	authcommon "github.com/lyonmu/quebec/app/auth/common"
	authgrpc "github.com/lyonmu/quebec/app/auth/grpc"
	"github.com/lyonmu/quebec/app/bootstrap"
	"github.com/lyonmu/quebec/app/lib"
	"github.com/lyonmu/quebec/app/logger"
	"github.com/sirupsen/logrus"

	commonvariable "github.com/lyonmu/quebec/app/common/variable"
)

var (
	config = kingpin.Flag("config", "config file path").Short('c').Default("config.yml").String()
)

func init() {
	kingpin.Parse()

	if !lib.PathExists(*config) {
		os.Exit(1)
	}

	if err := lib.LoadConfig(*config, &authcommon.System); err != nil {
		os.Exit(1)
	}

	logConfig := authcommon.System.Auth.Log

	logger.SetLogger(&logConfig)
	logrus.Info(authcommon.System)

	commonvariable.Redis = bootstrap.InitRedis(authcommon.System.Redis)
}

func main() {
	authgrpc.StartGRPCServer(authcommon.System.Auth.Port)
}
