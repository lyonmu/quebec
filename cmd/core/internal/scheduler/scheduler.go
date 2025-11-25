package scheduler

import (
	"time"

	"github.com/lyonmu/quebec/cmd/core/internal/global"
)

func StartSchedulerTask() {

	time.AfterFunc(15*time.Second, func() {
		global.Logger.Sugar().Info("启动定时任务")
		go func() {
			ticker := time.Tick(30 * time.Minute)
			for range ticker {
				DelOnlineuserTask()
			}
		}()

	})

}
