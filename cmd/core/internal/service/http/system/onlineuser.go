package system

import (
	"context"
	"time"

	"github.com/lyonmu/quebec/cmd/core/internal/dto/request"
	"github.com/lyonmu/quebec/cmd/core/internal/dto/response"
	"github.com/lyonmu/quebec/cmd/core/internal/ent"
	"github.com/lyonmu/quebec/cmd/core/internal/ent/coreonlineuser"
	"github.com/lyonmu/quebec/cmd/core/internal/ent/coreuser"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/pkg/code"
)

func (s *SystemSvc) ListOnlineUser(req *request.SystemOnlineUserListReq, ctx context.Context) (*response.SystemOnlineUserListResp, error) {
	var (
		total    int
		items    = make([]*response.SystemOnlineUserResp, 0)
		page     = (req.Page - 1) * req.PageSize
		pageSize = req.PageSize
		resp     = &response.SystemOnlineUserListResp{}
	)
	query := global.EntClient.CoreOnLineUser.Query().Where(coreonlineuser.DeletedAtIsNil())
	if len(req.UserID) > 0 {
		query = query.Where(coreonlineuser.UserIDEQ(req.UserID))
	}
	if len(req.AccessIP) > 0 {
		query = query.Where(coreonlineuser.AccessIPEQ(req.AccessIP))
	}
	if req.StartTime > 0 {
		query = query.Where(coreonlineuser.LastOperationTimeGTE(req.StartTime))
	}
	if req.EndTime > 0 {
		query = query.Where(coreonlineuser.LastOperationTimeLTE(req.EndTime))
	}

	total, err := query.Count(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("获取在线用户总数失败: %v", err)
		return nil, &code.Failed
	}

	rows, err := query.Offset(page).Limit(pageSize).
		WithOnLineFromUser(
			func(q *ent.CoreUserQuery) {
				q.Select(coreuser.FieldID, coreuser.FieldNickname).Where(coreuser.DeletedAtIsNil())
			},
		).
		All(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("获取在线用户列表失败: %v", err)
		return nil, &code.Failed
	}

	for _, row := range rows {
		item := response.SystemOnlineUserResp{}
		item.LoadDb(row)
		items = append(items, &item)
	}

	resp.Total = total
	resp.Items = items
	resp.Page = req.Page
	resp.PageSize = pageSize

	return resp, nil
}

func (s *SystemSvc) UserLabel(ctx context.Context) ([]*response.Options, error) {

	var (
		resp = make([]*response.Options, 0)
	)

	rows, err := global.EntClient.CoreOnLineUser.Query().
		WithOnLineFromUser(
			func(q *ent.CoreUserQuery) {
				q.Select(coreuser.FieldID, coreuser.FieldNickname).Where(coreuser.DeletedAtIsNil())
			},
		).
		Select(coreonlineuser.FieldID, coreonlineuser.FieldUserID).
		Where(coreonlineuser.DeletedAtIsNil()).
		All(ctx)

	if err != nil {
		global.Logger.Sugar().Errorf("获取在线用户标签失败: %v", err)
		return nil, &code.Failed
	}

	for _, v := range rows {

		if v.Edges.OnLineFromUser != nil {
			resp = append(resp, &response.Options{
				Label: v.Edges.OnLineFromUser.Nickname,
				Value: v.UserID,
			})
		}
	}
	return resp, nil
}

func (s *SystemSvc) Clearance(req *request.IdReq, ctx context.Context) error {

	if perr := global.JwtToolEntity.DeleteToken(req.ID, global.RedisCli); perr != nil {
		global.Logger.Sugar().Errorf("清除在线用户删除token失败: %v", perr)
		return &code.Failed
	}

	_, uerr := global.EntClient.CoreOnLineUser.Update().
		Where(coreonlineuser.UserIDEQ(req.ID), coreonlineuser.DeletedAtIsNil()).
		SetDeletedAt(time.Now()).
		Save(ctx)
	if uerr != nil {
		global.Logger.Sugar().Errorf("清除在线用户失败: %v", uerr)
		return &code.Failed
	}
	return nil
}
