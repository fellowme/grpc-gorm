package model

import (
	"grpc/tools/base_model"
)

type User struct {
	base_model.BaseModel
	UserPhone   string `gorm:"comment:'用户电话';type:varchar(18);index:user_phone" json:"user_phone"`
	UserName    string `gorm:"comment:'用户名字';type:varchar(50)" json:"user_name"`
	UserSex     int    `gorm:"comment:'用户性别 1男 2女 3 未知';type:int(1)" json:"user_sex"`
	UserAddress string `gorm:"comment:'用户地址';type:varchar(100)" json:"user_address"`
}

func (User) TableName() string {
	return "grpc_user"
}
