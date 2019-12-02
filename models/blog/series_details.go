package blog

import "blog/models/base"

/*
  @Author : lanyulei
*/

type SeriesDetails struct {
	base.Model
	Series  int    `gorm:"column:series_id;int(11);" json:"series_id" form:"series_id"`
	Title   string `gorm:"column:title;type:varchar(512)" json:"title" form:"title"`
	Content string `gorm:"column:content;type:longtext" json:"content" form:"content"`
	Awesome int    `gorm:"column:awesome;int(11)" json:"awesome" form:"awesome"`     // 赞
	View    int    `gorm:"column:view;int(11);default:0" json:"view" form:"view"`    // 访问量
	Parent  int    `gorm:"column:parent;int(11);" json:"parent" form:"parent"`       // 父ID
	Order   int    `gorm:"column:order_id;int(11);" json:"order_id" form:"order_id"` // 排序
}

func (SeriesDetails) TableName() string {
	return "blog_series_details"
}
