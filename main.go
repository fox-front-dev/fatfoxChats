package main

import (
	"fatfoxChats/router"
	"fatfoxChats/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMysql()
	utils.InitRedis()
	// utils.DB.AutoMigrate(&models.UserInfo{})
	// utils.DB.AutoMigrate(&models.Contact{})
	//utils.DB.AutoMigrate(&models.Message{})
	//utils.DB.AutoMigrate(&models.GrounpBasic{})
	// utils.DB.AutoMigrate(&models.SpaceContent{})
	// utils.DB.AutoMigrate(&models.SpaceComment{})
	// utils.DB.AutoMigrate(&models.SecondLevelComment{})
	r := router.Router()
	r.Run(":4000")
}
