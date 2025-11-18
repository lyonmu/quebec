package initialize

import (
	"github.com/lyonmu/quebec/cmd/core/internal/ent"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/cmd/core/internal/initialize/data"
)

func Init(client *ent.Client) error {
	if err := data.InitRole(client); err != nil {
		global.Logger.Sugar().Errorf("core_role表数据初始化失败: %v", err)
		return err
	}
	if err := data.InitUser(client); err != nil {
		global.Logger.Sugar().Errorf("core_user表数据初始化失败: %v", err)
		return err
	}
	return nil
}
