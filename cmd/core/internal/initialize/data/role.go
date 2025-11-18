package data

import (
	"context"
	"fmt"

	"github.com/lyonmu/quebec/cmd/core/internal/ent"
	"github.com/lyonmu/quebec/cmd/core/internal/ent/corerole"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/pkg/common"
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
			client.CoreRole.Create().SetID("1").SetName("system").SetRemark("系统管理员").SetStatus(int8(common.Yes)),
			client.CoreRole.Create().SetID("2").SetName("operations").SetRemark("运营管理员").SetStatus(int8(common.Yes)),
			client.CoreRole.Create().SetID("3").SetName("guest").SetRemark("访客").SetStatus(int8(common.Yes)),
			client.CoreRole.Create().SetID("4").SetName("user").SetRemark("普通用户").SetStatus(int8(common.Yes)),
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
