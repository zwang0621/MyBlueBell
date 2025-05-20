package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Init() (err error) {
	viper.SetConfigFile("conf.yaml")
	viper.AddConfigPath(".") // 指定查找配置文件的路径（这里使用相对路径）
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadInConfig() failed,err:%v\n", err)
		return err
	}

	//监控config文件
	viper.WatchConfig()
	//配置文件发生改变时会调用回调函数（配置热加载）
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("Config file changed:", in.Name)
	})
	return
}
