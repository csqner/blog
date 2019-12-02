package blog

import "blog/models/base"

/*
  @Author : lanyulei
*/

type Series struct {
	base.Model
	Title        string `gorm:"column:title;type:varchar(512)" json:"title" form:"title"`
	Introduction string `grom:"column:introduction;type:varchar(1024)" json:"introduction" form:"introduction"`
	Awesome      int    `gorm:"column:awesome;int(11);default:0" json:"awesome" form:"awesome"` //  赞
	View         int    `gorm:"column:view;int(11);default:0" json:"view" form:"view"`          // 访问量
	Image        string `gorm:"column:image;int(11);" json:"image" form:"image"`                // 背景图片
}

func (Series) TableName() string {
	return "blog_series"
}
