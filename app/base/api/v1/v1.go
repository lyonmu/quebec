package v1

import (
	"github.com/lyonmu/quebec/app/base/api/v1/system"
	"github.com/lyonmu/quebec/app/base/api/v1/user"
)

// @title                       quebec base API接口文档
// @description                 quebec base 模块后端
// @version                     v1
// @license.name  MIT
// @license.url   https://github.com/lyonmu/quebec/blob/master/LICENSE

type ApiV1Group struct {
	UserApi   user.ApiV1Group
	SystemApi system.ApiV1Group
}

var ApiV1GroupInstance = new(ApiV1Group)
