/*
  @Author : lanyulei
*/

package main

import (
	_ "blog/conf"
	"blog/models"
	"blog/pkg/connection"
	"blog/pkg/logger"
	"blog/router"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func init() {
	// 监控配置文件变化
	viper.OnConfigChange(func(e fsnotify.Event) {
		logger.Infof("Config file changed: %s", e.Name)
	})

	// 初始化数据库连接
	connection.DB.Init()
}

func main() {
	// 同步表结构
	models.AutoMigrateTable()

	g := gin.New()
	g.LoadHTMLGlob("template/*")
	g.Static("/static", "./static")

	// 加载路由
	router.Load(g)

	// 运行程序
	err := g.Run(":9090")
	if err != nil {
		logger.Error("启动失败")
		panic(err)
	}

	// 关闭数据库连接
	defer connection.DB.Close()
}
