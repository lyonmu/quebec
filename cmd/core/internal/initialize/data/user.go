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

func InitUser(client *ent.Client, roles []*ent.CoreRole) error {

	var system, operator, user = "", "", ""
	systemPassword, _ := encrypt.HashWithBcryptBytes(encrypt.HashWithSHA256Bytes([]byte("system@123456")))
	operatorPassword, _ := encrypt.HashWithBcryptBytes(encrypt.HashWithSHA256Bytes([]byte("operator@123456")))
	userPassword, _ := encrypt.HashWithBcryptBytes(encrypt.HashWithSHA256Bytes([]byte("user@123456")))

	for _, role := range roles {
		switch role.Name {
		case "系统管理员":
			system = role.ID
		case "运维操作员":
			operator = role.ID
		case "普通用户":
			user = role.ID
		}
	}

	ctx := context.Background()

	exists, err := client.CoreUser.Query().Where(coreuser.DeletedAtIsNil()).Exist(ctx)
	if err != nil {
		global.Logger.Sugar().Error("core_user表数据初始化失败")
		return err
	}
	if exists {
		global.Logger.Sugar().Info("core_user表数据无需初始化")
	} else {
		users, err := client.CoreUser.CreateBulk(
			client.CoreUser.Create().
				SetID(fmt.Sprintf("%d", global.Id.GenID())).
				SetUsername(encrypt.HashWithSHA256String("system")).
				SetPassword(string(systemPassword)).
				SetNickname("系统管理员").
				SetStatus(constant.Yes).
				SetRemark("system").
				SetRoleID(system).
				SetSystem(constant.Yes),
			client.CoreUser.Create().
				SetID(fmt.Sprintf("%d", global.Id.GenID())).
				SetUsername(encrypt.HashWithSHA256String("operator")).
				SetPassword(string(operatorPassword)).
				SetNickname("运维操作员").
				SetStatus(constant.Yes).
				SetRemark("operator").
				SetRoleID(operator).
				SetSystem(constant.Yes),
			client.CoreUser.Create().
				SetID(fmt.Sprintf("%d", global.Id.GenID())).
				SetUsername(encrypt.HashWithSHA256String("user")).
				SetPassword(string(userPassword)).
				SetNickname("普通用户").
				SetStatus(constant.Yes).
				SetRemark("user").
				SetRoleID(user).
				SetSystem(constant.Yes),
		).Save(ctx)
		if err != nil {
			return err
		}
		if users == nil {
			return fmt.Errorf("core_user表数据初始化失败")
		}
		global.Logger.Sugar().Info("core_user表数据初始化成功")
	}
	return nil
}
