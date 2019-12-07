package user

import "blog/models/base"

/*
  @Author : lanyulei
*/

type User struct {
	base.Model
	OpenId   string `gorm:"column:open_id;type:varchar(256)" json:"open_id" form:"open_id"`
	Nickname string `gorm:"column:nickname;type:varchar(256)" json:"nickname" form:"nickname"`
	Avatar   string `gorm:"column:avatar;type:varchar(256)" json:"avatar" form:"avatar"`
	Title    string `gorm:"column:title;type:varchar(128)" json:"title" form:"title"`
}

func (User) TableName() string {
	return "blog_user"
}
