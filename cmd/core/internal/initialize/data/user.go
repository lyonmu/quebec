package data

import (
	"context"
	"fmt"

	"github.com/lyonmu/quebec/cmd/core/internal/ent"
	"github.com/lyonmu/quebec/cmd/core/internal/ent/coreuser"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/pkg/common"
)

func InitUser(client *ent.Client) error {
	ctx := context.Background()

	exists, err := client.CoreUser.Query().Where(coreuser.DeletedAtIsNil()).Exist(ctx)
	if err != nil {
		global.Logger.Error("core_user表数据初始化失败")
		return err
	}
	if exists {
		global.Logger.Info("core_user表数据无需初始化")
	} else {
		users, err := client.CoreUser.CreateBulk(
			client.CoreUser.Create().SetID("1").SetUsername("system_1").SetPassword(common.DefaultPassword).SetNickname("系统管理员1").SetStatus(int8(common.Yes)).SetRemark("系统管理员1").SetRoleID("1"),
			client.CoreUser.Create().SetID("2").SetUsername("system_2").SetPassword(common.DefaultPassword).SetNickname("系统管理员2").SetStatus(int8(common.Yes)).SetRemark("系统管理员2").SetRoleID("1"),
			client.CoreUser.Create().SetID("3").SetUsername("system_3").SetPassword(common.DefaultPassword).SetNickname("系统管理员3").SetStatus(int8(common.Yes)).SetRemark("系统管理员3").SetRoleID("1"),
			client.CoreUser.Create().SetID("4").SetUsername("operations_1").SetPassword(common.DefaultPassword).SetNickname("运营管理员1").SetStatus(int8(common.Yes)).SetRemark("运营管理员1").SetRoleID("2"),
			client.CoreUser.Create().SetID("5").SetUsername("operations_2").SetPassword(common.DefaultPassword).SetNickname("运营管理员2").SetStatus(int8(common.Yes)).SetRemark("运营管理员2").SetRoleID("2"),
			client.CoreUser.Create().SetID("6").SetUsername("operations_3").SetPassword(common.DefaultPassword).SetNickname("运营管理员3").SetStatus(int8(common.Yes)).SetRemark("运营管理员3").SetRoleID("2"),
			client.CoreUser.Create().SetID("7").SetUsername("guest_1").SetPassword(common.DefaultPassword).SetNickname("访客1").SetStatus(int8(common.Yes)).SetRemark("访客1").SetRoleID("3"),
			client.CoreUser.Create().SetID("8").SetUsername("guest_2").SetPassword(common.DefaultPassword).SetNickname("访客2").SetStatus(int8(common.Yes)).SetRemark("访客2").SetRoleID("3"),
			client.CoreUser.Create().SetID("9").SetUsername("guest_3").SetPassword(common.DefaultPassword).SetNickname("访客3").SetStatus(int8(common.Yes)).SetRemark("访客3").SetRoleID("3"),
			client.CoreUser.Create().SetID("10").SetUsername("user_1").SetPassword(common.DefaultPassword).SetNickname("普通用户1").SetStatus(int8(common.Yes)).SetRemark("普通用户1").SetRoleID("4"),
			client.CoreUser.Create().SetID("11").SetUsername("user_2").SetPassword(common.DefaultPassword).SetNickname("普通用户2").SetStatus(int8(common.Yes)).SetRemark("普通用户2").SetRoleID("4"),
			client.CoreUser.Create().SetID("12").SetUsername("user_3").SetPassword(common.DefaultPassword).SetNickname("普通用户3").SetStatus(int8(common.Yes)).SetRemark("普通用户3").SetRoleID("4"),
		).Save(ctx)
		if err != nil {
			return err
		}
		if users == nil {
			return fmt.Errorf("core_user表数据初始化失败")
		}
		global.Logger.Info("core_user表数据初始化成功")
	}
	return nil
}
