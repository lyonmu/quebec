package auth

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

func main() {
	if !lib.PathExists(*config) {
		os.Exit(1)
	}
	
	if err := lib.LoadConfig(*config, &authcommon.Config); err != nil {
		os.Exit(1)
	}
	
	logConfig := authcommon.Config.Log
	
	logger.SetLogger(&logConfig)
	logrus.Info(authcommon.Config)
	
	commonvariable.Redis = bootstrap.InitRedis(authcommon.Config.Redis)
	
	authgrpc.StartGRPCServer(authcommon.Config.Port)
	
}
