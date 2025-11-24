package auth

import (
	"context"
	"encoding/json"
	"net/url"

	"strings"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/lyonmu/quebec/cmd/gateway/internal/global"
	"github.com/lyonmu/quebec/pkg/constant"
	codepb "google.golang.org/genproto/googleapis/rpc/code"
	rpcstatus "google.golang.org/genproto/googleapis/rpc/status"
)

// AuthEvent 实现 AuthorizationServer 接口
type AuthEvent struct{}

func NewAuthEvent() *AuthEvent {
	return &AuthEvent{}
}

func (a *AuthEvent) Check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {

	var (
		rawPath = req.GetAttributes().GetRequest().GetHttp().GetPath()
		headers = req.GetAttributes().GetRequest().GetHttp().GetHeaders()
		method  = req.GetAttributes().GetRequest().GetHttp().GetMethod()
		path    string
		query   string
	)

	// 将原生的请求路径拆分为 path 和 query
	splitPath := strings.Split(rawPath, "?")
	switch len(splitPath) {
	case 1:
		path = splitPath[0]
		query = ""
	case 2:
		path = splitPath[0]
		query = splitPath[1]
	default:
		path = rawPath
		query = ""
	}

	b, err := json.MarshalIndent(headers, "", "  ")
	if err == nil {
		global.Logger.Sugar().Debugf("inbound headers: %s", string(b))
	}

	if global.Cfg.Gateway.WhiteList[path] == strings.ToUpper(method) {
		global.Logger.Sugar().Debugf("method: %s, path: %s is in white list", method, path)
		return &authv3.CheckResponse{Status: &rpcstatus.Status{Code: int32(codepb.Code_OK)}}, nil
	}
	// 从 header 中获取 x-quebec-token
	authHeader, ok := headers[constant.ApiTokenName]
	if !ok || len(authHeader) == 0 {
		// 从 query 参数中获取 x-quebec-token
		parsedQuery, err := url.ParseQuery(query)
		if err != nil {
			global.Logger.Sugar().Errorf("failed to parse query: %v", err)
			return &authv3.CheckResponse{
				Status: &rpcstatus.Status{
					Code: int32(codepb.Code_UNAUTHENTICATED),
				},
			}, nil
		}
		authHeader = parsedQuery.Get(constant.ApiTokenName)
	}

	global.Logger.Sugar().Debugf("method: %s, path: %s, x-quebec-token is %s", method, path, authHeader)

	if len(authHeader) != 0 {
		return &authv3.CheckResponse{Status: &rpcstatus.Status{Code: int32(codepb.Code_OK)}}, nil
	} else {
		global.Logger.Sugar().Debugf("no x-quebec-token found in header or query")
		return &authv3.CheckResponse{
			Status: &rpcstatus.Status{
				Code: int32(codepb.Code_UNAUTHENTICATED),
			},
			HttpResponse: &authv3.CheckResponse_DeniedResponse{
				DeniedResponse: &authv3.DeniedHttpResponse{
					Status: &typev3.HttpStatus{
						Code: typev3.StatusCode_Unauthorized,
					},
					Body: "未授权，请先登录",
				},
			},
		}, nil
	}
}
