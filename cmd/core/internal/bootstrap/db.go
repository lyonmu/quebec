package bootstrap

import (
	"context"
	"strings"

	"entgo.io/ent/dialect"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lyonmu/quebec/cmd/core/internal/ent"
	"github.com/lyonmu/quebec/cmd/core/internal/ent/migrate"
	_ "github.com/lyonmu/quebec/cmd/core/internal/ent/runtime"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/pkg/config"
	"github.com/lyonmu/quebec/pkg/tools"
)

func InitMySQL(cfg config.MySQLConfig) error {
	// 确保数据库存在
	if err := cfg.EnsureDatabase(); err != nil {
		global.Logger.Sugar().Error("ensure database failed: %v", err)
		return err
	}

	options := make([]ent.Option, 0)
	if strings.ToLower(gin.Mode()) == gin.DebugMode {
		options = append(options, ent.Debug())
	}

	client, openErr := ent.Open(dialect.MySQL, cfg.DSN(), options...)
	if openErr != nil {
		global.Logger.Sugar().Error("init mysql client failed:: %v", openErr)
		return openErr
	}
	global.EntClient = client

	// 注册全局时间戳 hooks
	tools.RegisterTimeHooks(client)

	if createErr := client.Schema.Create(context.Background(), migrate.WithForeignKeys(false)); createErr != nil {
		global.Logger.Sugar().Error("failed to create mysql schema: %v", createErr)
		return createErr
	}

	return nil
}
