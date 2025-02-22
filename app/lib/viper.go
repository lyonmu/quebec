package lib

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// LoadConfig 加载并监听配置文件，实现文件内容热加载
func LoadConfig(m string, conf any) error {
	viper.SetConfigFile(m)
	viper.WatchConfig()
	var resp error
	viper.OnConfigChange(func(in fsnotify.Event) {
		if err := viper.Unmarshal(conf); err != nil {
			resp = err
		}
	})
	if err := viper.ReadInConfig(); err != nil {
		resp = err
	}
	if err := viper.Unmarshal(conf); err != nil {
		resp = err
	}
	return resp
}
