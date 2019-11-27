package blog

import (
	"blog/models/base"
)

/*
  @Author : lanyulei
*/

type Content struct {
	base.Model
	Title        string `gorm:"column:title;type:varchar(512)" json:"title" form:"title"`
	Introduction string `grom:"column:introduction;type:varchar(1024)" json:"introduction" form:"introduction"`
	Type         int    `gorm:"column:type_id;type:int(11)" json:"type_id" form:"type_id"`
	Content      string `gorm:"column:content;type:longtext" json:"content" form:"content"`
	Awesome      int    `gorm:"column:awesome;int(11)" json:"awesome" form:"awesome"` //  赞
	View         int    `gorm:"column:view;int(11)" json:"view" form:"view"`          // 访问量
}

func (Content) TableName() string {
	return "blog_content"
}
