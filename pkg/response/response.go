package response

/*
  @Author : lanyulei
*/

import (
	"blog/pkg/errno"
	"blog/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseCode struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func HtmlResponse(c *gin.Context, requestUrl string, errorContent string, jumpUrl string) {
	c.HTML(http.StatusOK, requestUrl, gin.H{
		"ErrorContent": errorContent,
		"JumpUrl":      jumpUrl,
	})
}

func Response(c *gin.Context, err error, data interface{}, errText string) {
	code, message := errno.DecodeErr(err)

	if errText != "" {
		message = errText
	}

	// write log
	if code != 10000 {
		logger.Error(message)
	} else if message != "OK" {
		logger.Info(message)
	}

	// always return http.StatusOK
	c.JSON(http.StatusOK, ResponseCode{
		Code:    code,
		Message: message,
		Data:    data,
	})
}
