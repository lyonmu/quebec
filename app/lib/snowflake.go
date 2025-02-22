package lib

import (
	"strconv"
	"sync"
	"time"
	
	"github.com/sony/sonyflake"
	
	commonvariable "github.com/lyonmu/quebec/app/common/variable"
)

var idInstance *sonyflake.Sonyflake

func MachineID() (uint16, error) {
	return commonvariable.Distributed.ID, nil
}

var once = new(sync.Once)

// InitDistributedId 初始化分布式雪花ID
func InitDistributedId() {
	once.Do(func() {
		settings := sonyflake.Settings{
			MachineID: MachineID,
			StartTime: time.Now(),
		}
		idInstance = sonyflake.NewSonyflake(settings) // 创建snowflake实例
	})
}

// ID 生成全局唯一ID
func ID() string {
	v, _ := idInstance.NextID()
	return strconv.FormatUint(v, 10)
}
