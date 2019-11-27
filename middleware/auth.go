/*
  @Author : lanyulei
*/

package middleware

import (
	"blog/pkg/auth"
	"blog/pkg/errno"
	. "blog/pkg/response"
	"github.com/gin-gonic/gin"
)

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
