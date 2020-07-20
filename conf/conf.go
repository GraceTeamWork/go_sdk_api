package conf

import (
	"fmt"
	"github.com/spf13/viper"
)

func init() {

	fmt.Println("初始化配置文件")

	//设置配置文件，不包含配置文件扩展
	viper.SetConfigName("conf")
	//设置配置文件类型
	viper.SetConfigType("yaml")
	//设置配置文件目录
	viper.AddConfigPath("./conf")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func GetConfig(key string) interface{} {

	return viper.Get(key)

}
