package models

import (
	"fatfoxChats/utils"

	"gorm.io/gorm"
)

// 空间评论主体
type SpaceContent struct {
	Id       uint64 `gorm:"primary_key;AUTO_INCREMENT" `
	Count    int    `json:"-"`
	Type     int    `json:"-" default:"1"` //0代表删除 1代表正常
	UserId   uint64
	Content  string
	ImageUrl string
	gorm.Model
}

// 空间评论主体下的一级评论
type SpaceComment struct {
	Id      uint64 `gorm:"primary_key;AUTO_INCREMENT" `
	Content string ` binding:"required"`
	Root    uint64 ` binding:"required"`
	Count   int
	gorm.Model
}

// 空间评论主体下的二级评论
type SecondLevelComment struct {
	Id           uint64 `gorm:"primary_key;AUTO_INCREMENT" `
	Content      string ` binding:"required"`
	CommentParse uint64 ` binding:"required"` // 指向一级评论
	Root         uint64 // 指向空间评论主体
	gorm.Model
}

func (table *SpaceComment) TableName() string {
	return "space_comment"
}

func (table *SpaceContent) TableName() string {
	return "space_content"
}
func (table *SecondLevelComment) TableName() string {
	return "second_level_comment"
}

// 添加说说
func AppendSpace(scontent SpaceContent) *gorm.DB {
	return utils.DB.Create(&scontent)
}

// 更新说说
func UpdateSpace(s SpaceContent) *gorm.DB {
	var mapVal = map[string]interface{}{
		"type": s.Type,
	}
	return utils.DB.Model(&s).Where("id=? and user_id=?", s.Id, s.UserId).Updates(mapVal)
}

// 查询说说
func SelectAllSpace(list []uint64, limit int, offset int) []SpaceContent {
	var spaces []SpaceContent
	utils.DB.Where("user_id in ?", list).Limit(limit).Offset(offset).Find(&spaces)
	return spaces
}
