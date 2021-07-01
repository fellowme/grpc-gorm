package etcdSetup

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"grpc/middleware"
	"grpc/tools/common_tool"
	"grpc/tools/settings"
	"strings"
	"time"
)

var etcdClient *clientv3.Client

func RegisterETCD(ttl int64) (err error) {
	msg := fmt.Sprintf("RegisterETCD 开始注册 etcd")
	middleware.MyLogger.Debug(msg)
	if etcdClient == nil {
		etcdClient, err = clientv3.New(clientv3.Config{
			Endpoints:   strings.Split(settings.AppSetting.EtcdSetting.ServerAddress, ","),
			DialTimeout: time.Duration(settings.AppSetting.EtcdSetting.DialTimeout),
		})
	}

	if err != nil {
		msg := fmt.Sprintf("etcd clientv3.New err= %v", err)
		middleware.MyLogger.Error(msg)
	}
	msg = fmt.Sprintf("RegisterETCD etcdClient 创建链接客户端")
	middleware.MyLogger.Debug(msg)
	timeTicker := time.NewTicker(time.Second * time.Duration(ttl))
	go func() {
		key := common_tool.GetEtcdFmtKey(settings.AppSetting.EtcdSetting.Schema, settings.AppSetting.ServiceName, settings.AppSetting.ServiceHost)
		for _ = range timeTicker.C {
			err := keepAlive(key, ttl)
			if err != nil {
				msg := fmt.Sprintf("etcd keepAlive err= %v", err)
				middleware.MyLogger.Error(msg)
			}
		}
	}()
	return nil
}

func keepAlive(key string, ttl int64) error {
	msg := fmt.Sprintf("RegisterETCD keepAlive 开始执行")
	middleware.MyLogger.Debug(msg)
	clientDeadline := time.Now().Add(time.Duration(settings.AppSetting.EtcdSetting.DialTimeout) * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
	defer cancel()
	request, err := etcdClient.Get(ctx, key)
	if err != nil {
		msg := fmt.Sprintf("etcd clientv3.Get err= %v", err)
		middleware.MyLogger.Error(msg)
	} else if request.Count == 0 {
		leaseGrantRes, err := etcdClient.Grant(ctx, ttl)
		msg := fmt.Sprintf("RegisterETCD Grant 开始执行")
		middleware.MyLogger.Debug(msg)
		if err != nil {
			msg := fmt.Sprintf("etcd etcdClient.Grant 失败err= %v", err)
			middleware.MyLogger.Error(msg)
			return err
		}
		_, err = etcdClient.Put(ctx, key, settings.AppSetting.Weight, clientv3.WithLease(leaseGrantRes.ID))
		msg = fmt.Sprintf("RegisterETCD etcdClient Put key= %s, value=%s", key, settings.AppSetting.Weight)
		middleware.MyLogger.Debug(msg)
		if err != nil {
			msg := fmt.Sprintf("etcd etcdClient.Put 失败err= %v", err)
			middleware.MyLogger.Error(msg)
			return err
		}
		_, err = etcdClient.KeepAlive(ctx, leaseGrantRes.ID)
		msg = fmt.Sprintf("RegisterETCD KeepAlive leaseGrantRes.ID=%v", leaseGrantRes.ID)
		middleware.MyLogger.Debug(msg)
		if err != nil {
			msg := fmt.Sprintf("etcd etcdClient.KeepAlive 失败err= %v", err)
			middleware.MyLogger.Error(msg)
			return err
		}
	}

	return nil
}

func UnRegisterETCD() {
	if etcdClient != nil {
		msg := fmt.Sprintf("etcdClient UnRegisterETCD 开始取消 etcd")
		middleware.MyLogger.Info(msg)
		key := common_tool.GetEtcdFmtKey(settings.AppSetting.EtcdSetting.Schema, settings.AppSetting.ServiceName, settings.AppSetting.ServiceHost)
		clientDeadline := time.Now().Add(time.Duration(settings.AppSetting.EtcdSetting.DialTimeout) * time.Second)
		ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
		defer cancel()
		_, err := etcdClient.Delete(ctx, key)
		if err != nil {
			msg := fmt.Sprintf("etcd etcdClient.Delete err= %v", err)
			middleware.MyLogger.Error(msg)
		}
		msg = fmt.Sprintf("etcdClient UnRegisterETCD 取消 etcd 成功")
		middleware.MyLogger.Info(msg)
	}
}
