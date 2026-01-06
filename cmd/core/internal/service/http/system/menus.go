package system

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	corecommon "github.com/lyonmu/quebec/cmd/core/internal/common"
	"github.com/lyonmu/quebec/cmd/core/internal/dto/request"
	"github.com/lyonmu/quebec/cmd/core/internal/dto/response"
	"github.com/lyonmu/quebec/cmd/core/internal/ent"
	"github.com/lyonmu/quebec/cmd/core/internal/ent/coredatarelationship"
	"github.com/lyonmu/quebec/cmd/core/internal/ent/coremenu"
	"github.com/lyonmu/quebec/cmd/core/internal/ent/corerole"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/pkg/code"
	"github.com/lyonmu/quebec/pkg/constant"
)

func (s *SystemSvc) MenuPage(ctx context.Context, req *request.SystemMenuPageReq) (*response.SystemMenuListResp, error) {
	var (
		total    int
		items    = make([]*response.SystemMenuResp, 0)
		page     = (req.Page - 1) * req.PageSize
		pageSize = req.PageSize
		resp     = &response.SystemMenuListResp{}
		query    = global.EntClient.CoreMenu.Query().Where(coremenu.DeletedAtIsNil())
	)

	if len(req.Name) > 0 {
		query = query.Where(coremenu.NameContains(req.Name))
	}

	if len(req.MenuCode) > 0 {
		query = query.Where(coremenu.MenuCodeContains(req.MenuCode))
	}

	if req.MenuType != 0 {
		query = query.Where(coremenu.MenuType(req.MenuType))
	}

	if req.Status != 0 {
		query = query.Where(coremenu.Status(req.Status))
	}

	if len(req.ParentMenuCode) > 0 {
		query = query.Where(coremenu.ParentMenuCode(req.ParentMenuCode))
	}

	total, err := query.Count(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("select core_menu failed: %s", err)
		return nil, &code.MenuQueryFailed
	}

	rows, err := query.Offset(page).Limit(pageSize).Order(coremenu.ByID(sql.OrderAsc())).All(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("select core_menu failed: %s", err)
		return nil, &code.MenuQueryFailed
	}

	for _, row := range rows {
		item := response.SystemMenuResp{}
		item.LoadDb(row)
		items = append(items, &item)
	}

	resp.Total = total
	resp.Items = items
	resp.Page = req.Page
	resp.PageSize = pageSize

	return resp, nil
}

func (s *SystemSvc) MenuList(ctx context.Context, req *request.SystemMenuListReq) ([]*response.SystemMenuResp, error) {
	var (
		items = make([]*response.SystemMenuResp, 0)
		query = global.EntClient.CoreMenu.Query().Where(coremenu.DeletedAtIsNil())
	)

	if len(req.Name) > 0 {
		query = query.Where(coremenu.NameContains(req.Name))
	}

	if len(req.MenuCode) > 0 {
		query = query.Where(coremenu.MenuCodeContains(req.MenuCode))
	}

	if req.MenuType != 0 {
		query = query.Where(coremenu.MenuType(req.MenuType))
	}

	if req.Status != 0 {
		query = query.Where(coremenu.Status(req.Status))
	}

	if len(req.ParentMenuCode) > 0 {
		query = query.Where(coremenu.ParentMenuCode(req.ParentMenuCode))
	}

	rows, err := query.Order(coremenu.ByID(sql.OrderAsc())).All(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("select core_menu failed: %s", err)
		return nil, &code.MenuQueryFailed
	}

	for _, row := range rows {
		item := response.SystemMenuResp{}
		item.LoadDb(row)
		items = append(items, &item)
	}

	return items, nil
}

func (s *SystemSvc) MenuTree(ctx context.Context) ([]*response.SystemMenuTreeResp, error) {
	// Get all root menus (parent_id is empty)
	rootMenus, err := global.EntClient.CoreMenu.Query().
		Where(coremenu.DeletedAtIsNil(), coremenu.ParentMenuCodeIsNil()).
		Order(coremenu.ByOrder()).
		All(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("select core_menu failed: %s", err)
		return nil, &code.MenuQueryFailed
	}

	return buildMenuTree(rootMenus, ctx), nil
}

func buildMenuTree(menus []*ent.CoreMenu, ctx context.Context) []*response.SystemMenuTreeResp {
	result := make([]*response.SystemMenuTreeResp, 0, len(menus))

	for _, menu := range menus {
		node := &response.SystemMenuTreeResp{
			ID:            menu.ID,
			Name:          menu.Name,
			MenuType:      menu.MenuType,
			ApiPath:       menu.APIPath,
			ApiPathMethod: menu.APIPathMethod,
			Order:         menu.Order,
			MenuCode:      menu.MenuCode,
			Status:        menu.Status,
			Children:      nil,
		}

		// Get children
		children, err := global.EntClient.CoreMenu.Query().
			Where(coremenu.DeletedAtIsNil(), coremenu.ParentMenuCode(menu.MenuCode)).
			Order(coremenu.ByOrder()).
			All(ctx)
		if err != nil {
			global.Logger.Sugar().Errorf("select core_menu children failed: %s", err)
			continue
		}

		if len(children) > 0 {
			node.Children = buildMenuTree(children, ctx)
		}

		result = append(result, node)
	}

	return result
}

func (s *SystemSvc) MenuLabel(ctx context.Context) ([]*response.Options, error) {
	var (
		resp = make([]*response.Options, 0)
	)

	rows, err := global.EntClient.CoreMenu.Query().
		Select(coremenu.FieldID, coremenu.FieldName).
		Where(coremenu.DeletedAtIsNil(), coremenu.Status(constant.Yes)).
		All(ctx)

	if err != nil {
		global.Logger.Sugar().Errorf("select core_menu failed: %s", err)
		return nil, &code.MenuQueryFailed
	}

	for _, v := range rows {
		resp = append(resp, &response.Options{
			Label: v.Name,
			Value: v.ID,
		})
	}

	return resp, nil
}

func (s *SystemSvc) MenuGetById(ctx context.Context, id string) (*response.SystemMenuResp, error) {
	row, err := global.EntClient.CoreMenu.Query().Where(coremenu.ID(id), coremenu.DeletedAtIsNil()).First(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("select core_menu failed: %s", err)
		return nil, &code.MenuNotExists
	}

	item := response.SystemMenuResp{}
	item.LoadDb(row)

	return &item, nil
}

func (s *SystemSvc) MenuEnable(ctx context.Context, id string, req *request.EnableReq) error {
	_, qerr := global.EntClient.CoreMenu.Query().Where(coremenu.ID(id), coremenu.DeletedAtIsNil()).First(ctx)
	if qerr != nil {
		global.Logger.Sugar().Errorf("select core_menu failed: %s", qerr)
		return &code.MenuNotExists
	}

	_, err := global.EntClient.CoreMenu.UpdateOneID(id).SetStatus(req.Status).Save(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("enable core_menu failed: %s", err)
		return &code.MenuEnableFailed
	}

	return nil
}

// Role Menu Binding

func (s *SystemSvc) GetRoleMenus(ctx context.Context, roleID string) ([]*response.SystemRoleMenuResp, error) {
	// Check if role exists
	roleExist, err := global.EntClient.CoreRole.Query().Where(corerole.ID(roleID), corerole.DeletedAtIsNil()).Exist(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("select core_role failed: %s", err)
		return nil, &code.RoleQueryFailed
	}
	if !roleExist {
		return nil, &code.RoleNotExists
	}

	// Get role-menu relationships
	relationships, err := global.EntClient.CoreDataRelationship.Query().
		Where(
			coredatarelationship.DeletedAtIsNil(),
			coredatarelationship.RoleID(roleID),
			coredatarelationship.DataRelationshipType(corecommon.DataRelationshipTypeRoleToMenu),
		).
		All(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("select core_data_relationship failed: %s", err)
		return nil, &code.MenuQueryFailed
	}

	result := make([]*response.SystemRoleMenuResp, 0, len(relationships))
	for _, rel := range relationships {
		resp := &response.SystemRoleMenuResp{
			MenuID: rel.MenuID,
		}
		// 使用 QueryMenu 获取关联的菜单信息
		menus, err := rel.QueryMenu().All(ctx)
		if err != nil {
			global.Logger.Sugar().Errorf("query menu failed: %s", err)
		} else if len(menus) > 0 {
			resp.MenuName = menus[0].Name
		}
		result = append(result, resp)
	}

	return result, nil
}

func (s *SystemSvc) BindRoleMenus(ctx context.Context, roleID string, req *request.SystemRoleMenuBindReq) error {
	// Check if role exists
	roleExist, err := global.EntClient.CoreRole.Query().Where(corerole.ID(roleID), corerole.DeletedAtIsNil()).Exist(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("select core_role failed: %s", err)
		return &code.RoleQueryFailed
	}
	if !roleExist {
		return &code.RoleNotExists
	}

	// Delete existing role-menu relationships
	_, err = global.EntClient.CoreDataRelationship.Delete().
		Where(
			coredatarelationship.DeletedAtIsNil(),
			coredatarelationship.RoleID(roleID),
			coredatarelationship.DataRelationshipType(corecommon.DataRelationshipTypeRoleToMenu),
		).
		Exec(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("delete core_data_relationship failed: %s", err)
		return &code.RoleMenuBindFailed
	}

	// Create new relationships
	bulk := make([]*ent.CoreDataRelationshipCreate, 0, len(req.MenuIDs))
	for _, menuID := range req.MenuIDs {
		bulk = append(bulk, global.EntClient.CoreDataRelationship.Create().
			SetDataRelationshipType(corecommon.DataRelationshipTypeRoleToMenu).
			SetMenuID(menuID).
			SetRoleID(roleID))
	}

	if len(bulk) > 0 {
		_, cerr := global.EntClient.CoreDataRelationship.CreateBulk(bulk...).Save(ctx)
		if cerr != nil {
			global.Logger.Sugar().Errorf("create core_data_relationship failed: %s", cerr)
			return &code.RoleMenuBindFailed
		}
	}

	return nil
}

func (s *SystemSvc) AddRoleMenu(ctx context.Context, roleID string, menuID string) error {
	// Check if role exists
	roleExist, err := global.EntClient.CoreRole.Query().Where(corerole.ID(roleID), corerole.DeletedAtIsNil()).Exist(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("select core_role failed: %s", err)
		return &code.RoleQueryFailed
	}
	if !roleExist {
		return &code.RoleNotExists
	}

	// Check if menu exists
	menuExist, err := global.EntClient.CoreMenu.Query().Where(coremenu.ID(menuID), coremenu.DeletedAtIsNil()).Exist(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("select core_menu failed: %s", err)
		return &code.MenuQueryFailed
	}
	if !menuExist {
		return &code.MenuNotExists
	}

	// Check if relationship already exists
	exist, err := global.EntClient.CoreDataRelationship.Query().
		Where(
			coredatarelationship.DeletedAtIsNil(),
			coredatarelationship.RoleID(roleID),
			coredatarelationship.MenuID(menuID),
			coredatarelationship.DataRelationshipType(corecommon.DataRelationshipTypeRoleToMenu),
		).
		Exist(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("select core_data_relationship failed: %s", err)
		return &code.RoleMenuBindFailed
	}
	if exist {
		return nil // Already exists, skip
	}

	// Create new relationship
	_, cerr := global.EntClient.CoreDataRelationship.Create().
		SetDataRelationshipType(corecommon.DataRelationshipTypeRoleToMenu).
		SetMenuID(menuID).
		SetRoleID(roleID).
		Save(ctx)
	if cerr != nil {
		global.Logger.Sugar().Errorf("create core_data_relationship failed: %s", cerr)
		return &code.RoleMenuBindFailed
	}

	return nil
}

func (s *SystemSvc) RemoveRoleMenu(ctx context.Context, roleID string, menuID string) error {
	// Check if relationship exists
	exist, err := global.EntClient.CoreDataRelationship.Query().
		Where(
			coredatarelationship.DeletedAtIsNil(),
			coredatarelationship.RoleID(roleID),
			coredatarelationship.MenuID(menuID),
			coredatarelationship.DataRelationshipType(corecommon.DataRelationshipTypeRoleToMenu),
		).
		First(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("select core_data_relationship failed: %s", err)
		return &code.RoleMenuBindFailed
	}

	// Delete the relationship
	_, derr := exist.Update().SetDeletedAt(time.Now()).Save(ctx)
	if derr != nil {
		global.Logger.Sugar().Errorf("delete core_data_relationship failed: %s", derr)
		return &code.RoleMenuBindFailed
	}

	return nil
}
