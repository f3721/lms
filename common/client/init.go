package client

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func ReadConfig(key string) string {
	work, _ := os.Getwd()
	//设置文件名和文件后缀
	viper.SetConfigName("api")
	viper.SetConfigType("yml")
	//配置文件所在的文件夹
	viper.AddConfigPath(work + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		return ""
	}
	api := viper.GetString(key)
	return api
}

func ReadConfigByPrefix(prefix, module string) string {
	key := fmt.Sprintf("%s.%s", prefix, module)

	work, _ := os.Getwd()
	//设置文件名和文件后缀
	viper.SetConfigName("api")
	viper.SetConfigType("yml")
	//配置文件所在的文件夹
	viper.AddConfigPath(work + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		return ""
	}
	api := viper.GetString(key)
	return api
}
