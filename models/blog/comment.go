package blog

import "blog/models/base"

/*
  @Author : lanyulei
*/

type Comment struct {
	base.Model
	ArticleId int    `gorm:"column:article_id;type:int(11)" json:"article_id" form:"article_id"`
	Type      int    `gorm:"column:type;type:int(11)" json:"type" form:"type"` // 文章类型，博客文章 1， 系列文章 2， 留言数据，3
	UserId    int    `gorm:"column:user_id;type:int(11)" json:"user_id" form:"user_id"`
	Content   string `gorm:"column:content;type:varchar(1024)" json:"content" form:"content"`
	Browser   string `gorm:"column:browser;type:varchar(512)" json:"browser" form:"browser"`
}

func (Comment) TableName() string {
	return "blog_comment"
}

type CommentStruct struct {
	Comment
	Nickname  string        `json:"nickname"`
	UserTitle string        `json:"user_title"`
	Avatar    string        `json:"avatar"`
	ReplyList []ReplyStruct `json:"reply_list"`
}
