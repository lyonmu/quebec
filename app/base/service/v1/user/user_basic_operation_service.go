package user

import (
	"context"
	commonresp "github.com/lyonmu/quebec/app/common/resp"
)

type BasicUserOperationServiceV1 struct{}

func (BasicUserOperationServiceV1) Login(ctx context.Context) *commonresp.HttpResponse {
	return &commonresp.Success
}

func (BasicUserOperationServiceV1) Signup(ctx context.Context) *commonresp.HttpResponse {
	return &commonresp.Success
}

func (BasicUserOperationServiceV1) Captcha(ctx context.Context) *commonresp.HttpResponse {
	return &commonresp.Success
}
