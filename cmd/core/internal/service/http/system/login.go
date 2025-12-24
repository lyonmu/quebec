package system

import (
	"context"
	"time"

	"github.com/lyonmu/quebec/cmd/core/internal/common"
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
	"github.com/mssola/useragent"
)

func (s *SystemSvc) Login(req *request.SystemLoginRequest, ua *useragent.UserAgent, uaStr, access_ip string, ctx context.Context) (*response.SystemInfoResponse, error) {
	var (
		resp response.SystemInfoResponse
	)

	if req.CaptchaId != "" && req.Captcha != "" {
		if !global.CaptchaGenerator.Verify(req.CaptchaId, req.Captcha, true) {
			global.Logger.Sugar().Error("用户登录验证码错误")
			return nil, &code.InvalidCaptcha
		}
	} else {
		global.Logger.Sugar().Errorf("用户登录验证码无效")
		return nil, &code.InvalidCaptcha
	}

	u, qerr := global.EntClient.CoreUser.Query().
		WithUserFromRole(func(q *ent.CoreRoleQuery) {
			q.Select(corerole.FieldID, corerole.FieldName, corerole.FieldStatus).Where(corerole.DeletedAtIsNil())
		}).
		Where(coreuser.UsernameEQ(req.Username)).
		Where(coreuser.DeletedAtIsNil()).
		First(ctx)

	if qerr != nil {
		global.Logger.Sugar().Errorf("用户登录查询数据库失败:", qerr)
		return nil, &code.UserNotExists
	}

	if u.Edges.UserFromRole != nil {
		if u.Edges.UserFromRole.Status != constant.Yes {
			global.Logger.Sugar().Errorf("角色 %s 被禁用", u.Edges.UserFromRole.Name)
			return nil, &code.RoleDisabled
		}
	} else {
		global.Logger.Sugar().Error("角色查询失败")
		return nil, &code.RoleNotExists
	}

	if encrypt.CompareWithBcryptString(u.Password, req.Password) {
		global.Logger.Sugar().Errorf("用户 %s 登录密码错误", req.Username)
		return nil, &code.UserPasswordError
	}

	token, terr := global.JwtToolEntity.GenerateToken(u.ID, global.Cfg.Core.Jwt.Sign, u.LastPasswordChange, global.Cfg.Core.Jwt.Cache)
	if terr != nil {
		global.Logger.Sugar().Errorf("用户 %s 登录生成token失败: %v", req.Username, terr)
		return nil, &code.Failed
	}

	if err := global.JwtToolEntity.StoreToken(token, u.ID, global.Cfg.Core.Jwt.Cache, global.RedisCli); err != nil {
		global.Logger.Sugar().Errorf("用户 %s 登录存储token失败: %v", req.Username, err)
		return nil, &code.Failed
	}

	browserName, browserVersion := ua.Browser()
	browserEngineName, browserEngineVersion := ua.Engine()

	tx, terr := global.EntClient.Tx(ctx)
	if terr != nil {
		global.Logger.Sugar().Errorf("Failed to start transaction: %v", terr)
		return nil, &code.Failed
	}

	defer func() {
		if terr != nil {
			global.Logger.Sugar().Errorf("Transaction rolled back due to error: %v", terr)
			_ = tx.Rollback()
		}
	}()

	exists, terr := tx.CoreOnLineUser.Query().Where(coreonlineuser.UserIDEQ(u.ID), coreonlineuser.DeletedAtIsNil()).Exist(ctx)
	if terr != nil {
		global.Logger.Sugar().Errorf("用户 %s 登录查询在线用户失败: %v", req.Username, terr)
		return nil, &code.Failed
	}

	if exists {
		_, terr = tx.CoreOnLineUser.Update().
			Where(coreonlineuser.UserIDEQ(u.ID), coreonlineuser.DeletedAtIsNil()).
			SetAccessIP(access_ip).
			SetLastOperationTime(time.Now().Unix()).
			SetOperationType(common.OperationLogin).
			SetOs(ua.OS()).
			SetPlatform(ua.Platform()).
			SetBrowserName(browserName).
			SetBrowserVersion(browserVersion).
			SetBrowserEngineName(browserEngineName).
			SetBrowserEngineVersion(browserEngineVersion).
			Save(ctx)
		if terr != nil {
			global.Logger.Sugar().Errorf("用户 %s 登录更新在线用户失败: %v", req.Username, terr)
			return nil, &code.Failed
		}
	} else {
		_, terr = tx.CoreOnLineUser.Create().
			SetUserID(u.ID).
			SetAccessIP(access_ip).
			SetLastOperationTime(time.Now().Unix()).
			SetOperationType(common.OperationLogin).
			SetOs(ua.OS()).
			SetPlatform(ua.Platform()).
			SetBrowserName(browserName).
			SetBrowserVersion(browserVersion).
			SetBrowserEngineName(browserEngineName).
			SetBrowserEngineVersion(browserEngineVersion).
			Save(ctx)
		if terr != nil {
			global.Logger.Sugar().Errorf("用户 %s 登录记录在线用户失败: %v", req.Username, terr)
			return nil, &code.Failed
		}
	}

	if _, terr = tx.CoreOperationLog.Create().
		SetUserID(u.ID).
		SetAccessIP(access_ip).
		SetOperationTime(time.Now().Unix()).
		SetOperationType(common.OperationLogin).
		SetOs(ua.OS()).
		SetPlatform(ua.Platform()).
		SetBrowserName(browserName).
		SetBrowserVersion(browserVersion).
		SetBrowserEngineName(browserEngineName).
		SetBrowserEngineVersion(browserEngineVersion).
		Save(ctx); terr != nil {
		global.Logger.Sugar().Errorf("Failed to create operation log: %v", terr)
		return nil, &code.Failed
	}

	if terr = tx.Commit(); terr != nil {
		global.Logger.Sugar().Errorf("Failed to commit transaction: %v", terr)
		return nil, &code.Failed
	}

	global.Logger.Sugar().Infof("用户 %s 登录成功", u.Username)

	resp.Username = u.Username
	resp.Nickname = u.Nickname
	resp.Token = token
	resp.RoleName = u.Edges.UserFromRole.Name

	return &resp, nil
}

func (s *SystemSvc) Logout(token string, ctx context.Context) error {

	claims, perr := global.JwtToolEntity.ParseToken(token, global.Cfg.Core.Jwt.Sign)
	if perr != nil {
		global.Logger.Sugar().Errorf("用户登出解析token失败: %v", perr)
		return &code.Failed
	}
	if perr := global.JwtToolEntity.DeleteToken(claims.UserId, global.RedisCli); perr != nil {
		global.Logger.Sugar().Errorf("用户登出删除token失败: %v", perr)
		return &code.Failed
	}

	_, uerr := global.EntClient.CoreOnLineUser.Update().
		Where(coreonlineuser.UserIDEQ(claims.UserId), coreonlineuser.DeletedAtIsNil()).
		SetDeletedAt(time.Now()).
		Save(ctx)
	if uerr != nil {
		global.Logger.Sugar().Errorf("用户登出更新在线用户失败: %v", uerr)
		return &code.Failed
	}
	return nil
}
