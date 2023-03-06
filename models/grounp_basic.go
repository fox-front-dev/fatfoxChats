package models

import "gorm.io/gorm"

// 群信息
type GrounpBasic struct {
	gorm.Model
	Name    string
	OwnerId uint64
	Icon    string
	Desc    string
	Type    int //预留
}

func (table *GrounpBasic) TableName() string {
	return "grounp_basic"
}
