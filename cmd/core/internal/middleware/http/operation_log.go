package http

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lyonmu/quebec/cmd/core/internal/common"
	"github.com/lyonmu/quebec/cmd/core/internal/dto/request"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/cmd/core/internal/service/http/system"
	"github.com/lyonmu/quebec/cmd/core/internal/utils"
	"github.com/lyonmu/quebec/pkg/constant"
	"github.com/lyonmu/quebec/pkg/tools"
)

type OperationLogMiddleware struct {
	systemSvc *system.SystemSvc
}

func NewOperationLogMiddleware(systemSvc *system.SystemSvc) *OperationLogMiddleware {
	return &OperationLogMiddleware{
		systemSvc: systemSvc,
	}
}

func (m *OperationLogMiddleware) getUserClaims(c *gin.Context) *utils.JwtClaims {
	token := c.GetHeader(constant.ApiTokenName)
	if token == "" {
		return nil
	}

	claims, err := global.JwtToolEntity.ParseToken(token, global.Cfg.Core.Jwt.Sign)
	if err != nil {
		return nil
	}

	return claims
}

func (m *OperationLogMiddleware) Handle(operationType common.OperationType) gin.HandlerFunc {
	return func(c *gin.Context) {
		cC := c.Copy()
		c.Next()
		claims := m.getUserClaims(cC)
		if len(claims.UserId) != 0 {
			m.recordOperation(cC, operationType, claims.UserId)
		}
	}
}

func (m *OperationLogMiddleware) recordOperation(c *gin.Context, operationType common.OperationType, userID string) {
	uaStr := c.Request.UserAgent()
	ua := tools.ParseUserAgent(uaStr)
	access_ip := c.ClientIP()
	browserName, browserVersion := ua.Browser()
	browserEngineName, browserEngineVersion := ua.Engine()

	req := &request.OperationLogReq{
		ID:                   userID,
		AccessIP:             access_ip,
		OperationTime:        time.Now().Unix(),
		OperationType:        operationType,
		Os:                   ua.OS(),
		Platform:             ua.Platform(),
		BrowserName:          browserName,
		BrowserVersion:       browserVersion,
		BrowserEngineName:    browserEngineName,
		BrowserEngineVersion: browserEngineVersion,
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				global.Logger.Sugar().Errorf("Operation log goroutine panicked: %v", r)
			}
		}()

		// 创建独立的 Context，避免主请求 Context 取消的影响
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := m.systemSvc.CreateOperationLogWithOnlineUserUpdate(ctx, req); err != nil {
			global.Logger.Sugar().Errorf("Failed to create operation log and update online user: %v", err)
		}

	}()
}
