package blog

import "blog/models/base"

/*
  @Author : lanyulei
*/

type Philosophy struct {
	base.Model
	Content string `gorm:"column:content;type:varchar(512)" json:"content" form:"content"`
}

func (Philosophy) TableName() string {
	return "blog_philosophy"
}
