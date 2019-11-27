/*
  @Author : lanyulei
*/

package connection

import (
	"blog/pkg/logger"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

type Database struct {
	Self *gorm.DB
}

var DB *Database

func openDB(username, password, addr, name string) *gorm.DB {
	config := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s",
		username,
		password,
		addr,
		name,
		true,
		"Local")

	db, err := gorm.Open("mysql", config)
	if err != nil {
		logger.Errorf("数据库连接失败，连接地址: %s，error: %s", viper.GetString(`db.addr`), err)
	}

	// set for db connection
	setupDB(db)

	return db
}

func setupDB(db *gorm.DB) {
	// 是否开启详细日志记录
	db.LogMode(viper.GetBool(`db.gorm.logMode`))

	// 设置最大打开连接数
	db.DB().SetMaxOpenConns(viper.GetInt(`db.gorm.maxOpenConn`))

	// 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用
	db.DB().SetMaxIdleConns(viper.GetInt(`db.gorm.maxIdleConn`))

	// 创建表的时候去掉复数
	db.SingularTable(viper.GetBool(`db.gorm.singularTable`))
}

func InitSelfDB() *gorm.DB {
	return openDB(viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.addr"),
		viper.GetString("db.name"))
}

func (db *Database) Init() {
	DB = &Database{
		Self: InitSelfDB(),
	}
}

func (db *Database) Close() {
	err := DB.Self.Close()
	if err != nil {
		logger.Errorf("关闭连接失败，错误信息: %s", err)
	}
}
