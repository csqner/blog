package blog

import (
	"blog/models/base"
	"blog/models/blog"
	"blog/models/user"
	"blog/pkg/connection"
	"blog/pkg/errno"
	. "blog/pkg/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
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

	// 最近留言
	type UserCommentStruct struct {
		blog.Comment
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	}

	var commentList []UserCommentStruct
	err = connection.DB.Self.Table("blog_comment").
		Joins("left join blog_user on blog_comment.user_id = blog_user.id").
		Select("blog_comment.*, blog_user.nickname, blog_user.avatar").
		Limit(6).
		Order("id desc").
		Scan(&commentList).Error
	if err != nil {
		HtmlResponse(c, "error.html", err.Error(), "/blog/index")
		return
	}

	// 统计评论人员
	var commentUser []struct {
		Nickname     string `json:"nickname"`
		Avatar       string `json:"avatar"`
		CommentCount int    `json:"comment_count"`
		Order        int    `json:"order"`
	}
	err = connection.DB.Self.Table("blog_comment").
		Joins("left join blog_user on blog_user.id = blog_comment.user_id").
		Select("blog_user.nickname, blog_user.avatar, count(*) as comment_count").
		Order("comment_count desc").
		Limit(3).
		Group("user_id").Scan(&commentUser).Error
	if err != nil {
		HtmlResponse(c, "error.html", err.Error(), "/blog/index")
		return
	}
	orderValue := 1
	for index, _ := range commentUser {
		commentUser[index].Order = orderValue
		orderValue += 1
	}

	// 近期访客
	var userInfoList []user.User
	err = connection.DB.Self.Model(&user.User{}).Order("created_at desc").Limit(20).Find(&userInfoList).Error
	if err != nil {
		HtmlResponse(c, "error.html", err.Error(), "/blog/index")
		return
	}

	// 获取哲理句子
	philosophyId := rand.Intn(145)
	if philosophyId == 0 {
		philosophyId += 1
	}
	var philosophyValue blog.Philosophy
	err = connection.DB.Self.Model(&philosophyValue).Where("id = ?", int(philosophyId)).Find(&philosophyValue).Error
	if err != nil {
		HtmlResponse(c, "error.html", err.Error(), "/blog/index")
		return
	}

	c.HTML(http.StatusOK, "Index.html", gin.H{
		"contentList":     contentValueList,
		"commentList":     commentList,
		"commentUser":     commentUser,
		"userInfoList":    userInfoList,
		"philosophyValue": philosophyValue,
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
		TypeName     string `json:"type_name"`
		CommentCount int    `json:"comment_count"`
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

	// 获取用户评论数
	var aricleCommentIdList []int
	err = connection.DB.Self.Model(&blog.Comment{}).Where("article_id = ? and type = 1", contentId).Pluck("id", &aricleCommentIdList).Error
	if err != nil {
		Response(c, errno.ErrSelectComment, nil, err.Error())
		return
	}
	content.CommentCount = len(aricleCommentIdList)

	var articleReplyCount int
	err = connection.DB.Self.Model(&blog.Reply{}).Where("comment_id in (?)", aricleCommentIdList).Count(&articleReplyCount).Error
	if err != nil {
		Response(c, errno.ErrSelectComment, nil, err.Error())
		return
	}

	content.CommentCount += articleReplyCount

	// 获取所有评论
	var commentList []blog.CommentStruct
	err = connection.DB.Self.Table("blog_comment").
		Joins("left join blog_user on blog_comment.user_id = blog_user.id").
		Select("blog_comment.*, blog_user.nickname, blog_user.avatar, blog_user.title as user_title").
		Where("blog_comment.deleted_at is null and blog_comment.article_id = ? and type = 1", contentId).
		Order("created_at desc").
		Scan(&commentList).Error
	if err != nil {
		HtmlResponse(c, "error.html", fmt.Sprintf("获取评论错误，%v", err.Error()), "/blog/list")
		return
	}

	var commentIdList []int
	err = connection.DB.Self.Model(&blog.Comment{}).Where("article_id = ? and type = 1", contentId).
		Pluck("id", &commentIdList).Error
	if err != nil {
		HtmlResponse(c, "error.html", fmt.Sprintf("获取评论错误，%v", err.Error()), "/blog/list")
		return
	}

	// 获取当前文章的所有回复
	var replyList []blog.ReplyStruct
	err = connection.DB.Self.Table("blog_reply").
		Joins("left join blog_user as s_blog_user on s_blog_user.id = blog_reply.source_user_id").
		Joins("left join blog_user as d_blog_user on d_blog_user.id = blog_reply.aims_user_id").
		Select("blog_reply.*, s_blog_user.id as s_user_id, s_blog_user.nickname as s_nickname, s_blog_user.title as s_user_title, s_blog_user.avatar as s_avatar, d_blog_user.id as d_user_id, d_blog_user.nickname as d_nickname, d_blog_user.title as d_user_title, d_blog_user.avatar as d_avatar").
		Order("created_at").
		Where("blog_reply.deleted_at is null and blog_reply.comment_id in (?)", commentIdList).
		Scan(&replyList).Error
	if err != nil {
		HtmlResponse(c, "error.html", fmt.Sprintf("获取评论错误，%v", err.Error()), "/blog/list")
		return
	}

	for _, reply := range replyList {
		for commentIndex, comment := range commentList {
			if reply.CommentId == comment.Id {
				commentList[commentIndex].ReplyList = append(commentList[commentIndex].ReplyList, reply)
			}
		}
	}

	// 获取上一篇文章
	var onContentCount int
	err = connection.DB.Self.Model(&blog.Content{}).Where("id > ?", contentId).Limit(1).Count(&onContentCount).Error
	if err != nil {
		HtmlResponse(c, "error.html", fmt.Sprintf("获取上文失败，%v", err.Error()), "/blog/list")
		return
	}

	var onContent blog.Content
	if onContentCount > 0 {
		err := connection.DB.Self.Model(&blog.Content{}).Where("id > ?", contentId).Limit(1).Find(&onContent).Error
		if err != nil {
			HtmlResponse(c, "error.html", fmt.Sprintf("获取上文失败，%v", err.Error()), "/blog/list")
			return
		}
	}

	// 获取下一篇文章
	var underContentCount int
	err = connection.DB.Self.Model(&blog.Content{}).Where("id < ?", contentId).Limit(1).Count(&underContentCount).Error
	if err != nil {
		HtmlResponse(c, "error.html", fmt.Sprintf("获取下文失败，%v", err.Error()), "/blog/list")
		return
	}

	var underContent blog.Content
	if underContentCount > 0 {
		err := connection.DB.Self.Model(&blog.Content{}).Where("id < ?", contentId).Limit(1).Find(&underContent).Error
		if err != nil {
			HtmlResponse(c, "error.html", fmt.Sprintf("获取上文失败，%v", err.Error()), "/blog/list")
			return
		}
	}

	c.HTML(http.StatusOK, "ArticleDetails.html", gin.H{
		"content":      content,
		"comment":      commentList,
		"onContent":    onContent,
		"underContent": underContent,
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
	var commentCount int
	err = connection.DB.Self.Model(&blog.Comment{}).Count(&commentCount).Error
	if err != nil {
		Response(c, errno.ErrSelectComment, nil, err.Error())
		return
	}
	var replyCount int
	err = connection.DB.Self.Model(&blog.Reply{}).Count(&replyCount).Error
	if err != nil {
		Response(c, errno.ErrSelectComment, nil, err.Error())
		return
	}
	commentCount += replyCount
	ArchiveDataList.CommentCount = commentCount

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

func AddAwesomeHandler(c *gin.Context) {
	contentId := c.DefaultQuery("content_id", "")
	if contentId == "" {
		Response(c, errno.ErrBind, nil, "参数不正确")
		return
	}

	// 查询原来的赞
	var contentValue blog.Content
	err := connection.DB.Self.Model(&contentValue).Where("id = ?", contentId).Find(&contentValue).Error
	if err != nil {
		Response(c, errno.ErrSelectDetails, nil, err.Error())
		return
	}

	// 更新赞的数量
	err = connection.DB.Self.Model(&blog.Content{}).Where("id = ?", contentId).
		Update("awesome", contentValue.Awesome+1).Error
	if err != nil {
		Response(c, errno.ErrUpdateAwesome, nil, err.Error())
		return
	}

	Response(c, nil, nil, "")
}

func SiteMapHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "siteMap.html", gin.H{})
}
