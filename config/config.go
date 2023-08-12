package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"monitor/global"
)

func InitConfig() {
	config := viper.New()
	config.SetConfigFile("config.yaml")
	config.SetConfigType("yaml")
	err := config.ReadInConfig()
	if err != nil {
		panic(any("读取配置失败:" + err.Error()))
	}
	config.WatchConfig()
	config.OnConfigChange(func(e fsnotify.Event) {
		global.Log(global.LOG_INFO, "config file changed:"+e.Name)
		if err = config.Unmarshal(&global.Config); err != nil {
			global.Log(global.LOG_ERROR, "配置修改失败："+err.Error())
		}
	})
	if err = config.Unmarshal(&global.Config); err != nil {
		global.Log(global.LOG_ERROR, "配置解析失败："+err.Error())
	}
}
