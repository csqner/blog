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
	"time"
)

/*
  @Author : lanyulei
*/

// 首页
func IndexHandler(c *gin.Context) {
	// 获取访问量最高的文章
	var contentValueList []blog.Content
	err := connection.DB.Self.Model(&blog.Content{}).Limit(6).Order("view desc").Find(&contentValueList).Error
	if err != nil {
		Response(c, errno.ErrSelectList, nil, err.Error())
		return
	}

	c.HTML(http.StatusOK, "Index.html", gin.H{
		"content": contentValueList,
	})
}

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
		Order("created_at desc").
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

	// 获取访问量最高的文章
	var contentValueList []blog.Content
	err = connection.DB.Self.Model(&blog.Content{}).Limit(10).Order("view desc").Find(&contentValueList).Error
	if err != nil {
		Response(c, errno.ErrSelectList, nil, err.Error())
		return
	}

	c.HTML(http.StatusOK, "Article.html", gin.H{
		"content":     contentList,
		"types":       typeList,
		"count":       contentCount,
		"contentList": contentValueList,
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

// 文章详情
func ContentDetailsHandler(c *gin.Context) {
	contentId := c.DefaultQuery("content_id", "")
	if contentId == "" {
		Response(c, errno.InternalServerError, nil, "content_id的参数不存在，请确认")
		return
	}

	var contentTmp blog.Content
	err := connection.DB.Self.Model(&contentTmp).Where("id = ?", contentId).Find(&contentTmp).Error
	if err != nil {
		Response(c, errno.ErrSelectDetails, nil, err.Error())
		return
	}

	// 更新访问量
	err = connection.DB.Self.Model(&blog.Content{}).Where("id = ?", contentId).Update("view", contentTmp.View+1).Error
	if err != nil {
		Response(c, errno.ErrUpdateView, nil, err.Error())
		return
	}

	type contentModel struct {
		blog.Content
		TypeName string `json:"type_name"`
	}

	var content contentModel
	err = connection.DB.Self.Model(&blog.Content{}).
		Joins("left join blog_type on blog_content.type_id = blog_type.id").
		Select("blog_content.*, blog_type.name as type_name").
		Where("blog_content.id = ?", contentId).
		Find(&content).Error
	if err != nil {
		Response(c, errno.ErrSelectDetails, nil, err.Error())
		return
	}

	c.HTML(http.StatusOK, "ArticleDetails.html", gin.H{
		"content": content,
	})
}

// 文章归档
func ArchiveHandler(c *gin.Context) {
	startData := c.DefaultQuery("date", "")
	var monthCount int
	if startData != "" {
		dateValue := strings.Split(startData, "-")
		year, _ := strconv.Atoi(dateValue[0])
		month, _ := strconv.Atoi(dateValue[1])
		yearTmp := 0
		for i := year; i <= time.Now().Year(); i++ {
			if i < time.Now().Year() {
				yearTmp = i
				for m := month; m <= 12; m++ {
					monthCount += 1
				}
			} else if yearTmp != 0 && i == time.Now().Year() {
				for m := 1; m <= int(time.Now().Month()); m++ {
					monthCount += 1
				}
			} else if yearTmp == 0 && i == time.Now().Year() {
				for m := month; m <= int(time.Now().Month()); m++ {
					monthCount += 1
				}
			}
		}
	} else {
		monthCount = 6
	}

	// 归档数据
	var contentList []blog.ContentListStruct
	err := connection.DB.Self.Table("blog_content").
		Select("id, title, date_format(created_at, '%Y-%m') as created_at").
		Where("date_sub(current_date(), interval ? month) <= date(created_at)", monthCount).
		Order("created_at desc").
		Scan(&contentList).Error
	if err != nil {
		Response(c, errno.ErrSelectList, nil, err.Error())
		return
	}

	var ArchiveDataList blog.ArchiveData
	for _, content := range contentList {
		if content.CreatedAt != "" {
			IsExist := false
			// 判断是否存在，存在则追加
			for index, Article := range ArchiveDataList.ArticleGroup {
				if Article.ArchiveDataString == content.CreatedAt {
					IsExist = true
					ArchiveDataList.ArticleGroup[index].ArticleData = append(ArchiveDataList.ArticleGroup[index].ArticleData, content)
					break
				}
			}

			// 添加新数据
			if IsExist == false {
				ArchiveDataList.ArticleGroup = append(ArchiveDataList.ArticleGroup, blog.ArticleDataStruct{
					ArchiveDataString: content.CreatedAt,
					ArticleData:       []blog.ContentListStruct{content},
				})
			}
		}
	}

	// 获取分类总数
	err = connection.DB.Self.Model(&blog.Type{}).Count(&ArchiveDataList.TypeCount).Error
	if err != nil {
		Response(c, errno.ErrTypeCount, nil, err.Error())
		return
	}

	// 获取文章总数
	err = connection.DB.Self.Model(&blog.Content{}).Count(&ArchiveDataList.ArticleCount).Error
	if err != nil {
		Response(c, errno.ErrTotalCount, nil, err.Error())
		return
	}

	// 获取评价总数

	// 获取赞的总数
	var awesomeCount struct {
		Awesome int `json:"awesome"`
	}
	err = connection.DB.Self.Table("blog_content").Select("sum(`awesome`) as awesome").Scan(&awesomeCount).Error
	if err != nil {
		Response(c, errno.ErrAwesomeCount, nil, err.Error())
		return
	}
	ArchiveDataList.AwesomeCount = awesomeCount.Awesome

	// 获取阅读总数
	var articleCount struct {
		SumView int `json:"sum_view"`
	}
	err = connection.DB.Self.Table("blog_content").Select("sum(`view`) as sum_view").Scan(&articleCount).Error
	if err != nil {
		Response(c, errno.ErrTotalView, nil, err.Error())
		return
	}
	ArchiveDataList.ViewCount = articleCount.SumView

	// 获取最近的更新时间
	var articleValue blog.Content
	err = connection.DB.Self.Model(&articleValue).Limit(1).Order("updated_at desc").Find(&articleValue).Error
	if err != nil {
		Response(c, errno.ErrLastTime, nil, err.Error())
		return
	}
	ArchiveDataList.LastUpdateTime = articleValue.UpdatedAt

	// 获取第一篇文章的创建时间
	var articleValueFirst blog.Content
	err = connection.DB.Self.Model(&articleValueFirst).Limit(1).Order("created_at").Find(&articleValueFirst).Error
	if err != nil {
		Response(c, errno.ErrSelectDetails, nil, err.Error())
		return
	}

	var dateValue []blog.DateStruct
	frontYear := 0
	for i := articleValueFirst.CreatedAt.Year(); i <= time.Now().Year(); i++ {
		var monthList []blog.MonthStruct
		if i < time.Now().Year() {
			frontYear = i
			for m := int(articleValueFirst.CreatedAt.Month()); m <= 12; m++ {
				monthList = append(monthList, blog.MonthStruct{
					Month: m,
				})
			}
		} else if frontYear == 0 && i == time.Now().Year() {
			for m := int(articleValueFirst.CreatedAt.Month()); m <= int(time.Now().Month()); m++ {
				monthList = append(monthList, blog.MonthStruct{
					Month: m,
				})
			}
		} else if frontYear != 0 && i == time.Now().Year() {
			for m := 1; m <= int(time.Now().Month()); m++ {
				monthList = append(monthList, blog.MonthStruct{
					Month: m,
				})
			}
		}
		if len(monthList) > 0 {
			dateValue = append(dateValue, blog.DateStruct{
				Year:   i,
				Months: monthList,
			})
		}
	}

	c.HTML(http.StatusOK, "Archive.html", gin.H{
		"content": ArchiveDataList,
		"date":    dateValue,
	})
}

// 关于
func AboutHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "About.html", gin.H{})
}

// 博主
func AuthorHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "Author.html", gin.H{})
}

// 友链
func LinksHandler(c *gin.Context) {
	var linksList []blog.Links
	err := connection.DB.Self.Model(&blog.Links{}).Find(&linksList).Error
	if err != nil {
		Response(c, errno.ErrSelectLinks, nil, err.Error())
		return
	}

	c.HTML(http.StatusOK, "FriendlyLink.html", gin.H{
		"links": linksList,
	})
}
