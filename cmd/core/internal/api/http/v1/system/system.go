package system

import "github.com/lyonmu/quebec/cmd/core/internal/service/http/system"

type SystemV1ApiGroup struct{}

var (
	systemsvc = system.SystemSvc{}
)

// GetSystemSvc 返回系统服务实例
func GetSystemSvc() *system.SystemSvc {
	return &systemsvc
}
