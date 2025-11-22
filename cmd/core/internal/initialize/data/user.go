package data

import (
	"context"
	"fmt"

	"github.com/lyonmu/quebec/cmd/core/internal/ent"
	"github.com/lyonmu/quebec/cmd/core/internal/ent/coreuser"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/pkg/constant"
	"github.com/lyonmu/quebec/pkg/tools/encrypt"
)

func InitUser(client *ent.Client, role_id string) error {
	ctx := context.Background()

	password, _ := encrypt.HashWithBcryptBytes(encrypt.HashWithSHA256Bytes([]byte("Quebec@123456")))

	exists, err := client.CoreUser.Query().Where(coreuser.DeletedAtIsNil()).Exist(ctx)
	if err != nil {
		global.Logger.Error("core_user表数据初始化失败")
		return err
	}
	if exists {
		global.Logger.Info("core_user表数据无需初始化")
	} else {
		users, err := client.CoreUser.CreateBulk(
			client.CoreUser.Create().
				SetID(fmt.Sprintf("%d", global.Id.GenID())).
				SetUsername(encrypt.HashWithSHA256String("system")).
				SetPassword(string(password)).
				SetNickname("system").
				SetStatus(constant.Yes).
				SetRemark("system").
				SetRoleID(role_id),
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
