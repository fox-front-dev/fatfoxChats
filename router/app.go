package router

import (
	"fatfoxChats/service"
	"fatfoxChats/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.Use(Cors())
	// 创建用户
	r.POST("/user/createUser", service.CreateUserService)
	// 注销用户
	r.POST("/user/delUser", service.DelectUser)
	// 登录
	r.POST("/user/loginUser", service.LoginUser)
	//更新用户
	r.POST("/user/modifyInfo", service.UpdataUser)
	// 发送消息ws
	r.GET("/user/sendUserMsg", service.SendUserMsg)
	// -------------好友操作--------------
	// 查询用户的好友
	r.GET("/user/searchFriends", service.SearchFriends)
	// 查询用户支持多个id查询
	r.GET("/user/selectUserList", service.SearchUserLists)
	// 查询用户信息,手机号码查询（添加好友使用）
	r.POST("/user/selectUserByPhone", service.SelectUserByPhone)
	// 添加好友
	r.POST("/user/addFriend", service.AddFriends_service)
	// 删除好友
	r.POST("/user/delFriend", service.DelFriends)
	// 同意添加好友
	r.POST("/user/agreeAdd", service.AgreeAdd)
	// 查询用户好友申请
	r.POST("/user/relationship", service.Relationship)
	// -------------空间说说操作--------------
	// 获取七牛token
	r.GET("/file/qiniutoken", utils.GETQINIUTOKEN)
	// 增说说
	r.POST("/space/addSpaceRecord", ValidToken, service.AddSpaceRecord)
	// 修改说说状态
	r.POST("/space/updateSpaceRecord", ValidToken, service.UpdateSpaceRecord)
	// 空间说说
	r.POST("/space/getAllSpaceContent", ValidToken, service.GetAllSpaceContent)
	return r
}
func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
	}
}
func ValidToken(c *gin.Context) {
	sToken := c.GetHeader("AUTHORIZATION")
	cla, err := utils.ParseToken(sToken)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1001,
			"msg":  "账号异常",
		})
		c.Abort()
		return
	}
	c.Set("Id", cla.Id)
	c.Next()
}
