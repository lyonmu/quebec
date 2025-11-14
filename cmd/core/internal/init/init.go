package init

import (
	"sync"

	"github.com/lyonmu/quebec/cmd/core/internal/global"
)

var once sync.Once

func Init() {

	once.Do(func() {
		InitMySQL(global.Cfg.MySQL)
		global.Redis = global.Cfg.Redis.Client()
	})

}
