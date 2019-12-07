package blog

import (
	"blog/models/user"
	"blog/pkg/connection"
	"blog/pkg/errno"
	"blog/pkg/logger"
	"blog/pkg/login"
	. "blog/pkg/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"net/url"
)

/*
  @Author : lanyulei
*/

func ToLoginHandler(c *gin.Context) {
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", viper.GetString("qq.appId"))
	params.Add("state", "lanyulei")
	str := fmt.Sprintf("%s&redirect_uri=%s", params.Encode(), viper.GetString("qq.redirectURI"))
	loginURL := fmt.Sprintf("%s?%s", "https://graph.qq.com/oauth2.0/authorize", str)
	c.Redirect(http.StatusMovedPermanently, loginURL)
}

func CallbackHandler(c *gin.Context) {
	code := c.DefaultQuery("code", "")
	params := url.Values{}
	params.Add("grant_type", "authorization_code")
	params.Add("client_id", viper.GetString("qq.appId"))
	params.Add("client_secret", viper.GetString("qq.appKey"))
	params.Add("code", code)
	str := fmt.Sprintf("%s&redirect_uri=%s", params.Encode(), viper.GetString("qq.redirectURI"))
	loginURL := fmt.Sprintf("%s?%s", "https://graph.qq.com/oauth2.0/token", str)

	response, err := http.Get(loginURL)
	if err != nil {
		logger.Error(err)
		return
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			logger.Error(err)
			return
		}
	}()

	bs, _ := ioutil.ReadAll(response.Body)
	body := string(bs)

	resultMap := login.ConvertToMap(body)

	info := &login.PrivateInfo{}
	info.AccessToken = resultMap["access_token"]
	info.RefreshToken = resultMap["refresh_token"]
	info.ExpiresIn = resultMap["expires_in"]

	login.GetOpenId(info, c)
}

// 获取用户数据
func UserInfoHandler(c *gin.Context) {
	var userInfo user.User
	if login.HasSession(c) {
		openId := login.GetSessionOpenId(c)
		err := connection.DB.Self.Model(&user.User{}).Where("open_id = ?", openId).Find(&userInfo).Error
		if err != nil {
			Response(c, errno.ErrSelectUser, nil, err.Error())
			return
		}
	} else {
		Response(c, errno.ErrNotLogin, nil, "")
		return
	}
	Response(c, nil, userInfo, "")
}

// 判断账号是否登陆
func IsLoginHandler(c *gin.Context) {
	loginStatus := login.HasSession(c)
	Response(c, nil, loginStatus, "")
}
