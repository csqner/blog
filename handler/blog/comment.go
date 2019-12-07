package blog

import (
	"blog/models/blog"
	"blog/models/user"
	"blog/pkg/connection"
	"blog/pkg/errno"
	"blog/pkg/login"
	. "blog/pkg/response"
	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

// 保存评论
func SaveCommentHandler(c *gin.Context) {
	var params blog.Comment
	if err := c.Bind(&params); err != nil {
		Response(c, errno.ErrBind, nil, err.Error())
		return
	}

	openId := login.GetSessionOpenId(c)
	var userInfo user.User
	err := connection.DB.Self.Model(&user.User{}).Where("open_id = ?", openId).Find(&userInfo).Error
	if err != nil {
		Response(c, errno.ErrSelectUser, nil, err.Error())
		return
	}

	params.UserId = userInfo.Id
	err = connection.DB.Self.Model(&blog.Comment{}).Create(&params).Error
	if err != nil {
		Response(c, errno.ErrCreateComment, nil, err.Error())
		return
	}

	Response(c, nil, nil, "")
}

// 保存回复
func SaveReplyHandler(c *gin.Context) {
	var params blog.Reply
	if err := c.Bind(&params); err != nil {
		Response(c, errno.ErrBind, nil, err.Error())
		return
	}

	openId := login.GetSessionOpenId(c)
	var userInfo user.User
	err := connection.DB.Self.Model(&user.User{}).Where("open_id = ?", openId).Find(&userInfo).Error
	if err != nil {
		Response(c, errno.ErrSelectUser, nil, err.Error())
		return
	}

	params.SourceUserId = userInfo.Id
	err = connection.DB.Self.Model(&blog.Reply{}).Create(&params).Error
	if err != nil {
		Response(c, errno.ErrCreateComment, nil, err.Error())
		return
	}

	Response(c, nil, nil, "")
}
