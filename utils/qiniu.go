package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

func GETQINIUTOKEN(c *gin.Context) {
	token := c.GetHeader("AUTHORIZATION")
	cla, err := ParseToken(token)

	if err != nil {
		c.JSON(200, gin.H{
			"code": 1001,
			"msg":  "账号异常",
		})
		return
	}
	if cla.Id <= 0 {
		c.JSON(200, gin.H{
			"code": 1001,
			"msg":  "账号异常",
		})
		return
	}
	accessKey := "BeqLsFg2f7UOgriGMEIyVsAmtHLRki6xnrKqa0xf"
	secretKey := "SFJkgY46KTk0R64dhHHDlJaj0OtzeR6p2AzKnJZE"
	bucket := "fatfox-chat-space"
	putPolicy := storage.PutPolicy{
		Scope:   bucket,
		Expires: 300,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	c.JSON(200, gin.H{
		"code":  1000,
		"msg":   "",
		"qiniu": upToken,
	})

}
