package inits

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const (
	ConfigFile = "config/config.yaml"
)

func Viper() *viper.Viper {
	v := viper.New()
	v.SetConfigFile(ConfigFile)
	err2 := v.ReadInConfig()
	if err2 != nil {
		panic(fmt.Errorf("viper配置文件出错: %s \n", err2.Error()))
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("viper配置文件名:", e.Name)
	})
	fmt.Println(`============================配置文件路径为: ` + ConfigFile + "\n")
	return v
}
