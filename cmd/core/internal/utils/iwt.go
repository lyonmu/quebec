package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lyonmu/quebec/cmd/core/internal/common"
	"github.com/lyonmu/quebec/pkg/tools/encrypt"
	"github.com/redis/go-redis/v9"
)

type JwtClaims struct {
	UserId             string `json:"user_id"`
	LastPasswordChange int64  `json:"last_password_change"`
	jwt.RegisteredClaims
}

type JwtTool struct{}

func (JwtTool) GenerateToken(id, sign string, lastPasswordChange, cache int64) (string, error) {

	// Create claims with multiple fields populated
	claims := JwtClaims{
		id,
		lastPasswordChange,
		jwt.RegisteredClaims{
			// token 的过期时间，超过此时间则不应被接受。
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cache) * time.Hour)),
			// token 的签发时间，表示 token 何时被创建。
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// token 生效的时间，在此时间之前，token 不应被接受。
			NotBefore: jwt.NewNumericDate(time.Now()),
			// 指定 token 的发行者，通常是创建并签发 token 的实体。
			Issuer: "quebec",
			// jti 唯一标识符，用于防止重放攻击。
			ID: uuid.Must(uuid.NewV7()).String(),
		},
	}
	ss, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(encrypt.HashWithSHA256String(sign)))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func (JwtTool) ParseToken(t, sign string) (*JwtClaims, error) {

	token, err := jwt.ParseWithClaims(t, &JwtClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(encrypt.HashWithSHA256String(sign)), nil
	})
	if err != nil {
		return nil, err
	} else if claims, ok := token.Claims.(*JwtClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("unknown claims type, cannot proceed")
	}
}

func (JwtTool) IsExpired(t, sign string) bool {
	token, err := jwt.ParseWithClaims(t, &JwtClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(encrypt.HashWithSHA256String(sign)), nil
	})
	if err != nil {
		return true
	} else if claims, ok := token.Claims.(*JwtClaims); ok {
		return claims.ExpiresAt.Time.Before(time.Now())
	} else {
		return true
	}
}

func (JwtTool) StoreToken(t, id string, cache int64, cli redis.UniversalClient) error {
	return cli.Set(context.Background(), fmt.Sprintf(common.TokenCache, id), t, time.Duration(cache)*time.Hour).Err()
}

func (JwtTool) DeleteToken(id string, cli redis.UniversalClient) error {
	return cli.Del(context.Background(), fmt.Sprintf(common.TokenCache, id)).Err()
}

func (JwtTool) GetToken(id string, cli redis.UniversalClient) (string, error) {
	return cli.Get(context.Background(), fmt.Sprintf(common.TokenCache, id)).Result()
}
