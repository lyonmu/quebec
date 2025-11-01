package global

import (
	"github.com/lyonmu/quebec/cmd/gateway/internal/config"
	"go.uber.org/zap"
)

var (
	Cfg    config.Config
	Logger *zap.Logger
)
