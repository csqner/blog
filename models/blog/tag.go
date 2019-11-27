package blog

import (
	"blog/models/base"
)

/*
  @Author : lanyulei
*/

type Tag struct {
	base.Model
	Name string `gorm:"column:name;type:varchar(256)" json:"name" form:"name"`
}

func (Tag) TableName() string {
	return "blog_tag"
}
