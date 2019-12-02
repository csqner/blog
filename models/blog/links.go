package blog

import "blog/models/base"

/*
  @Author : lanyulei
  @Desc : 友情链接
*/

type Links struct {
	base.Model
	Name        string `gorm:"column:name;type:varchar(256)" json:"name" form:"name"`
	Address     string `gorm:"column:address;type:varchar(512)" json:"address" form:"address"`
	Icon        string `gorm:"column:icon;type:varchar(512)" json:"icon" form:"icon"`
	Description string `gorm:"column:description;type:varchar(1024)" json:"description" form:"description"`
}

func (Links) TableName() string {
	return "blog_links"
}
