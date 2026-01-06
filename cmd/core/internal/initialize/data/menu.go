package data

import (
	"context"
	"fmt"

	"github.com/lyonmu/quebec/cmd/core/internal/common"
	"github.com/lyonmu/quebec/cmd/core/internal/ent"
	"github.com/lyonmu/quebec/cmd/core/internal/ent/coremenu"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
)

func InitMenu(client *ent.Client) error {

	var (
		menus = []*ent.CoreMenuCreate{
			client.CoreMenu.Create().SetID(fmt.Sprintf("%d", global.Id.GenID())).SetName("仪表盘").SetMenuCode("dashboard").SetMenuType(common.MenuTypeDirectory).SetAPIPath("/core/api/v1/dashboard").SetAPIPathMethod("GET").SetOrder(1),
			client.CoreMenu.Create().SetID(fmt.Sprintf("%d", global.Id.GenID())).SetName("网关管理").SetMenuCode("gateway_management").SetMenuType(common.MenuTypeDirectory).SetOrder(2),
			client.CoreMenu.Create().SetID(fmt.Sprintf("%d", global.Id.GenID())).SetName("系统管理").SetMenuCode("system_management").SetMenuType(common.MenuTypeDirectory).SetOrder(3),

			client.CoreMenu.Create().SetID(fmt.Sprintf("%d", global.Id.GenID())).SetName("实例列表").SetMenuCode("instances").SetMenuType(common.MenuTypeMenu).SetOrder(1).SetParentMenuCode("gateway_management"),
			client.CoreMenu.Create().SetID(fmt.Sprintf("%d", global.Id.GenID())).SetName("L4 代理").SetMenuCode("l4_proxy").SetMenuType(common.MenuTypeMenu).SetOrder(2).SetParentMenuCode("gateway_management"),
			client.CoreMenu.Create().SetID(fmt.Sprintf("%d", global.Id.GenID())).SetName("L7 代理").SetMenuCode("l7_proxy").SetMenuType(common.MenuTypeMenu).SetOrder(3).SetParentMenuCode("gateway_management"),
			client.CoreMenu.Create().SetID(fmt.Sprintf("%d", global.Id.GenID())).SetName("证书管理").SetMenuCode("certificates").SetMenuType(common.MenuTypeMenu).SetOrder(4).SetParentMenuCode("gateway_management"),

			client.CoreMenu.Create().SetID(fmt.Sprintf("%d", global.Id.GenID())).SetName("用户管理").SetMenuCode("user_management").SetMenuType(common.MenuTypeMenu).SetOrder(1).SetParentMenuCode("system_management"),
			client.CoreMenu.Create().SetID(fmt.Sprintf("%d", global.Id.GenID())).SetName("在线用户").SetMenuCode("online_users").SetMenuType(common.MenuTypeMenu).SetOrder(2).SetParentMenuCode("system_management"),
			client.CoreMenu.Create().SetID(fmt.Sprintf("%d", global.Id.GenID())).SetName("角色管理").SetMenuCode("role_management").SetMenuType(common.MenuTypeMenu).SetOrder(3).SetParentMenuCode("system_management"),
			client.CoreMenu.Create().SetID(fmt.Sprintf("%d", global.Id.GenID())).SetName("菜单管理").SetMenuCode("menu_management").SetMenuType(common.MenuTypeMenu).SetOrder(4).SetParentMenuCode("system_management"),
			client.CoreMenu.Create().SetID(fmt.Sprintf("%d", global.Id.GenID())).SetName("操作日志").SetMenuCode("operation_logs").SetMenuType(common.MenuTypeMenu).SetOrder(5).SetParentMenuCode("system_management"),

			client.CoreMenu.Create().SetID(fmt.Sprintf("%d", global.Id.GenID())).SetName("添加用户").SetMenuCode("add_user").SetMenuType(common.MenuTypeButton).SetAPIPath("/core/api/v1/system/user").SetAPIPathMethod("POST").SetOrder(1).SetParentMenuCode("user_management"),
			client.CoreMenu.Create().SetID(fmt.Sprintf("%d", global.Id.GenID())).SetName("添加角色").SetMenuCode("add_role").SetMenuType(common.MenuTypeButton).SetAPIPath("/core/api/v1/system/role").SetAPIPathMethod("POST").SetOrder(1).SetParentMenuCode("role_management"),
		}

		ctx = context.Background()
	)

	exists, err := client.CoreMenu.Query().Where(coremenu.DeletedAtIsNil()).Exist(ctx)
	if err != nil {
		global.Logger.Sugar().Error("core_menu表数据初始化失败")
		return err
	}
	if exists {
		global.Logger.Sugar().Info("core_menu表数据无需初始化")
	} else {
		menus, err := client.CoreMenu.CreateBulk(menus...).Save(ctx)
		if err != nil {
			return err
		}
		if menus == nil {
			return fmt.Errorf("core_menu表数据初始化失败")
		}
		global.Logger.Sugar().Info("core_menu表数据初始化成功")
	}
	return nil
}
