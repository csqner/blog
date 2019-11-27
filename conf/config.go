/*
  @Author : lanyulei
*/

package conf

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// config path
var configDir = "conf/file"

type Config struct {
	Name string
}

func init() {
	// 全局初始化viper
	cfg := pflag.StringP("config", "c", "", "配置文件的路径")
	pflag.Parse()
	err := Init(*cfg)
	if err != nil {
		panic(err)
	}
}

func Init(cfg string) error {
	c := Config{
		Name: cfg,
	}
	if err := c.initConfig(); err != nil {
		return err
	}
	return nil
}

func (c *Config) initConfig() error {
	if c.Name != "" {
		viper.SetConfigFile(c.Name)
	} else {
		viper.AddConfigPath(configDir)
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")

	// 热加载配置信息
	viper.WatchConfig()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
