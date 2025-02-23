package v1

import "github.com/lyonmu/quebec/app/base/service/v1/user"

type ServiceV1Group struct {
	UserService user.ServiceV1Group
}

var ServiceV1GroupInstance = new(ServiceV1Group)
