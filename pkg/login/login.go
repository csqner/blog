package login

import (
	"blog/models/user"
	"blog/pkg/connection"
	"blog/pkg/logger"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

/*
  @Author : lanyulei
*/

type PrivateInfo struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
}

func GetOpenId(info *PrivateInfo, c *gin.Context) {
	resp, err := http.Get(fmt.Sprintf("%s?access_token=%s", "https://graph.qq.com/oauth2.0/me", info.AccessToken))
	if err != nil {
		logger.Error(err)
		return
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			logger.Error(err)
			return
		}
	}()

	bs, _ := ioutil.ReadAll(resp.Body)
	body := string(bs)
	info.OpenId = body[45:77]

	GetUserInfo(info, c)
}

func GetUserInfo(info *PrivateInfo, c *gin.Context) {
	params := url.Values{}
	params.Add("access_token", info.AccessToken)
	params.Add("openid", info.OpenId)
	params.Add("oauth_consumer_key", viper.GetString("qq.appId"))

	uri := fmt.Sprintf("https://graph.qq.com/user/get_user_info?%s", params.Encode())
	resp, err := http.Get(uri)
	if err != nil {
		logger.Error(err)
		return
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			logger.Error(err)
			return
		}
	}()

	userInfo, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return
	}

	userMap := make(map[string]interface{})
	err = json.Unmarshal(userInfo, &userMap)
	if err != nil {
		logger.Error(err)
		return
	}

	// 1. 查询OpenId是否存在
	var userCount int
	err = connection.DB.Self.Model(&user.User{}).Where("open_id = ?", info.OpenId).Count(&userCount).Error
	if err != nil {
		logger.Errorf("查询OpenId失败，错误：%v", err)
		return
	}

	userData := user.User{
		OpenId:   info.OpenId,
		Nickname: userMap["nickname"].(string),
		Avatar:   userMap["figureurl_qq"].(string),
	}

	if userCount > 0 {
		// 更新数据
		err = connection.DB.Self.Model(&user.User{}).Where("open_id = ?", info.OpenId).Updates(userData).Error
		if err != nil {
			logger.Errorf("更新用户数据失败，错误：%v", err)
			return
		}
	} else {
		// 写入数据
		err = connection.DB.Self.Model(&user.User{}).Create(&userData).Error
		if err != nil {
			logger.Errorf("创建用户数据失败，错误：%v", err)
			return
		}
	}

	// 写入Session
	SaveAuthSession(c, info.OpenId)

	c.Redirect(http.StatusMovedPermanently, "/")
}

func ConvertToMap(str string) map[string]string {
	var resultMap = make(map[string]string)
	values := strings.Split(str, "&")
	for _, value := range values {
		vs := strings.Split(value, "=")
		resultMap[vs[0]] = vs[1]
	}
	return resultMap
}
