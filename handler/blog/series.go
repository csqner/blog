package blog

import (
	"blog/models/blog"
	"blog/pkg/connection"
	. "blog/pkg/response"
	"blog/pkg/tree"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"strconv"
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

	detailsId := c.DefaultQuery("details_id", "0")
	searchValue := c.DefaultQuery("searchValue", "")

	// 获取左侧菜单
	searchValueDB := connection.DB.Self
	if searchValue != "" {
		searchValueDB = searchValueDB.Where("title like ?", fmt.Sprintf("%%%s%%", searchValue))
	}
	var seriesDetailsList []blog.SeriesDetailsStruct
	err := searchValueDB.Table("blog_series_details").
		Select("id, title, parent, order_id").
		Order("id, order_id").
		Where("series_id = ?", seriesId).Scan(&seriesDetailsList).Error
	if err != nil {
		HtmlResponse(c, "error.html", err.Error(), "/blog/series")
		return
	}

	treeValue := tree.SuperSeriesTree(seriesDetailsList, 0)

	buf := bytes.NewBufferString("")
	detailsIdInt, _ := strconv.Atoi(detailsId)
	// 获取默认展示节点
	var defaultSeriesValue blog.SeriesDetailsStruct
	if detailsId == "0" {
		if len(seriesDetailsList) > 0 {
			detailsIdInt = seriesDetailsList[0].Id
			defaultSeriesValue.Id = seriesDetailsList[0].Id
		}
	}
	var onDetailsValue blog.Tree
	var underDetailsValue blog.Tree
	onAndunder := map[string]blog.Tree{
		"onContent":    onDetailsValue,
		"underContent": underDetailsValue,
	}
	tree.GetDocumentTree(treeValue, buf, detailsIdInt, onAndunder, seriesId)

	// 文档信息
	var seriesValue blog.Series
	err = connection.DB.Self.Model(&seriesValue).Where("id = ?", seriesId).Find(&seriesValue).Error
	if err != nil {
		HtmlResponse(c, "error.html", err.Error(), "/blog/index")
		return
	}

	// 获取详细内容
	var seriesDetailsValue blog.SeriesDetails
	err = connection.DB.Self.Model(&seriesDetailsValue).Where("id = ?", detailsIdInt).Find(&seriesDetailsValue).Error
	if err != nil {
		HtmlResponse(c, "error.html", err.Error(), "/blog/series")
		return
	}

	// 更新阅读量
	err = connection.DB.Self.Model(&seriesDetailsValue).Where("id = ?", detailsIdInt).Update("view", seriesDetailsValue.View+1).Error
	if err != nil {
		HtmlResponse(c, "error.html", err.Error(), "/blog/series")
		return
	}

	// 获取文档评论
	var commentList []blog.CommentStruct
	err = connection.DB.Self.Table("blog_comment").
		Joins("left join blog_user on blog_comment.user_id = blog_user.id").
		Select("blog_comment.*, blog_user.nickname, blog_user.avatar, blog_user.title as user_title").
		Where("blog_comment.deleted_at is null and blog_comment.article_id = ? and type = 2", detailsIdInt).
		Order("created_at desc").
		Scan(&commentList).Error
	if err != nil {
		HtmlResponse(c, "error.html", fmt.Sprintf("获取评论错误，%v", err.Error()), "/blog/list")
		return
	}

	var commentIdList []int
	err = connection.DB.Self.Model(&blog.Comment{}).Where("article_id = ? and type = 2", detailsIdInt).
		Pluck("id", &commentIdList).Error
	if err != nil {
		HtmlResponse(c, "error.html", fmt.Sprintf("获取评论错误，%v", err.Error()), "/blog/list")
		return
	}

	// 获取当前文章的所有回复
	var replyList []blog.ReplyStruct
	err = connection.DB.Self.Table("blog_reply").
		Joins("left join blog_user as s_blog_user on s_blog_user.id = blog_reply.source_user_id").
		Joins("left join blog_user as d_blog_user on d_blog_user.id = blog_reply.aims_user_id").
		Select("blog_reply.*, s_blog_user.id as s_user_id, s_blog_user.nickname as s_nickname, s_blog_user.title as s_user_title, s_blog_user.avatar as s_avatar, d_blog_user.id as d_user_id, d_blog_user.nickname as d_nickname, d_blog_user.title as d_user_title, d_blog_user.avatar as d_avatar").
		Order("created_at").
		Where("blog_reply.deleted_at is null and blog_reply.comment_id in (?)", commentIdList).
		Scan(&replyList).Error
	if err != nil {
		HtmlResponse(c, "error.html", fmt.Sprintf("获取评论错误，%v", err.Error()), "/blog/list")
		return
	}

	for _, reply := range replyList {
		for commentIndex, comment := range commentList {
			if reply.CommentId == comment.Id {
				commentList[commentIndex].ReplyList = append(commentList[commentIndex].ReplyList, reply)
			}
		}
	}

	c.HTML(http.StatusOK, "SeriesDetails.html", gin.H{
		"titleTree":          template.HTML(buf.String()),
		"seriesValue":        seriesValue,
		"seriesDetailsValue": seriesDetailsValue,
		"onDetailsValue":     onAndunder["onContent"],
		"underDetailsValue":  onAndunder["underContent"],
		"comment":            commentList,
	})
}
