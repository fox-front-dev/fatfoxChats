package models

import (
	"fatfoxChats/utils"
	"fmt"

	"gorm.io/gorm"
)

// 人员关系
type Contact struct {
	gorm.Model
	OwnerId  uint64 //谁的关系
	TargetId uint64 //对应的人员id
	Type     int    //对应的类型0未添加,1已添加,2申请中,3等待同意
	Desc     string //预留字段
}

func (table *Contact) TableName() string {
	return "contact"
}

// 根据id list寻找好友的id
func FindFriendsById(id string) ([]Contact, error) {
	var contact []Contact
	if err := utils.DB.Where("owner_id = ? and type=1", id).Find(&contact).Error; err != nil {
		return contact, err
	}
	return contact, nil
}

// 判断是否是好友
func IsFriends(MId string, FId string) bool {
	var contact Contact
	utils.DB.Where("owner_id=? and target_id=?", MId, FId).First(&contact)
	if contact.Type == 1 {
		return true
	}
	return false
}

// 查询好友关系
func SelectFriendsState(MId uint64, FId uint64) Contact {
	var contact Contact
	utils.DB.Where("owner_id=? and target_id=? ", MId, FId).First(&contact)
	return contact
}

// 查询申请好友关系表
func SelectFriendsRelationship(MId uint64, Type int) []Contact {
	var contact []Contact
	utils.DB.Where("owner_id=? and type=? ", MId, Type).Find(&contact)
	fmt.Println(contact)
	return contact
}

// 更新好友关系
func InsertRelations(OwnerId uint64, TargetId uint64, state int) *gorm.DB {
	var contact Contact
	utils.DB.Where("owner_id = ? and target_id=?", OwnerId, TargetId).First(&contact)
	contact.OwnerId = OwnerId
	contact.TargetId = TargetId
	contact.Type = state
	return utils.DB.Model(&contact).Where("id=?", contact.ID).Save(&contact)
}
