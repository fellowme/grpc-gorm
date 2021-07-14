package dao

import (
	"gorm.io/gorm"
	"grpc-gorm/model"
	"grpc-gorm/service"
)

type (
	UserDaoInterface interface {
		GetUserDao(userId int32) service.UserResponse
		CreateUserDao(user service.UserRequest) error
		UpdateUserDao(user service.UserRequest) error
		DeleteUserDao(userId int32) error
		UserListCountDao(param service.UserIdListRequest) service.UserListResponse
	}
)

type UserDao struct {
	db *gorm.DB
}

func GetUserDao(mysqlDB *gorm.DB) UserDao {
	return UserDao{db: mysqlDB}
}

func (ud UserDao) GetUserDao(userId int32) service.UserResponse {
	var userInfo service.UserResponse
	ud.db.Model(&model.User{}).Where("id = ?", userId).First(&userInfo)
	return userInfo
}

func (ud UserDao) CreateUserDao(user service.UserRequest) error {
	if err := ud.db.Table("grpc_user").Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (ud UserDao) UpdateUserDao(user service.UserRequest) error {
	if err := ud.db.Table("grpc_user").Where("id=?", user.Id).Updates(&user).Error; err != nil {
		return err
	}
	return nil
}

func (ud UserDao) DeleteUserDao(userId int32) error {
	if err := ud.db.Model(&model.User{}).Where("id=?", userId).Update("is_deleted", true).Error; err != nil {
		return err
	}
	return nil
}

func (ud UserDao) UserListCountDao(param service.UserIdListRequest) (userCountList service.UserListResponse) {
	var userList []*service.UserResponse
	var totalCount int64
	userDb := ud.db.Table("grpc_user")
	if len(param.UserId) != 0 {
		userDb.Where("id in ? ", param.UserId)
	}
	userDb.Count(&totalCount)
	userDb.Offset(int((param.Page - 1) * param.PageSize)).Limit(int(param.PageSize)).Find(&userList)
	userCountList.TotalCount = totalCount
	userCountList.UserList = userList
	return userCountList
}
