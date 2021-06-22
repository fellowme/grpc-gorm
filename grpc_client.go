package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
	"grpc/service"
	"grpc/tools/etcdSetup"
	"log"
	"time"
)

func main() {
	// 1. 新建连接，端口是服务端开放的8082端口
	// 并且添加grpc.WithInsecure()，不然没有证书会报错
	// 注册 etcd
	client := etcdSetup.GetEtcdClient()
	builder := &etcdSetup.Builder{
		Client: client,
	}
	resolver.Register(builder)
	// 设置超时间
	ctx := context.Background()
	clientDeadline := time.Now().Add(time.Duration(100) * time.Millisecond)
	ctx, cancel := context.WithDeadline(ctx, clientDeadline)
	defer cancel()
	// 负载均衡
	conn, err := grpc.DialContext(ctx, "127.0.0.1:8002",
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	// 退出时关闭链接
	defer conn.Close()

	// 2. User.pb.go中的NewUserServiceClient方法
	productServiceClient := service.NewUserServiceClient(conn)

	// 3. 直接像调用本地方法一样调用GetUserInfo方法
	resp, err := productServiceClient.CreateUser(ctx, &service.UserRequest{
		UserPhone:   "15710",
		UserName:    "上海",
		UserSex:     "2",
		UserAddress: "河南",
	})
	//resp, err := productServiceClient.UpdateUser(context.Background(), &service.UserRequest{
	//	Id:          1,
	//	UserPhone:   "789",
	//	UserName:    "555",
	//	UserSex:     "2",
	//	UserAddress: "heanan ",
	//})
	//resp, err := productServiceClient.UserList(context.Background(), &service.UserIdListRequest{
	//	Page: 1, PageSize: 20})
	//resp, err := productServiceClient.DeleteUser(context.Background(), &service.UserRequest{
	//	Id: 3,
	//})
	//resp, err := productServiceClient.UserList(context.Background(), &service.UserIdListRequest{
	//	Page: 1, PageSize: 2})
	if err != nil {
		log.Fatal("调用gRPC方法错误: ", err)
	}
	fmt.Println("调用gRPC方法成功，ProdStock = ", resp)

}
