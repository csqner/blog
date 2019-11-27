package routers

import (
	"blog/handler/blog"
	"fmt"
	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

func BlogRouter(g *gin.Engine) {
	blogRouterGroup := fmt.Sprintf("%s", "/blog")
	blogs := g.Group(blogRouterGroup)
	{
		blogs.GET("/list", blog.ListHandler)
		blogs.POST("", blog.SaveContentHandler)
	}
}
