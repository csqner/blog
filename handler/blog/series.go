package blog

import (
	"blog/models/blog"
	"blog/pkg/connection"
	. "blog/pkg/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
  @Author : lanyulei
*/

// 系列教程
func SeriesHandler(c *gin.Context) {
	var seriesList []blog.Series
	err := connection.DB.Self.Model(&blog.Series{}).Find(&seriesList).Error
	if err != nil {
		HtmlResponse(c, "error.html", err.Error(), "/blog/index")
		return
	}

	c.HTML(http.StatusOK, "Series.html", gin.H{
		"series": seriesList,
	})
}

// 文档详情
func SeriesDetailsHandler(c *gin.Context) {
	seriesId := c.DefaultQuery("series_id", "")
	if seriesId == "" {
		HtmlResponse(c, "error.html", "series_id参数不存在，请确认", "/blog/series")
		return
	}

	var seriesDetailsList []struct {
		Id     int    `json:"id"`
		Title  string `json:"title"`
		Parent int    `json:"parent"`
		Order  int    `json:"order_id"`
	}

	err := connection.DB.Self.Table("blog_series_details").
		Select("id, title, parent, order_id").
		Where("series_id = ?", seriesId).Scan(&seriesDetailsList).Error
	if err != nil {
		HtmlResponse(c, "error.html", err.Error(), "/blog/series")
		return
	}

	fmt.Println(seriesDetailsList)

	c.HTML(http.StatusOK, "SeriesDetails.html", gin.H{})
}
