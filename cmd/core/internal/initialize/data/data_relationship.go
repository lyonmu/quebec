package data

import (
	"context"
	"fmt"

	"github.com/lyonmu/quebec/cmd/core/internal/ent"
	"github.com/lyonmu/quebec/cmd/core/internal/ent/coredatarelationship"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
)

func InitDataRelationship(client *ent.Client) error {
	ctx := context.Background()

	exists, err := client.CoreDataRelationship.Query().Where(coredatarelationship.DeletedAtIsNil()).Exist(ctx)
	if err != nil {
		global.Logger.Error("core_data_relationship表数据初始化失败")
		return err
	}
	if exists {
		global.Logger.Info("core_data_relationship表数据无需初始化")
	} else {
		relationships, err := client.CoreDataRelationship.CreateBulk(
		// client.CoreDataRelationship.Create().SetID("1").SetDataRelationshipType(common.DataRelationshipTypeRoleToMenu).SetMenuID("1").SetRoleID("1"),
		// client.CoreDataRelationship.Create().SetID("2").SetDataRelationshipType(common.DataRelationshipTypeUserToRole).SetUserID("1").SetRoleID("1"),
		// client.CoreDataRelationship.Create().SetID("3").SetDataRelationshipType(common.DataRelationshipTypeManyToMany).SetMenuID("1").SetRoleID("1"),
		).Save(ctx)
		if err != nil {
			return err
		}
		if relationships == nil {
			return fmt.Errorf("core_data_relationship表数据初始化失败")
		}
		global.Logger.Info("core_data_relationship表数据初始化成功")
	}
	return nil
}
