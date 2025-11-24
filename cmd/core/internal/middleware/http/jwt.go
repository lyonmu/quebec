package http

import (
	"github.com/gin-gonic/gin"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/pkg/code"
	"github.com/lyonmu/quebec/pkg/constant"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(constant.ApiTokenName)
		if token == "" {
			global.Logger.Sugar().Errorf("jwt auth failed: token is empty")
			code.Unauthorized.Unauthorized(c)
			c.Abort()
			return
		}
		claims, err := global.JwtToolEntity.ParseToken(token, global.Cfg.Core.Jwt.Sign)
		if err != nil {
			global.Logger.Sugar().Errorf("jwt auth failed: parse token failed: %v", err)
			code.Unauthorized.Unauthorized(c)
			c.Abort()
			return
		}
		if global.JwtToolEntity.IsExpired(token, global.Cfg.Core.Jwt.Sign) {
			global.Logger.Sugar().Errorf("jwt auth failed: token is expired")
			code.Unauthorized.Unauthorized(c)
			c.Abort()
			return
		}
		storedToken, err := global.JwtToolEntity.GetToken(claims.UserId, global.RedisCli)
		if err != nil || storedToken == "" || storedToken != token {
			global.Logger.Sugar().Errorf("jwt auth failed: token mismatch")
			code.Unauthorized.Unauthorized(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
