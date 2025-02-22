package variable

import (
	"github.com/redis/go-redis/v9"
	
	commonconfig "github.com/lyonmu/quebec/app/common/config"
)

var (
	Redis redis.UniversalClient
	
	Distributed commonconfig.DistributedConfig
)
