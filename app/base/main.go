package main

import (
	"os"

	"github.com/alecthomas/kingpin/v2"
	basecommon "github.com/lyonmu/quebec/app/base/common"
	baseinit "github.com/lyonmu/quebec/app/base/init"
	"github.com/lyonmu/quebec/app/bootstrap"
	commonvariable "github.com/lyonmu/quebec/app/common/variable"
	"github.com/lyonmu/quebec/app/lib"
	"github.com/lyonmu/quebec/app/logger"
	"github.com/sirupsen/logrus"
)

var (
	config = kingpin.Flag("config", "config file path").Short('c').Default("config.yml").String()
)

func init() {
	kingpin.Parse()
	if !lib.PathExists(*config) {
		os.Exit(1)
	}
	if err := lib.LoadConfig(*config, &basecommon.System); err != nil {
		os.Exit(1)
	}
	logConfig := basecommon.System.Base.Log
	logger.SetLogger(&logConfig)
	logrus.Info(basecommon.System)
	commonvariable.Redis = bootstrap.InitRedis(basecommon.System.Redis)
}

func main() {
	baseinit.Start()
}
