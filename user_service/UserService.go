package user_service

import (
	"context"
	"github.com/jinzhu/gorm"
	"grpc/dao"
	"grpc/service"
)

type UserService struct {
	ud dao.UserDaoInterface
}

func GetUserService(mysqlDB *gorm.DB) UserService {
	userDao := dao.GetUserDao(mysqlDB)
	return UserService{ud: &userDao}
}

func (us *UserService) GetUser(ctx context.Context, request *service.UserRequest) (*service.UserResponse, error) {
	data := us.ud.GetUserDao(request.Id)
	return &data, nil
}

func (us *UserService) CreateUser(cxt context.Context, request *service.UserRequest) (*service.MessageResponse, error) {
	err := us.ud.CreateUserDao(*request)
	if err != nil {
		return &service.MessageResponse{}, err
	}
	return &service.MessageResponse{
		Message: "创建成功",
	}, nil
}

func (us *UserService) UpdateUser(cxt context.Context, request *service.UserRequest) (*service.MessageResponse, error) {
	err := us.ud.UpdateUserDao(*request)
	if err != nil {
		return &service.MessageResponse{}, err
	}
	return &service.MessageResponse{
		Message: "更新成功",
	}, nil
}
func (us *UserService) DeleteUser(ctx context.Context, request *service.UserRequest) (*service.MessageResponse, error) {
	err := us.ud.DeleteUserDao(request.Id)
	if err != nil {
		return &service.MessageResponse{}, err
	}
	return &service.MessageResponse{
		Message: "删除成功",
	}, nil
}

func (us *UserService) UserList(ctx context.Context, request *service.UserIdListRequest) (*service.UserListResponse, error) {
	data := us.ud.UserListCountDao(*request)
	return &data, nil
}
