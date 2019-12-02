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
		blogs.GET("/index", blog.IndexHandler)
		blogs.GET("/list", blog.ListHandler)
		blogs.GET("/archive", blog.ArchiveHandler)
		blogs.GET("/about", blog.AboutHandler)
		blogs.GET("/author", blog.AuthorHandler)
		blogs.GET("/links", blog.LinksHandler)
		blogs.GET("/series", blog.SeriesHandler)
		blogs.GET("/series_details", blog.SeriesDetailsHandler)
		blogs.POST("", blog.SaveContentHandler)
		blogs.GET("/details", blog.ContentDetailsHandler)
	}
}
