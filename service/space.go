package service

import (
	"fatfoxChats/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 添加说说
func AddSpaceRecord(c *gin.Context) {
	var param models.SpaceContent
	err := c.ShouldBind(&param)
	if param.Content == "" {
		c.JSON(200, gin.H{
			"code": 1001,
			"msg":  "请填写内容！！！",
		})
		return
	}
	cId, _ := c.Get("Id")
	Id := cId.(uint64)
	param.UserId = Id
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1001,
			"msg":  "参数错误！！！",
		})
		return
	}
	models.AppendSpace(param)
	c.JSON(200, gin.H{
		"code": 1000,
		"msg":  "添加成功！！！",
	})
}

// 更改说说状态
func UpdateSpaceRecord(c *gin.Context) {
	var spaceComment models.SpaceContent
	c.ShouldBind(&spaceComment)
	if spaceComment.Id <= 0 {
		c.JSON(200, gin.H{
			"code": 1001,
			"msg":  "Id错误",
		})
		return
	}
	if (spaceComment.Type == 1) || (spaceComment.Type == 0) {
		id, _ := c.Get("Id")
		spaceComment.UserId = id.(uint64)
		models.UpdateSpace(spaceComment)
		c.JSON(200, gin.H{
			"code": 1000,
			"msg":  "修改成功",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 1001,
		"msg":  "Type 必填项目",
	})
}

// 获取所有好友的空间说说
func GetAllSpaceContent(c *gin.Context) {
	type Size struct {
		Limit  int
		Offset int
	}
	var size Size
	err := c.ShouldBind(&size)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1001,
			"msg":  "参数错误！！",
		})
		return
	}
	id, _ := c.Get("Id")
	SId := strconv.FormatUint(id.(uint64), 10)
	userList, err := models.FindFriendsById(SId)
	if err != nil {
		return
	}
	var Idlist []uint64
	Idlist = append(Idlist, id.(uint64))

	for _, v := range userList {
		Idlist = append(Idlist, v.TargetId)
	}
	spacesList := models.SelectAllSpace(Idlist, size.Limit, size.Offset)
	c.JSON(200, gin.H{
		"code": 1000,
		"msg":  "",
		"data": spacesList,
	})
}
