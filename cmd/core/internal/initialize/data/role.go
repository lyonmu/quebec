package data

import (
	"context"
	"fmt"

	"github.com/lyonmu/quebec/cmd/core/internal/ent"
	"github.com/lyonmu/quebec/cmd/core/internal/ent/corerole"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/pkg/constant"
)

func InitRole(client *ent.Client) error {
	ctx := context.Background()

	exists, err := client.CoreRole.Query().Where(corerole.DeletedAtIsNil()).Exist(ctx)
	if err != nil {
		global.Logger.Error("core_role表数据初始化失败")
		return err
	}

	if exists {
		global.Logger.Info("core_role表数据无需初始化")
	} else {
		roles, err := client.CoreRole.CreateBulk(
			client.CoreRole.Create().
				SetID("1").
				SetName("system").
				SetRemark("system").
				SetStatus(constant.Yes),
		).Save(ctx)
		if err != nil {
			return err
		}
		if roles == nil {
			return fmt.Errorf("roles is nil")
		}
		global.Logger.Info("core_role表数据初始化成功")
	}

	return nil

}
