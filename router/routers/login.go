package routers

import (
	"blog/handler/blog"
	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

func LoginRouter(g *gin.Engine) {
	g.GET("/toLogin", blog.ToLoginHandler)
	g.GET("/login/callback", blog.CallbackHandler)
	g.POST("/userInfo", blog.UserInfoHandler)
	g.POST("/isLogin", blog.IsLoginHandler)
}
