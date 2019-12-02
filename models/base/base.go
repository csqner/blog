/*
  @Author : lanyulei
*/

package base

import (
	"blog/utils"
)

type Model struct {
	Id        int             `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id" form:"id"`
	CreatedAt utils.JSONTime  `gorm:"column:created_at" json:"created_at" form:"created_at"`
	UpdatedAt utils.JSONTime  `gorm:"column:updated_at" json:"updated_at" form:"updated_at"`
	DeletedAt *utils.JSONTime `gorm:"column:deleted_at" sql:"index" json:"-"`
}

type ListRequest struct {
	Offset int    `form:"offset"`
	Limit  int    `form:"limit"`
	Item   string `form:"item"`
}
