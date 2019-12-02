/*
  @Author : lanyulei
*/

package router

import (
	"blog/router/routers"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// 加载路由
func Load(g *gin.Engine) {
	// 404
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "404 页面不存在")
	})

	// pprof router
	pprof.Register(g)

	// cors， 跨域
	config := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowOrigins:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	g.Use(cors.New(config))

	// blog
	routers.BlogRouter(g)

	// blog type
	routers.TypeRouter(g)
}
