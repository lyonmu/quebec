package init

import (
	"context"
	"fmt"
	"strings"

	"entgo.io/ent/dialect"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lyonmu/quebec/cmd/core/internal/ent"
	"github.com/lyonmu/quebec/cmd/core/internal/ent/migrate"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/pkg/config"
)

func InitMySQL(cfg config.MySQLConfig) error {
	options := make([]ent.Option, 0)
	if strings.ToLower(gin.Mode()) == gin.DebugMode {
		options = append(options, ent.Debug())
	}

	client, openErr := ent.Open(dialect.MySQL, cfg.DSN(), options...)
	if openErr != nil {
		return fmt.Errorf("init mysql client failed:: %v", openErr)
	}
	global.EntClient = client

	if createErr := client.Schema.Create(context.Background(), migrate.WithForeignKeys(false)); createErr != nil {
		return fmt.Errorf("failed to create mysql schema: %v", createErr)
	}

	return nil
}
