package system

import (
	"context"
	"time"

	"github.com/lyonmu/quebec/cmd/core/internal/dto/request"
	"github.com/lyonmu/quebec/cmd/core/internal/dto/response"
	"github.com/lyonmu/quebec/cmd/core/internal/ent/coreuser"
	"github.com/lyonmu/quebec/cmd/core/internal/ent/corerole"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/pkg/code"
	"github.com/lyonmu/quebec/pkg/constant"
)

func (s *SystemSvc) RolePage(ctx context.Context, req *request.SystemRolePageReq) (*response.SystemRoleListResp, error) {
	var (
		total    int
		items    = make([]*response.SystemRoleResp, 0)
		page     = (req.Page - 1) * req.PageSize
		pageSize = req.PageSize
		resp     = &response.SystemRoleListResp{}
		query    = global.EntClient.CoreRole.Query().Where(corerole.DeletedAtIsNil())
	)

	if len(req.Name) > 0 {
		query = query.Where(corerole.NameContains(req.Name))
	}

	if req.Status != 0 {
		query = query.Where(corerole.Status(req.Status))
	}

	total, err := query.Count(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("select core_role failed: %s", err)
		return nil, &code.RoleQueryFailed
	}

	rows, err := query.Offset(page).Limit(pageSize).All(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("select core_role failed: %s", err)
		return nil, &code.RoleQueryFailed
	}

	for _, row := range rows {
		item := response.SystemRoleResp{}
		item.LoadDb(row)
		// 统计该角色的用户数量
		count, err := row.QueryRoleToUser().Where(coreuser.DeletedAtIsNil()).Count(ctx)
		if err != nil {
			global.Logger.Sugar().Errorf("count role users failed: %s", err)
		} else {
			item.UsersCount = count
		}
		items = append(items, &item)
	}

	resp.Total = total
	resp.Items = items
	resp.Page = req.Page
	resp.PageSize = pageSize

	return resp, nil
}

func (s *SystemSvc) RoleList(ctx context.Context, req *request.SystemRoleListReq) (*response.SystemRoleListResp, error) {
	var (
		items = make([]*response.SystemRoleResp, 0)
		resp  = &response.SystemRoleListResp{}
		query = global.EntClient.CoreRole.Query().Where(corerole.DeletedAtIsNil())
	)

	if len(req.Name) > 0 {
		query = query.Where(corerole.NameContains(req.Name))
	}

	if req.Status != 0 {
		query = query.Where(corerole.Status(req.Status))
	}

	rows, err := query.All(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("select core_role failed: %s", err)
		return nil, &code.RoleQueryFailed
	}

	for _, row := range rows {
		item := response.SystemRoleResp{}
		item.LoadDb(row)
		// 统计该角色的用户数量
		count, err := row.QueryRoleToUser().Where(coreuser.DeletedAtIsNil()).Count(ctx)
		if err != nil {
			global.Logger.Sugar().Errorf("count role users failed: %s", err)
		} else {
			item.UsersCount = count
		}
		items = append(items, &item)
	}

	resp.Items = items

	return resp, nil
}

func (s *SystemSvc) RoleLabel(ctx context.Context) ([]*response.Options, error) {

	var (
		resp = make([]*response.Options, 0)
	)

	rows, err := global.EntClient.CoreRole.Query().
		Select(corerole.FieldID, corerole.FieldName). // 指定需要的字段
		Where(corerole.DeletedAtIsNil(), corerole.Status(constant.Yes)).
		All(ctx)

	if err != nil {
		global.Logger.Sugar().Errorf("select core_role failed: %s", err)
		return nil, &code.RoleQueryFailed
	}

	for _, v := range rows {
		resp = append(resp, &response.Options{
			Label: v.Name,
			Value: v.ID,
		})
	}

	return resp, nil
}

func (s *SystemSvc) RoleAdd(ctx context.Context, req *request.SystemRoleAddReq) error {

	exist, err := global.EntClient.CoreRole.Query().Where(corerole.Name(req.Name), corerole.DeletedAtIsNil()).Exist(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("select core_role failed: %s", err)
		return &code.RoleQueryFailed
	}

	if exist {
		return &code.RoleNameDuplicate
	}

	_, cerr := global.EntClient.CoreRole.Create().
		SetName(req.Name).
		SetNillableRemark(req.Remark).
		Save(ctx)

	if cerr != nil {
		global.Logger.Sugar().Errorf("add core_role failed: %s", cerr)
		return &code.RoleAddFailed
	}

	return nil
}

func (s *SystemSvc) RoleUpdate(ctx context.Context, id string, req *request.SystemRoleUpdateReq) error {

	exist, err := global.EntClient.CoreRole.Query().Where(corerole.Name(*req.Name), corerole.IDNEQ(id), corerole.DeletedAtIsNil()).Exist(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("select core_role failed: %s", err)
		return &code.RoleQueryFailed
	}

	if exist {
		return &code.RoleNameDuplicate
	}

	_, uerr := global.EntClient.CoreRole.
		UpdateOneID(id).
		SetNillableName(req.Name).
		SetNillableRemark(req.Remark).
		Save(ctx)
	if uerr != nil {
		global.Logger.Sugar().Errorf("update core_role failed: %s", uerr)
		return &code.RoleEditFailed
	}

	return nil
}

func (s *SystemSvc) RoleDelete(ctx context.Context, id string) error {

	row, err := global.EntClient.CoreRole.Query().Where(corerole.ID(id), corerole.DeletedAtIsNil()).First(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("select core_role failed: %s", err)
		return &code.RoleNotExists
	}

	if row.System == constant.Yes {
		return &code.RoleSystemNotAllow
	}

	_, derr := row.Update().SetDeletedAt(time.Now()).Save(ctx)
	if derr != nil {
		global.Logger.Sugar().Errorf("delete core_role failed: %s", derr)
		return &code.RoleDelFailed
	}

	return nil
}

func (s *SystemSvc) RoleGetById(ctx context.Context, id string) (*response.SystemRoleResp, error) {

	row, err := global.EntClient.CoreRole.Query().Where(corerole.ID(id), corerole.DeletedAtIsNil()).First(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("select core_role failed: %s", err)
		return nil, &code.RoleNotExists
	}

	item := response.SystemRoleResp{}
	item.LoadDb(row)

	return &item, nil
}

func (s *SystemSvc) RoleEnable(ctx context.Context, id string, req *request.EnableReq) error {
	row, qerr := global.EntClient.CoreRole.Query().Where(corerole.ID(id), corerole.DeletedAtIsNil()).First(ctx)
	if qerr != nil {
		global.Logger.Sugar().Errorf("select core_role failed: %s", qerr)
		return &code.RoleNotExists
	}

	if row.System == constant.Yes {
		return &code.RoleSystemNotAllow
	}

	_, err := global.EntClient.CoreRole.UpdateOneID(id).SetStatus(req.Status).Save(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("enable core_role failed: %s", err)
		return &code.RoleEnableFailed
	}

	return nil
}
