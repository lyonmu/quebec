package data

import (
	"context"
	"fmt"

	"github.com/lyonmu/quebec/cmd/core/internal/ent"
	"github.com/lyonmu/quebec/cmd/core/internal/ent/corerole"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/pkg/constant"
)

func InitRole(client *ent.Client) ([]*ent.CoreRole, error) {
	ctx := context.Background()
	var resp []*ent.CoreRole

	exists, err := client.CoreRole.Query().Where(corerole.DeletedAtIsNil()).Exist(ctx)
	if err != nil {
		global.Logger.Sugar().Error("core_role表数据初始化失败")
		return nil, err
	}

	if exists {
		global.Logger.Sugar().Info("core_role表数据无需初始化")
	} else {
		roles, err := client.CoreRole.CreateBulk(
			client.CoreRole.Create().
				SetID(fmt.Sprintf("%d", global.Id.GenID())).
				SetName("system").
				SetRemark("system").
				SetStatus(constant.Yes),
		).Save(ctx)
		if err != nil {
			return nil, err
		}
		if roles == nil {
			return nil, fmt.Errorf("roles is nil")
		}
		resp = roles
		global.Logger.Sugar().Info("core_role表数据初始化成功")
	}

	return resp, nil

}
