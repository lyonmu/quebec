package system

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/lyonmu/quebec/cmd/core/internal/common"
	"github.com/lyonmu/quebec/cmd/core/internal/dto/request"
	"github.com/lyonmu/quebec/cmd/core/internal/dto/response"
	"github.com/lyonmu/quebec/cmd/core/internal/ent"
	"github.com/lyonmu/quebec/cmd/core/internal/ent/coreonlineuser"
	"github.com/lyonmu/quebec/cmd/core/internal/ent/coreoperationlog"
	"github.com/lyonmu/quebec/cmd/core/internal/ent/coreuser"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/pkg/code"
)

// CreateOperationLogWithOnlineUserUpdate 使用事务方式同时创建操作日志和更新在线用户信息
func (s *SystemSvc) CreateOperationLogWithOnlineUserUpdate(ctx context.Context, req *request.OperationLogReq) error {
	tx, err := global.EntClient.Tx(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("Failed to start transaction: %v", err)
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if err != nil {
			global.Logger.Sugar().Errorf("Transaction rolled back due to error: %v", err)
			_ = tx.Rollback()
		}
	}()

	// 1. 创建操作日志
	if _, err := tx.CoreOperationLog.Create().
		SetUserID(req.ID).
		SetAccessIP(req.AccessIP).
		SetOperationTime(req.OperationTime).
		SetOperationType(req.OperationType).
		SetOs(req.Os).
		SetPlatform(req.Platform).
		SetBrowserName(req.BrowserName).
		SetBrowserVersion(req.BrowserVersion).
		SetBrowserEngineName(req.BrowserEngineName).
		SetBrowserEngineVersion(req.BrowserEngineVersion).
		Save(ctx); err != nil {
		global.Logger.Sugar().Errorf("Failed to create operation log: %v", err)
		return fmt.Errorf("failed to create operation log: %w", err)
	}

	// 2. 更新在线用户
	if _, err := tx.CoreOnLineUser.Update().
		Where(coreonlineuser.UserID(req.ID)).
		SetLastOperationTime(req.OperationTime).
		SetOperationType(req.OperationType).
		Save(ctx); err != nil {
		global.Logger.Sugar().Errorf("Failed to update online user: %v", err)
		return fmt.Errorf("failed to update online user: %w", err)
	}

	// 3. 提交事务
	if err := tx.Commit(); err != nil {
		global.Logger.Sugar().Errorf("Failed to commit transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *SystemSvc) OperationLogPage(ctx context.Context, req *request.OperationLogPageReq) (*response.SystemOperationLogListResp, error) {
	var (
		total    int
		items    = make([]*response.SystemOperationLogResp, 0)
		page     = (req.Page - 1) * req.PageSize
		pageSize = req.PageSize
		resp     = &response.SystemOperationLogListResp{}
		query    = global.EntClient.CoreOperationLog.Query()
	)

	if len(req.ID) > 0 {
		query = query.Where(coreoperationlog.UserID(req.ID))
	}

	if req.OperationType != 0 {
		query = query.Where(coreoperationlog.OperationType(common.OperationType(req.OperationType)))
	}

	if req.StartTime > 0 {
		query = query.Where(coreoperationlog.OperationTimeGTE(req.StartTime))
	}

	if req.EndTime > 0 {
		query = query.Where(coreoperationlog.OperationTimeLTE(req.EndTime))
	}

	total, err := query.Count(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("select operation log failed: %s", err)
		return nil, &code.LogQueryFailed
	}

	rows, err := query.Offset(page).Limit(pageSize).
		WithOperationLogFromUser(
			func(q *ent.CoreUserQuery) {
				q.Select(coreuser.FieldID, coreuser.FieldUsername, coreuser.FieldNickname).Where(coreuser.DeletedAtIsNil())
			},
		).Order(coreoperationlog.ByOperationTime(sql.OrderDesc())).All(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("select operation log failed: %s", err)
		return nil, &code.LogQueryFailed
	}

	for _, row := range rows {
		item := response.SystemOperationLogResp{}
		item.LoadDb(row)
		items = append(items, &item)
	}

	resp.Total = total
	resp.Items = items
	resp.Page = req.Page
	resp.PageSize = pageSize

	return resp, nil
}
