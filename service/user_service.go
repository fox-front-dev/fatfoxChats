package service

import (
	"fatfoxChats/models"
	"fatfoxChats/utils"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func SelectUser(c *gin.Context) {
	Id := c.Param("Id")
	user := models.FindUserById(Id)
	id := strconv.FormatUint(user.Id, 10)

	if strings.EqualFold(id, Id) {
		c.JSON(200, gin.H{
			"code": 1000,
			"data": user,
			"msg":  "",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 1001,
		"msg":  "请稍后再试",
	})
}

// @Summary ping GetUserList
// @Tags 获取用户列表
// @Success 200 {string} data
// @Router /user/getUserList [get]

func CreateUserService(c *gin.Context) {
	var user models.UserInfo
	e := c.ShouldBind(&user)
	if e != nil {
		c.JSON(200, gin.H{
			"code":   1001,
			"msg":    "请填写完整参数!",
			"params": user,
		})
		return
	}
	users := models.FindUserByPhone(user.Phone)
	if !strings.EqualFold(users.Phone, user.Phone) {
		rand.Seed(time.Now().Unix())
		ranNum := strconv.Itoa(rand.Int())
		user.Salt = ranNum
		user.Password = utils.Md5Encode(user.Password + ranNum)
		models.CreateUser(user)
		c.JSON(200, gin.H{
			"code": 1000,
			"msg":  "创建用户成功",
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "该账号已存在!",
	})
	return
}

func DelectUser(c *gin.Context) {
	type Users struct {
		Id    uint64 `binding:"required"`
		Pwd   string `binding:"required"`
		Phone string `binding:"required"`
	}
	var u Users
	err := c.ShouldBind(&u)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1001,
			"msg":  "请填写完整!",
		})
		return
	}
	user := models.FindUserByPhone(u.Phone)
	i1 := user.Id
	i2 := u.Id
	if i1 == i2 {
		b := utils.ValidPassWord(u.Pwd, user.Salt, user.Password)
		if b {
			models.DelectUser(u.Id, user.Password, u.Phone)
			c.JSON(200, gin.H{
				"code": 200,
				"msg":  "注销成功！！！",
			})
			return
		}
		c.JSON(200, gin.H{
			"code": 1001,
			"msg":  "注销失败,请稍后再试",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "暂无该用户！！！",
	})
}

// 登录
func LoginUser(c *gin.Context) {
	type userLoginParams struct {
		Phone    string `binding:"required"`
		Password string `binding:"required"`
	}
	var Params userLoginParams
	err := c.ShouldBind(&Params)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1001,
			"msg":  "请填写完整",
		})
		return
	}
	users := models.FindUserByPhone(Params.Phone)
	b := utils.ValidPassWord(Params.Password, users.Salt, users.Password)
	if !b {
		c.JSON(200, gin.H{
			"code": 1002,
			"msg":  "账号或者密码错误",
		})
		return
	}
	token, tokenErr := utils.GetToken(users.Id)
	if tokenErr != nil {
		c.JSON(200, gin.H{
			"code": 1003,
			"msg":  "token获取失败,请稍后再试",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":  1000,
		"data":  users,
		"msg":   "",
		"token": token,
	})
}

func UpdataUser(c *gin.Context) {
	var user models.UserInfo
	c.ShouldBind(&user)
	id := strconv.FormatUint(user.Id, 10)
	users := models.FindUserById(id)
	b := utils.ValidPassWord(user.Password, users.Salt, users.Password)
	if b {
		models.UpdataUserInfo(user)
		c.JSON(200, gin.H{
			"code": 1000,
			"msg":  "修改成功。。。",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 1001,
		"msg":  "修改失败。。。",
	})
	// UpdataUserInfo
}

func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
func SearchFriends(c *gin.Context) {
	s := c.Query("userId")
	var friend []models.UserInfo
	contact, err := models.FindFriendsById(s)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1001,
			"data": "",
			"msg":  "暂无该好友",
		})
		return
	}
	for _, v := range contact {
		id := strconv.FormatUint(v.TargetId, 10)
		if v.Type == 1 {
			friend = append(friend, models.FindUserById(id))
		}
	}
	c.JSON(200, gin.H{
		"code": 1000,
		"data": friend,
		"msg":  "",
	})
}
func SearchUserLists(c *gin.Context) {
	selectId := strings.Split(c.Query("selectId"), ",")
	userList := models.FindUserByIdList(selectId)
	c.JSON(200, gin.H{
		"code": 1000,
		"data": userList,
		"msg":  "",
	})
}
func SelectUserByPhone(c *gin.Context) {
	type Phone struct {
		Phone string `binding:"required"`
	}
	var phone Phone
	err := c.ShouldBind(&phone)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1001,
			"msg":  "请填写完整",
		})
		return
	}
	userinfo := models.FindUserByPhone(phone.Phone)
	if userinfo.Id == 0 {
		c.JSON(200, gin.H{
			"code": 1000,
			"data": "",
			"msg":  "暂无该用户",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 1000,
		"data": userinfo,
		"msg":  "",
	})
}

// 添加好友
func AddFriends_service(c *gin.Context) {
	type AddFriendsParams struct {
		MId uint64
		FId uint64
	}
	var addFriendsParams AddFriendsParams
	err := c.ShouldBind(&addFriendsParams)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1001,
			"msg":  "参数不正确",
		})
		return
	}
	mc := models.SelectFriendsState(addFriendsParams.MId, addFriendsParams.FId)
	if mc.Type == 1 {
		c.JSON(200, gin.H{
			"code": 1001,
			"msg":  "该用户已经是你好友了",
		})
		return
	}

	models.InsertRelations(addFriendsParams.MId, addFriendsParams.FId, 2)
	models.InsertRelations(addFriendsParams.FId, addFriendsParams.MId, 3)
	c.JSON(200, gin.H{
		"code": 1000,
		"msg":  "请等待对方同意",
	})
}
func DelFriends(c *gin.Context) {
	type AddFriendsParams struct {
		MId uint64
		FId uint64
	}
	var addFriendsParams AddFriendsParams
	err := c.ShouldBind(&addFriendsParams)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1001,
			"msg":  "参数不正确",
		})
		return
	}
	models.InsertRelations(addFriendsParams.MId, addFriendsParams.FId, 0)
	models.InsertRelations(addFriendsParams.FId, addFriendsParams.MId, 0)
	c.JSON(200, gin.H{
		"code": 1000,
		"msg":  "删除成功",
	})
}
func AgreeAdd(c *gin.Context) {
	type AddFriendsParams struct {
		MId uint64
		FId uint64
	}
	var addFriendsParams AddFriendsParams
	err := c.ShouldBind(&addFriendsParams)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1001,
			"msg":  "参数不正确",
		})
		return
	}
	models.InsertRelations(addFriendsParams.MId, addFriendsParams.FId, 1)
	models.InsertRelations(addFriendsParams.FId, addFriendsParams.MId, 1)
	c.JSON(200, gin.H{
		"code": 1000,
		"msg":  "添加成功",
	})
}

// 查询好友关系
func Relationship(c *gin.Context) {
	type Params struct {
		MId  uint64
		Type int
	}
	var p Params
	c.ShouldBind(&p)
	user := models.SelectFriendsRelationship(p.MId, p.Type)
	c.JSON(200, gin.H{
		"code": 1000,
		"data": user,
		"msg":  "",
	})
}
