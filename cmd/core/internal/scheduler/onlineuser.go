package scheduler

import (
	"context"
	"fmt"

	"github.com/lyonmu/quebec/cmd/core/internal/common"
	"github.com/lyonmu/quebec/cmd/core/internal/ent/coreonlineuser"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
)

// 根据redis中存储的在线用户信息，删除数据库中已过期的在线用户记录
func DelOnlineuserTask() {

	var (
		ctx = context.Background()
	)
	// 1、获取数据库中所有在线用户记录
	u, qerr := global.EntClient.CoreOnLineUser.Query().
		Select(coreonlineuser.FieldUserID).
		Where(coreonlineuser.DeletedAtIsNil()).
		All(ctx)
	if qerr != nil {
		global.Logger.Sugar().Errorf("删除在线用户查询数据库失败: %v", qerr)
		return
	}
	// 2、遍历每条记录，检查其在Redis中是否存在
	for _, user := range u {
		exists, err := global.RedisCli.Exists(ctx, fmt.Sprintf(common.TokenCache, user.UserID)).Result()
		if err != nil {
			global.Logger.Sugar().Errorf("删除在线用户查询Redis失败: %v", err)
			continue
		}
		if exists == 0 {
			_, derr := global.EntClient.CoreOnLineUser.
				Delete().
				Where(coreonlineuser.UserIDEQ(user.UserID)).
				Exec(ctx)
			if derr != nil {
				global.Logger.Sugar().Errorf("删除在线用户数据库记录失败: %v", derr)
				continue
			}
		}
	}

	_, derr := global.EntClient.CoreOnLineUser.Delete().Where(coreonlineuser.DeletedAtNotNil()).Exec(ctx)
	if derr != nil {
		global.Logger.Sugar().Errorf("删除已过期在线用户数据库记录失败: %v", derr)
	}

}
