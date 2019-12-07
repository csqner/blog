package middleware

/*
  @Author : lanyulei
*/

import (
	"blog/pkg/auth"
	"blog/pkg/errno"
	. "blog/pkg/response"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// session中间件
func AuthSessionMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionValue := session.Get("openId")
		if sessionValue == nil {
			Response(c, errno.ErrNotLogin, nil, "")
			c.Abort()
			return
		}
		c.Set("openId", sessionValue)
		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析验证token
		if _, err := auth.ParseRequest(c); err != nil {
			Response(c, errno.ErrTokenInvalid, nil, "")
			c.Abort()
			return
		}
		c.Next()
	}
}
