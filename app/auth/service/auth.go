package service

import (
	"context"
	"strings"

	"github.com/sirupsen/logrus"

	envoyapiv3core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoyserviceauthv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"
)

type Users map[string]string

type server struct {
	users Users
}

// Check checks if a key could retrieve a user from a list of users.
func (u Users) Check(key string) (bool, string) {
	value, ok := u[key]
	if !ok {
		return false, ""
	}
	return ok, value
}

var _ envoyserviceauthv3.AuthorizationServer = &server{}

// New creates a new authorization server.
func New(users Users) envoyserviceauthv3.AuthorizationServer {
	return &server{users}
}

// Check implements authorization's Check interface which performs authorization check based on the
// attributes associated with the incoming request.
func (s *server) Check(
	ctx context.Context,
	req *envoyserviceauthv3.CheckRequest) (*envoyserviceauthv3.CheckResponse, error) {
	authorization := req.Attributes.Request.Http.Headers["authorization"]
	logrus.Info("authorization: ", authorization)

	extracted := strings.Fields(authorization)
	if len(extracted) == 2 && extracted[0] == "Bearer" {
		valid, user := s.users.Check(extracted[1])
		if valid {
			return &envoyserviceauthv3.CheckResponse{
				HttpResponse: &envoyserviceauthv3.CheckResponse_OkResponse{
					OkResponse: &envoyserviceauthv3.OkHttpResponse{
						Headers: []*envoyapiv3core.HeaderValueOption{
							{
								Append: &wrappers.BoolValue{Value: false},
								Header: &envoyapiv3core.HeaderValue{
									// For a successful request, the authorization server sets the
									// x-current-user value.
									Key:   "x-current-user",
									Value: user,
								},
							},
						},
					},
				},
				Status: &status.Status{
					Code: int32(code.Code_OK),
				},
			}, nil
		}
	}

	return &envoyserviceauthv3.CheckResponse{
		Status: &status.Status{
			Code: int32(code.Code_PERMISSION_DENIED),
		},
	}, nil
}
