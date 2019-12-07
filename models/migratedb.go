/*
  @Author : lanyulei
*/

package models

import (
	"blog/models/blog"
	"blog/models/user"
	"blog/pkg/connection"
)

func AutoMigrateTable() {
	connection.DB.Self.AutoMigrate(
		// blog
		&blog.Content{},
		&blog.Tag{},
		&blog.Type{},
		&blog.ContentTag{},
		&blog.Links{},
		&blog.Series{},
		&blog.SeriesDetails{},
		&blog.Comment{},
		&blog.Reply{},

		// user
		&user.User{},
	)
}
