package system

import (
	"context"
	"time"

	"github.com/lyonmu/quebec/cmd/core/internal/common"
	"github.com/lyonmu/quebec/cmd/core/internal/dto/request"
	"github.com/lyonmu/quebec/cmd/core/internal/dto/response"
	"github.com/lyonmu/quebec/cmd/core/internal/ent/coreuser"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/pkg/code"
	"github.com/lyonmu/quebec/pkg/tools/encrypt"
	"github.com/mssola/useragent"
)

func (s *SystemSvc) Login(req *request.SystemLoginRequest, ua *useragent.UserAgent, access_ip string, ctx context.Context) (*response.SystemInfoResponse, error) {
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
		Where(coreuser.UsernameEQ(req.Username)).
		Where(coreuser.DeletedAtIsNil()).
		First(ctx)

	if qerr != nil {
		global.Logger.Sugar().Errorf("用户登录查询数据库失败:", qerr)
		return nil, &code.UserNotExists
	}
	if encrypt.CompareWithBcryptString(u.Password, req.Password) {
		global.Logger.Sugar().Errorf("用户 %s 登录密码错误", req.Username)
		return nil, &code.UserPasswordError
	}

	token, terr := global.JwtToolEntity.GenerateToken(u.ID, u.Username, u.LastPasswordChange, global.Cfg.Core.Jwt.Cache)
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

	// TODO 如果已经存在在线用户记录,则更新记录,否则创建记录
	_, oerr := global.EntClient.CoreOnLineUser.Create().
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
	if oerr != nil {
		global.Logger.Sugar().Errorf("用户 %s 登录记录在线用户失败: %v", req.Username, oerr)
		return nil, &code.Failed
	}

	global.Logger.Sugar().Infof("用户 %s 登录成功", u.Username)

	resp.Username = u.Username
	resp.Token = token

	return &resp, nil
}
