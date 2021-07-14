package main

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"google.golang.org/grpc"
	"grpc-gorm/middleware/deadline"
	"grpc-gorm/model"
	"grpc-gorm/service"
	"grpc-gorm/tools/etcdSetup"
	"grpc-gorm/tools/logSetup"
	"grpc-gorm/tools/mysqlSetup"
	"grpc-gorm/tools/settings"
	"grpc-gorm/user_service"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	gRpcService := createGRpc()
	defer mysqlSetup.CloseDB()
	defer gRpcService.GracefulStop()
	address := fmt.Sprintf("%s", settings.AppSetting.ServiceHost)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logSetup.MyLogger.Error(fmt.Sprintf("服务监听端口失败 %V", err))
		return
	}
	if err := etcdSetup.RegisterETCD(settings.AppSetting.EtcdSetting.TTL); err != nil {
		logSetup.MyLogger.Error(fmt.Sprintf("注册失败 %V", err))
		return
	}
	logSetup.MyLogger.Debug(fmt.Sprintf("etcdSetup 注册成功"))
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		etcdSetup.UnRegisterETCD()
		if i, ok := s.(syscall.Signal); ok {
			logSetup.MyLogger.Info("grpc server 终止启动")
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}
	}()
	err = gRpcService.Serve(listener)
	logSetup.MyLogger.Info("grpc server 启动")
	if err != nil {
		logSetup.MyLogger.Error(fmt.Sprintf("gRpcService服务监听端口失败 %V", err))
		return
	}

}

func createGRpc() *grpc.Server {
	settings.SettingSetUp()
	zapLogger := logSetup.InitLogger(settings.AppSetting.LogPath, settings.AppSetting.LevelInt)

	db := mysqlSetup.SetUp()
	db.AutoMigrate(&model.User{})
	db.Migrator().CreateIndex(&model.User{}, "UserPhone")
	UserService := user_service.GetUserService(db)
	gRpcService := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(zapLogger),
			deadline.UnaryServerInterceptor(),
		)))
	service.RegisterUserServiceServer(gRpcService, &UserService)

	return gRpcService
}
