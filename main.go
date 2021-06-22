package main

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"google.golang.org/grpc"
	"grpc/middleware"
	"grpc/model"
	"grpc/service"
	"grpc/tools/etcdSetup"
	"grpc/tools/mysqlSetup"
	"grpc/tools/settings"
	"grpc/user_service"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	gRpcService := createGRpc()
	defer mysqlSetup.CloseDB()
	defer gRpcService.GracefulStop()
	address := fmt.Sprintf("%s:%d", settings.AppSetting.ServiceHost, settings.AppSetting.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		middleware.MyLogger.Error(fmt.Sprintf("服务监听端口失败 %V", err))
		return
	}
	if err := etcdSetup.RegisterETCD(5); err != nil {
		middleware.MyLogger.Error(fmt.Sprintf("注册失败 %V", err))
		return
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		etcdSetup.UnRegisterETCD()
		if i, ok := s.(syscall.Signal); ok {
			middleware.MyLogger.Info("grpc server 终止启动")
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}
	}()
	err = gRpcService.Serve(listener)
	middleware.MyLogger.Info("grpc server 启动")
	if err != nil {
		middleware.MyLogger.Error(fmt.Sprintf("gRpcService服务监听端口失败 %V", err))
		return
	}

}

func createGRpc() *grpc.Server {
	settings.SettingSetUp()
	zapLogger := middleware.InitLogger("D:\\grpc-gin-gorm\\log\\grpc.log", 0)

	db := mysqlSetup.SetUp()

	db.AutoMigrate(&model.User{})

	UserService := user_service.GetUserService(db)
	gRpcService := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(zapLogger),
			middleware.UnaryServerInterceptor(),
		)))
	service.RegisterUserServiceServer(gRpcService, &UserService)

	return gRpcService
}
