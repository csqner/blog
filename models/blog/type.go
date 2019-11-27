package blog

import (
	"blog/models/base"
)

/*
  @Author : lanyulei
*/

type Type struct {
	base.Model
	Name string `gorm:"column:name;type:varchar(256)" json:"name" form:"name"`
}

func (Type) TableName() string {
	return "blog_type"
}
