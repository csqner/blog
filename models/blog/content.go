package blog

import (
	"blog/models/base"
	"blog/utils"
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
	Awesome      int    `gorm:"column:awesome;int(11)" json:"awesome" form:"awesome"`  //  赞
	View         int    `gorm:"column:view;int(11);default:0" json:"view" form:"view"` // 访问量
	Image        string `gorm:"column:image;varchar(1024);" json:"image" form:"image"` // 文章图片
}

func (Content) TableName() string {
	return "blog_content"
}

type ContentListStruct struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
}

type ArticleDataStruct struct {
	ArchiveDataString string              `json:"archive_data_string"`
	ArticleData       []ContentListStruct `json:"article_data"`
}

type ArchiveData struct {
	TypeCount      int                 `json:"type_count"`       // 分类总数
	AwesomeCount   int                 `json:"awesome_count"`    // 总赞数
	ArticleCount   int                 `json:"article_count"`    // 文章总数
	CommentCount   int                 `json:"commend_count"`    // 评论总数
	ViewCount      int                 `json:"view_count"`       // 文章总访问量
	LastUpdateTime utils.JSONTime      `json:"last_update_time"` // 最后更新时间
	ArticleGroup   []ArticleDataStruct `json:"article_group"`
}

type DateStruct struct {
	Year   int `json:"year"`
	Months []MonthStruct
}

type MonthStruct struct {
	Month int `json:"month"`
}
