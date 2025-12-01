package global

import (
	"github.com/lyonmu/quebec/cmd/gateway/internal/config"
	"github.com/lyonmu/quebec/pkg/tools"
	"go.uber.org/zap"
)

var (
	Cfg    config.Config
	Logger *zap.Logger
	Id     tools.IDGenerator
)
