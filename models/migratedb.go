/*
  @Author : lanyulei
*/

package models

import (
	"blog/models/blog"
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
	)
}
