package routers

import (
	"blog/handler/blog"
	"blog/middleware"
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
		blogs.POST("/save-comment", middleware.AuthSessionMiddle(), blog.SaveCommentHandler)
		blogs.POST("/save-reply", middleware.AuthSessionMiddle(), blog.SaveReplyHandler)
		blogs.GET("/add-awesome", blog.AddAwesomeHandler)
		blogs.GET("/leave", blog.LeaveHandler)
	}
}
