package user

import servicev1 "github.com/lyonmu/quebec/app/base/service/v1"

type ApiV1Group struct {
	BasicUserOperationApiV1
}

var (
	userService = servicev1.ServiceV1GroupInstance.UserService
)
