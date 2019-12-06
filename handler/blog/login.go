package blog

import (
	"blog/pkg/logger"
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

// 1. Get Authorization Code
func GetAuthCode(c *gin.Context) {
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", viper.GetString("qq.appId"))
	params.Add("state", "lanyulei")
	str := fmt.Sprintf("%s&redirect_uri=%s", params.Encode(), viper.GetString("qq.redirectURI"))
	loginURL := fmt.Sprintf("%s?%s", "https://graph.qq.com/oauth2.0/authorize", str)
	c.Redirect(http.StatusMovedPermanently, loginURL)
}

// 2. Get Access Token
func GetToken(c *gin.Context) {
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

	resultMap := convertToMap(body)

	info := &PrivateInfo{}
	info.AccessToken = resultMap["access_token"]
	info.RefreshToken = resultMap["refresh_token"]
	info.ExpiresIn = resultMap["expires_in"]

	GetOpenId(info, c)
}

// 3. Get OpenId
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

// 4. Get User info
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

	c.Redirect(http.StatusMovedPermanently, "/")
}

func convertToMap(str string) map[string]string {
	var resultMap = make(map[string]string)
	values := strings.Split(str, "&")
	for _, value := range values {
		vs := strings.Split(value, "=")
		resultMap[vs[0]] = vs[1]
	}
	return resultMap
}
