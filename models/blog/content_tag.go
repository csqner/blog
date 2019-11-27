package blog

import (
	"blog/models/base"
)

/*
  @Author : lanyulei
*/

type ContentTag struct {
	base.Model
	Content int `gorm:"column:content_id;type:int(11)" json:"content_id" form:"content_id"`
	Tag     int `gorm:"column:tag_id;type:int(11)" json:"tag_id" form:"tag_id"`
}

func (ContentTag) TableName() string {
	return "blog_content_tag"
}
