package blog

import (
	"blog/models/blog"
	"blog/pkg/connection"
	"blog/pkg/errno"
	. "blog/pkg/response"
	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

// 获取文章分类列表
func TypeListHandler(c *gin.Context) {
	var (
		typeList []*blog.Type
		err      error
	)
	err = connection.DB.Self.Model(&blog.Type{}).Find(&typeList).Error
	if err != nil {
		Response(c, errno.ErrSelectTypeList, nil, err.Error())
		return
	}

	Response(c, nil, typeList, "")
}

func CreateTypeHandler(c *gin.Context) {
	// 绑定参数
	var typeParams blog.Type
	err := c.ShouldBind(&typeParams)
	if err != nil {
		Response(c, errno.ErrBind, nil, err.Error())
		return
	}

	// 判断分类是否存在
	var typeCount int
	err = connection.DB.Self.Model(&blog.Type{}).Where("name = ?", typeParams.Name).Count(&typeCount).Error
	if err != nil {
		Response(c, errno.ErrSelectTypeList, nil, err.Error())
		return
	}

	if typeCount > 0 {
		Response(c, errno.ErrExists, nil, "分类名称已存在，请重新输入")
		return
	}

	// 创建文章分类
	err = connection.DB.Self.Create(&typeParams).Error
	if err != nil {
		Response(c, errno.ErrCreateType, nil, err.Error())
		return
	}

	Response(c, nil, typeParams, "")
}

func UpdateTypeHandler(c *gin.Context) {
	// 绑定参数
	var typeParams blog.Type
	err := c.ShouldBind(&typeParams)
	if err != nil {
		Response(c, errno.ErrBind, nil, err.Error())
		return
	}

	// 判断分类是否存在
	var typeCount int
	err = connection.DB.Self.Model(&blog.Type{}).Where("name = ?", typeParams.Name).Count(&typeCount).Error
	if err != nil {
		Response(c, errno.ErrSelectTypeList, nil, err.Error())
		return
	}

	if typeCount > 0 {
		Response(c, errno.ErrExists, nil, "分类名称已存在，请重新输入")
		return
	}

	// 创建文章分类
	err = connection.DB.Self.Model(&blog.Type{}).Where("id = ?", typeParams.Id).Update("name", typeParams.Name).Error
	if err != nil {
		Response(c, errno.ErrUpdateType, nil, err.Error())
		return
	}

	Response(c, nil, typeParams, "")
}
