package utils

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lyonmu/quebec/pkg/config"
)

func TestJwtGenerateToken(t *testing.T) {

	mySigningKey := []byte("quebec_secret_key")
	type MyCustomClaims struct {
		Name string `json:"name"`
		jwt.RegisteredClaims
	}

	// Create claims with multiple fields populated
	claims := MyCustomClaims{
		"test",
		jwt.RegisteredClaims{
			// token 的过期时间，超过此时间则不应被接受。
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			// token 的签发时间，表示 token 何时被创建。
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// token 生效的时间，在此时间之前，token 不应被接受。
			NotBefore: jwt.NewNumericDate(time.Now()),
			// 指定 token 的发行者，通常是创建并签发 token 的实体。
			Issuer: "quebec",
			// 指定 token 的主题，通常是 token 所代表的用户或实体。
			Subject: "user",
			// 唯一标识符，用于防止重放攻击。
			ID: "1",
			// 接收该 token 的预期受众。
			Audience: []string{"frontend", "api"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	fmt.Println(ss, err)
}

func TestJwtParseToken(t *testing.T) {

	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJpc3MiOiJ0ZXN0IiwiYXVkIjoic2luZ2xlIn0.QAWg1vGvnqRuCFTMcPkjZljXHh8U3L_qUjszOtQbeaA"

	type MyCustomClaims struct {
		Foo string `json:"foo"`
		jwt.RegisteredClaims
	}

	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (any, error) {
		return []byte("AllYourBase"), nil
	})
	if err != nil {
		log.Fatal(err)
	} else if claims, ok := token.Claims.(*MyCustomClaims); ok {
		fmt.Println(claims.Foo, claims.Issuer)
	} else {
		log.Fatal("unknown claims type, cannot proceed")
	}
}

func TestJwtTool(t *testing.T) {

	testJwt := JwtTool{}

	ss, _ := testJwt.GenerateToken("12345", "quebec", time.Now().Unix(), 24)
	log.Printf("token  is : %s", ss)

	claims, err := testJwt.ParseToken(ss, "quebec")
	if err != nil {
		log.Printf("parse token error: %v", err)
	} else {
		log.Printf("info is : %+v", claims)
	}
}

func TestIsExpired(t *testing.T) {
	testJwt := JwtTool{}
	ss := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjI2ODQzNTQ1NyIsImxhc3RfcGFzc3dvcmRfY2hhbmdlIjoxNzYzODEzODE4LCJpc3MiOiJxdWViZWMiLCJleHAiOjE3NjY0MDcyODUsIm5iZiI6MTc2MzgxNTI4NSwiaWF0IjoxNzYzODE1Mjg1fQ.XQ7xS1xub1at7HXr1ht6xrd06PxF2FuL6GqTTHcTEGs"

	expired := testJwt.IsExpired(ss, "quebec")
	log.Printf("is expired: %v", expired)
}

func TestJwtToolParseToken(t *testing.T) {
	testJwt := JwtTool{}
	ss := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNzM4MTk3NTI4IiwibGFzdF9wYXNzd29yZF9jaGFuZ2UiOjE3NjM5NDY2MzAsImlzcyI6InF1ZWJlYyIsImV4cCI6MTc2NDAzNzI1OCwibmJmIjoxNzYzOTUwODU4LCJpYXQiOjE3NjM5NTA4NTgsImp0aSI6IjAxOWFiM2FhLTEwMTAtNzcxNC04NGJlLWE1MTg2ODE0ODRhMiJ9._dRKDFrUblWQK1xQotdddnVDWxhV_j21CkMyTEZsmQU"
	claims, err := testJwt.ParseToken(ss, "quebec")
	if err != nil {
		log.Printf("parse token error: %v", err)
	} else {
		log.Printf("info is : %+v", claims)
	}
}

func TestStoreToken(t *testing.T) {
	testJwt := JwtTool{}
	rc := config.RedisConfig{
		Host:     []string{"127.0.0.1:6379"},
		Password: "root",
		DB:       1,
	}
	rcli := rc.Client("test-module")

	ss := "1231231232132"

	if err := testJwt.StoreToken(ss, "123456", 24, rcli); err != nil {
		log.Printf("store token error: %v", err)
	} else {
		log.Printf("store token success")
	}
}

func TestGetToken(t *testing.T) {
	testJwt := JwtTool{}
	rc := config.RedisConfig{
		Host:     []string{"127.0.0.1:6379"},
		Password: "root",
		DB:       1,
	}
	rcli := rc.Client("test-module")
	token, err := testJwt.GetToken("123456", rcli)
	if err != nil {
		log.Printf("get token error: %v", err)
	} else {
		log.Printf("get token success: %s", token)
	}
}

func TestDeleteToken(t *testing.T) {
	testJwt := JwtTool{}
	rc := config.RedisConfig{
		Host:     []string{"127.0.0.1:6379"},
		Password: "root",
		DB:       1,
	}
	rcli := rc.Client("test-module")
	if err := testJwt.DeleteToken("123456", rcli); err != nil {
		log.Printf("del token error: %v", err)
	} else {
		log.Printf("del token success")
	}
}
