package models

import (
	"fatfoxChats/utils"

	"gorm.io/gorm"
)

type UserInfo struct {
	Id            uint64 `gorm:"primary_key;AUTO_INCREMENT" `
	AccName       string `gorm:"size:16"  binding:"required" `
	Name          string `gorm:"size:10"  binding:"required" `
	Phone         string `gorm:"size:11"  binding:"required"`
	Password      string `gorm:"size:128"  binding:"required" json:"-"`
	Sex           string `gorm:"size:1"   binding:"required"`
	Email         string `gorm:"size:128" binding:"required" `
	HeadPortrait  string
	Identity      string `json:"-"`
	ClientIp      string `gorm:"size:128"`
	LoginOutTime  uint64
	HeartbeatTime uint64 `json:"-"`
	LoginTime     uint64
	IsLogout      bool
	IsAdmin       bool   `json:"-"`
	Salt          string `gorm:"default:0" json:"-"`
}

func (table *UserInfo) TableName() string {
	return "user_info"
}

// 注册用户
func CreateUser(user UserInfo) *gorm.DB {
	return utils.DB.Create(&user)
}

// 注销用户
func DelectUser(uid uint64, pwd string, phone string) *gorm.DB {
	var user UserInfo
	return utils.DB.Where("id=? and password=? and phone=?", uid, pwd, phone).Delete(&user)
}

// 更新用户信息
func UpdataUserInfo(u UserInfo) *gorm.DB {
	var user UserInfo
	return utils.DB.Model(&user).Where("Id=?", u.Id).Updates(u)
}

// 查找用户根据手机号码
func FindUserByPhone(phone string) UserInfo {
	var user UserInfo
	utils.DB.Where("phone = ?", phone).First(&user)
	return user
}

// 查找用户根据id
func FindUserById(u string) UserInfo {
	var user UserInfo
	utils.DB.Where(" id=?", u).First(&user)
	return user
}

// 根据 idList 查找多位好友
func FindUserByIdList(Id []string) []UserInfo {
	var user []UserInfo
	utils.DB.Where("id IN ?", Id).Find(&user)
	return user
}

// 登录状态
func IsLoginState(Id string, State bool) *gorm.DB {
	var user UserInfo
	return utils.DB.Model(&user).Where("id = ?", Id).Updates(map[string]interface{}{"IsLogout": State})
}
