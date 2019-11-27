package blog

import (
	"blog/models/base"
	"blog/models/blog"
	"blog/pkg/connection"
	"blog/pkg/errno"
	. "blog/pkg/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

/*
  @Author : lanyulei
*/

// 获取文章列表
func ListHandler(c *gin.Context) {
	typeParam := c.DefaultQuery("type", "")
	var param base.ListRequest
	if err := c.Bind(&param); err != nil {
		Response(c, errno.ErrBind, nil, err.Error())
		return
	}

	var err error

	if param.Offset > 0 {
		param.Offset = param.Offset - 1
	}

	if param.Limit == 0 {
		param.Limit = 5
	}

	type contentTypeValue struct {
		blog.Content
		TypeName string `json:"type_name"`
	}
	var contentList []contentTypeValue
	db := connection.DB.Self
	if typeParam != "" {
		db = db.Where("blog_content.type_id = ?", typeParam)
	}
	err = db.Table("blog_content").
		Select("blog_content.*, blog_type.name as type_name").
		Joins("left join blog_type on blog_content.type_id = blog_type.id").
		Where("blog_content.title like ?", fmt.Sprintf("%%%s%%", param.Item)).
		Offset(param.Offset * param.Limit).
		Limit(param.Limit).
		Scan(&contentList).Error

	if err != nil {
		Response(c, errno.ErrSelectList, nil, err.Error())
		return
	}

	type typeModels struct {
		Count    int    `json:"count"`
		TypeId   int    `json:"type_id"`
		TypeName string `json:"type_name"`
	}

	// 获取分类
	var typeList []typeModels
	err = connection.DB.Self.Table("blog_content").
		Joins("left join blog_type on blog_content.type_id = blog_type.id").
		Where("blog_content.type_id in (?)", connection.DB.Self.Table("blog_type").Select("id").QueryExpr()).
		Select("blog_content.type_id, blog_type.name as type_name, count(blog_content.id) as count").Group("blog_content.type_id").
		Order("count desc").Limit(5).
		Scan(&typeList).Error
	if err != nil {
		Response(c, errno.ErrSelectTypeList, nil, err.Error())
		return
	}

	// 获取所有数据的数量
	var contentCount int
	countDB := connection.DB.Self
	if typeParam != "" {
		countDB = countDB.Where("blog_content.type_id = ?", typeParam)
	}
	if param.Item != "" {
		countDB = countDB.Where("blog_content.title like ?", fmt.Sprintf("%%%s%%", param.Item))
	}
	err = countDB.Model(&blog.Content{}).Count(&contentCount).Error
	if err != nil {
		Response(c, errno.ErrSelectList, nil, err.Error())
		return
	}

	c.HTML(http.StatusOK, "Article.html", gin.H{
		"content": contentList,
		"types":   typeList,
		"count":   contentCount,
	})
}

// 新建文章
func SaveContentHandler(c *gin.Context) {
	type contentParams struct {
		Title        string `json:"title"`
		Introduction string `json:"introduction"`
		Tags         string `json:"tags"`
		TypeValue    string `json:"type"`
		Content      string `json:"content"`
	}

	var params contentParams
	err := c.ShouldBindJSON(&params)
	if err != nil {
		Response(c, errno.ErrBind, nil, err.Error())
		return
	}

	// 创建事物
	tx := connection.DB.Self.Begin()

	// 创建文章
	typeValueInt, err := strconv.Atoi(params.TypeValue)
	if err != nil {
		Response(c, errno.InternalServerError, nil, err.Error())
		return
	}

	contentValue := blog.Content{
		Title:        params.Title,
		Introduction: params.Introduction,
		Content:      params.Content,
		Type:         typeValueInt,
	}

	err = tx.Create(&contentValue).Error
	if err != nil {
		tx.Rollback()
		Response(c, errno.ErrCreateBlog, nil, err.Error())
		return
	}

	// 标签处理
	tagList := strings.Split(params.Tags, ",")
	for _, tag := range tagList {
		tag = strings.TrimSpace(tag)

		// 获取标签数量
		var tagCount int
		err = tx.Model(&blog.Tag{}).Where("name = ?", tag).Count(&tagCount).Error
		if err != nil {
			tx.Rollback()
			Response(c, errno.ErrSelectTag, nil, err.Error())
			return
		}

		var tagValue blog.Tag
		if tagCount > 0 {
			// 绑定以后的标签数据
			err = tx.Model(&blog.Tag{}).Where("name = ?", tag).Find(&tagValue).Error
			if err != nil {
				tx.Rollback()
				Response(c, errno.ErrSelectTag, nil, err.Error())
				return
			}
		} else {
			// 新建标签数据
			tagValue.Name = tag
			err = tx.Create(&tagValue).Error
			if err != nil {
				tx.Rollback()
				Response(c, errno.ErrCreateTag, nil, err.Error())
				return
			}
		}

		// 创建标签与文章之前的关联
		contentTag := blog.ContentTag{
			Content: contentValue.Id,
			Tag:     tagValue.Id,
		}

		err := tx.Create(&contentTag).Error
		if err != nil {
			tx.Rollback()
			Response(c, errno.ErrCreateBlogTag, nil, err.Error())
			return
		}
	}

	tx.Commit()

	Response(c, nil, nil, "")
}
