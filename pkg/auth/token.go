/*
  @Author : lanyulei
*/

package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"time"
)

var (
	ErrMissingHeader = errors.New("Header中没有找到Authorization")
)

type Context struct {
	ID       uint64
	Username string
}

// secretFunc 验证密钥格式
func secretFunc(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	}
}

// 解密token获取对应的用户信息
func Parse(tokenString string, secret string) (*Context, error) {
	ctx := &Context{}

	// 验证token
	token, err := jwt.Parse(tokenString, secretFunc(secret))

	// 验证失败
	if err != nil {
		return ctx, err

		// 获取token对应的数据，并且验证token是否有效
	} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ctx.ID = uint64(claims["id"].(float64))
		ctx.Username = claims["username"].(string)
		return ctx, nil

	} else {
		return ctx, err
	}
}

// 从Header获取token信息
func ParseRequest(c *gin.Context) (*Context, error) {
	// 获取header中的token
	authorization := c.Request.Header.Get("Authorization")

	// 获取jwt secret
	secret := viper.GetString(`jwt.secret`)

	if len(authorization) == 0 {
		return &Context{}, ErrMissingHeader
	}

	return Parse(authorization, secret)
}

func Sign(ctx *gin.Context, c Context, secret string) (tokenString string, err error) {

	if secret == "" {
		secret = viper.GetString(`jwt.secret`)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       c.ID,
		"username": c.Username,
		"nbf":      time.Now().Unix(),
		"iat":      time.Now().Unix(),
	})

	tokenString, err = token.SignedString([]byte(secret))

	return
}
