package blog

import "blog/models/base"

/*
  @Author : lanyulei
  @Desc : 回复表
*/

type Reply struct {
	base.Model
	CommentId    int    `gorm:"column:comment_id;type:int(11)" json:"comment_id" form:"comment_id"` // 评论ID
	Type         int    `gorm:"column:type;type:int(11)" json:"type" form:"type"`                   // 文章回复：1，留言回复：2
	LeaveId      int    `gorm:"column:leave_id;type:int(11)" json:"leave_id" form:"leave_id"`       // 留言ID
	SourceUserId int    `gorm:"column:source_user_id;type:int(11)" json:"source_user_id" form:"source_user_id"`
	AimsUserId   int    `gorm:"column:aims_user_id;type:int(11)" json:"aims_user_id" form:"aims_user_id"`
	Content      string `gorm:"column:content;type:varchar(1024)" json:"content" form:"content"`
	Browser      string `gorm:"column:browser;type:varchar(512)" json:"browser" form:"browser"`
}

func (Reply) TableName() string {
	return "blog_reply"
}
