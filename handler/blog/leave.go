package blog

import (
	"blog/models/base"
	"blog/models/blog"
	"blog/pkg/connection"
	"blog/pkg/errno"
	. "blog/pkg/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
  @Author : lanyulei
*/

func LeaveHandler(c *gin.Context) {
	var param base.ListRequest
	if err := c.Bind(&param); err != nil {
		Response(c, errno.ErrBind, nil, err.Error())
		return
	}

	var err error

	if param.Offset > 0 {
		param.Offset = param.Offset - 1
	}

	if param.Limit == 0 {
		param.Limit = 5
	}

	type replyStruct struct {
		blog.Reply
		SUserId    int    `json:"s_user_id"`
		SNickname  string `json:"s_nickname"`
		SUserTitle string `json:"s_user_title"`
		SAvatar    string `json:"s_avatar"`
		NUserId    int    `json:"n_user_id"`
		DNickname  string `json:"d_nickname"`
		DUserTitle string `json:"d_user_title"`
		DAvatar    string `json:"d_avatar"`
	}

	type commentStruct struct {
		blog.Comment
		Nickname  string        `json:"nickname"`
		UserTitle string        `json:"user_title"`
		Avatar    string        `json:"avatar"`
		ReplyList []replyStruct `json:"reply_list"`
	}
	var commentList []commentStruct
	err = connection.DB.Self.Table("blog_comment").
		Joins("left join blog_user on blog_comment.user_id = blog_user.id").
		Select("blog_comment.*, blog_user.nickname, blog_user.avatar, blog_user.title as user_title").
		Where("blog_comment.deleted_at is null and blog_comment.type = 3").
		Offset(param.Offset * param.Limit).
		Limit(param.Limit).
		Order("created_at desc").
		Scan(&commentList).Error
	if err != nil {
		HtmlResponse(c, "error.html", fmt.Sprintf("获取评论错误，%v", err.Error()), "/blog/list")
		return
	}

	var commentIdList []int
	err = connection.DB.Self.Model(&blog.Comment{}).
		Pluck("id", &commentIdList).Error
	if err != nil {
		HtmlResponse(c, "error.html", fmt.Sprintf("获取评论错误，%v", err.Error()), "/blog/list")
		return
	}

	// 获取当前文章的所有回复
	var replyList []replyStruct
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

	var dataCount int
	err = connection.DB.Self.Model(&blog.Comment{}).Where("type = 3").Count(&dataCount).Error
	if err != nil {
		HtmlResponse(c, "error.html", fmt.Sprintf("获取留言错误，%v", err.Error()), "/blog/index")
		return
	}

	c.HTML(http.StatusOK, "MessageBoard.html", gin.H{
		"comment":   commentList,
		"dataCount": dataCount,
	})
}
