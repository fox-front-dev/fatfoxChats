package service

import "github.com/gin-gonic/gin"

// @BasePath /api/v1

// GetIndex
// @Tags 获取index
// @Success 200 {string} ok
// @Router /index [get]
func GetIndex(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "成功",
	})
}
