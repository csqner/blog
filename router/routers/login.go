package routers

import (
	"blog/handler/blog"
	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

func LoginRouter(g *gin.Engine) {
	g.GET("/toLogin", blog.GetAuthCode)
	g.GET("/qqLogin", blog.GetToken)
}
