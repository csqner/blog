package login

/*
  @Author : lanyulei
*/

import (
	"blog/pkg/logger"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// gin session key
const KEY = "647983B6B6A4ED0DEF81666268BBDFB6"

// 使用 Cookie 保存 session
func EnableCookieSession() gin.HandlerFunc {
	store := cookie.NewStore([]byte(KEY))
	return sessions.Sessions("BLOG", store)
}

// 登陆时都需要保存seesion信息
func SaveAuthSession(c *gin.Context, openId string) {
	session := sessions.Default(c)
	session.Set("openId", openId)
	err := session.Save()
	if err != nil {
		logger.Errorf("Session保存失败，错误：%v", err)
	}
}

// 退出时清除session
func ClearAuthSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		logger.Errorf("清除Session失败，错误：%v", err)
	}
}

// 判断Session是否存在
func HasSession(c *gin.Context) bool {
	session := sessions.Default(c)
	if sessionValue := session.Get("openId"); sessionValue == nil {
		return false
	}
	return true
}

// 获取Session
func GetSessionOpenId(c *gin.Context) string {
	session := sessions.Default(c)
	sessionValue := session.Get("openId")
	if sessionValue == nil {
		return ""
	}
	return sessionValue.(string)
}
