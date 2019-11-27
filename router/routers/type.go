package routers

import (
	"blog/handler/blog"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

/*
  @Author : lanyulei
*/

func TypeRouter(g *gin.Engine) {
	typeRouterGroup := fmt.Sprintf("%s%s", viper.GetString(`api.version`), "/type")

	// user
	//blogs := g.Group(blogRouterGroup, middleware.AuthMiddleware())
	types := g.Group(typeRouterGroup)
	{
		types.GET("", blog.TypeListHandler)
		types.POST("", blog.CreateTypeHandler)
		types.POST("/update", blog.UpdateTypeHandler)
	}
}
