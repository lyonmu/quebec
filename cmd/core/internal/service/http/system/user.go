package system

import (
	"context"
	"fmt"
	"time"

	"github.com/lyonmu/quebec/cmd/core/internal/dto/request"
	"github.com/lyonmu/quebec/cmd/core/internal/dto/response"
	"github.com/lyonmu/quebec/cmd/core/internal/ent"
	"github.com/lyonmu/quebec/cmd/core/internal/ent/coreonlineuser"
	"github.com/lyonmu/quebec/cmd/core/internal/ent/corerole"
	"github.com/lyonmu/quebec/cmd/core/internal/ent/coreuser"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/pkg/code"
	"github.com/lyonmu/quebec/pkg/constant"
	"github.com/lyonmu/quebec/pkg/tools/encrypt"
)

func (s *SystemSvc) UserPage(ctx context.Context, req *request.SystemUserPageReq) (*response.SystemUserListResp, error) {
	var (
		total    int
		items    = make([]*response.SystemUserResp, 0)
		page     = (req.Page - 1) * req.PageSize
		pageSize = req.PageSize
		resp     = &response.SystemUserListResp{}
		query    = global.EntClient.CoreUser.Query().Where(coreuser.DeletedAtIsNil())
	)

	if len(req.Username) > 0 {
		query = query.Where(coreuser.UsernameContains(req.Username))
	}

	if req.Status != 0 {
		query = query.Where(coreuser.Status(req.Status))
	}

	if len(req.RoleID) > 0 {
		query = query.Where(coreuser.RoleID(req.RoleID))
	}

	if len(req.Email) > 0 {
		query = query.Where(coreuser.EmailContains(req.Email))
	}

	if len(req.Nickname) > 0 {
		query = query.Where(coreuser.NicknameContains(req.Nickname))
	}

	total, err := query.Count(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("select core_user failed: %s", err)
		return nil, &code.UserQueryFailed
	}

	rows, err := query.Offset(page).Limit(pageSize).
		WithUserFromRole(
			func(q *ent.CoreRoleQuery) {
				q.Select(corerole.FieldID, corerole.FieldName).Where(corerole.DeletedAtIsNil())
			},
		).All(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("select core_user failed: %s", err)
		return nil, &code.UserQueryFailed
	}

	for _, row := range rows {
		item := response.SystemUserResp{}
		item.LoadDb(row)
		items = append(items, &item)
	}

	resp.Total = total
	resp.Items = items
	resp.Page = req.Page
	resp.PageSize = pageSize

	return resp, nil
}

func (s *SystemSvc) UserList(ctx context.Context, req *request.SystemUserListReq) (*response.SystemUserListResp, error) {
	var (
		items = make([]*response.SystemUserResp, 0)
		resp  = &response.SystemUserListResp{}
		query = global.EntClient.CoreUser.Query().Where(coreuser.DeletedAtIsNil())
	)

	if len(req.Username) > 0 {
		query = query.Where(coreuser.UsernameContains(req.Username))
	}

	if req.Status != 0 {
		query = query.Where(coreuser.Status(req.Status))
	}

	rows, err := query.
		WithUserFromRole(
			func(q *ent.CoreRoleQuery) {
				q.Select(corerole.FieldID, corerole.FieldName).Where(corerole.DeletedAtIsNil())
			},
		).All(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("select core_user failed: %s", err)
		return nil, &code.UserQueryFailed
	}

	for _, row := range rows {
		item := response.SystemUserResp{}
		item.LoadDb(row)
		items = append(items, &item)
	}

	resp.Items = items

	return resp, nil
}

func (s *SystemSvc) UserLabel(ctx context.Context) ([]*response.Options, error) {

	var (
		resp = make([]*response.Options, 0)
	)

	rows, err := global.EntClient.CoreUser.Query().
		Select(coreuser.FieldID, coreuser.FieldNickname, coreuser.FieldUsername). // 指定需要的字段
		Where(coreuser.DeletedAtIsNil(), coreuser.Status(constant.Yes)).
		All(ctx)

	if err != nil {
		global.Logger.Sugar().Errorf("select core_user failed: %s", err)
		return nil, &code.UserQueryFailed
	}

	for _, v := range rows {
		resp = append(resp, &response.Options{
			Label: fmt.Sprintf("%s(%s)", v.Nickname, v.Username),
			Value: v.ID,
		})
	}

	return resp, nil
}

func (s *SystemSvc) UserAdd(ctx context.Context, req *request.SystemUserAddReq) error {

	roleExist, rerr := global.EntClient.CoreRole.Query().Where(corerole.ID(req.RoleID), corerole.DeletedAtIsNil()).Exist(ctx)
	if rerr != nil {
		global.Logger.Sugar().Errorf("select core_role failed: %s", rerr)
		return &code.RoleQueryFailed
	}
	if !roleExist {
		return &code.RoleNotExists
	}

	password, err := encrypt.HashWithBcryptString(req.Password)
	if err != nil {
		global.Logger.Sugar().Errorf("hash password failed: %s", err)
		return &code.UserAddFailed
	}

	exist, err := global.EntClient.CoreUser.Query().Where(coreuser.Username(req.Username), coreuser.DeletedAtIsNil()).Exist(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("select core_user failed: %s", err)
		return &code.UserQueryFailed
	}

	if exist {
		return &code.UserNameDuplicate
	}

	_, cerr := global.EntClient.CoreUser.Create().
		SetUsername(req.Username).
		SetNickname(req.Nickname).
		SetRoleID(req.RoleID).
		SetPassword(password).
		SetNillableRemark(req.Remark).
		SetNillableEmail(req.Email).
		SetNillableStatus(req.Status).
		Save(ctx)

	if cerr != nil {
		global.Logger.Sugar().Errorf("add core_user failed: %s", cerr)
		return &code.UserAddFailed
	}

	return nil
}

func (s *SystemSvc) UserUpdate(ctx context.Context, id string, req *request.SystemUserUpdateReq) error {

	if len(*req.RoleID) > 0 {
		roleExist, rerr := global.EntClient.CoreRole.Query().Where(corerole.ID(*req.RoleID), corerole.DeletedAtIsNil()).Exist(ctx)
		if rerr != nil {
			global.Logger.Sugar().Errorf("select core_role failed: %s", rerr)
			return &code.RoleQueryFailed
		}
		if !roleExist {
			return &code.RoleNotExists
		}
	}

	if len(*req.Username) > 0 {
		exist, err := global.EntClient.CoreUser.Query().Where(coreuser.Username(*req.Username), coreuser.IDNEQ(id), coreuser.DeletedAtIsNil()).Exist(ctx)
		if err != nil {
			global.Logger.Sugar().Errorf("select core_user failed: %s", err)
			return &code.UserQueryFailed
		}
		if exist {
			return &code.UserNameDuplicate
		}
	}

	_, uerr := global.EntClient.CoreUser.
		UpdateOneID(id).
		SetNillableUsername(req.Username).
		SetNillableNickname(req.Nickname).
		SetNillableEmail(req.Email).
		SetNillableRoleID(req.RoleID).
		SetNillablePassword(req.Email).
		SetNillableRemark(req.Remark).
		Save(ctx)
	if uerr != nil {
		global.Logger.Sugar().Errorf("update core_user failed: %s", uerr)
		return &code.UserEditFailed
	}

	return nil
}

func (s *SystemSvc) UserDelete(ctx context.Context, id string) error {

	row, err := global.EntClient.CoreUser.Query().Where(coreuser.ID(id), coreuser.DeletedAtIsNil()).First(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("select core_user failed: %s", err)
		return &code.UserNotExists
	}

	if row.System == constant.Yes {
		return &code.UserSystemNotAllow
	}

	_, derr := row.Update().SetDeletedAt(time.Now()).Save(ctx)
	if derr != nil {
		global.Logger.Sugar().Errorf("delete core_user failed: %s", derr)
		return &code.UserDelFailed
	}

	return nil
}

func (s *SystemSvc) UserGetById(ctx context.Context, id string) (*response.SystemUserResp, error) {

	row, err := global.EntClient.CoreUser.Query().
		WithUserFromRole(
			func(q *ent.CoreRoleQuery) {
				q.Select(corerole.FieldID, corerole.FieldName).Where(corerole.DeletedAtIsNil())
			},
		).WithOnLineToUser(
		func(q *ent.CoreOnLineUserQuery) {
			q.Select(coreonlineuser.FieldID, coreonlineuser.FieldOperationType, coreonlineuser.FieldLastOperationTime).Where(coreonlineuser.DeletedAtIsNil())
		},
	).Where(coreuser.ID(id), coreuser.DeletedAtIsNil()).First(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("select core_user failed: %s", err)
		return nil, &code.UserNotExists
	}

	item := response.SystemUserResp{}
	item.LoadDb(row)

	return &item, nil
}

func (s *SystemSvc) UserEnable(ctx context.Context, id string, req *request.EnableReq) error {
	row, qerr := global.EntClient.CoreUser.Query().Where(coreuser.ID(id), coreuser.DeletedAtIsNil()).First(ctx)
	if qerr != nil {
		global.Logger.Sugar().Errorf("select core_user failed: %s", qerr)
		return &code.UserNotExists
	}

	if row.System == constant.Yes {
		return &code.UserSystemNotAllow
	}

	_, err := global.EntClient.CoreRole.UpdateOneID(id).SetStatus(req.Status).Save(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("enable core_user failed: %s", err)
		return &code.UserEnableFailed
	}

	return nil
}

func (s *SystemSvc) UserEditPassword(ctx context.Context, id string, req *request.SystemUserEditPasswordReq) error {

	row, qerr := global.EntClient.CoreUser.Query().Where(coreuser.ID(id), coreuser.DeletedAtIsNil()).First(ctx)
	if qerr != nil {
		global.Logger.Sugar().Errorf("select core_user failed: %s", qerr)
		return &code.UserNotExists
	}

	if row.Password != req.PrePassword {
		return &code.UserPasswordError
	}
	password, err := encrypt.HashWithBcryptString(req.NewPassword)
	if err != nil {
		global.Logger.Sugar().Errorf("hash password failed: %s", err)
		return &code.UserPasswordEditFailed
	}

	_, err = global.EntClient.CoreUser.UpdateOneID(id).SetPassword(password).Save(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("update core_user failed: %s", err)
		return &code.UserPasswordEditFailed
	}

	return nil
}
