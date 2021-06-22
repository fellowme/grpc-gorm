package etcdSetup

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"grpc/tools/common_tool"
	"grpc/tools/settings"
	"strings"
	"time"
)

var etcdClient *clientv3.Client

func RegisterETCD(ttl int64) (err error) {
	if etcdClient == nil {
		etcdClient, err = clientv3.New(clientv3.Config{
			Endpoints:   strings.Split(settings.AppSetting.EtcdSetting.ServerAddress, ","),
			DialTimeout: 0,
		})
	}
	if err != nil {
		fmt.Print("clientv3.New err= ", err)
	}
	timeTicker := time.NewTicker(time.Second * time.Duration(ttl))
	go func() {
		key := common_tool.GetEtcdFmtKey(settings.AppSetting.EtcdSetting.Schema, settings.AppSetting.ServiceName, settings.AppSetting.EtcdSetting.ServerAddress)
		for {
			request, err := etcdClient.Get(context.Background(), key)
			if err != nil {
				fmt.Print("etcdClient.Get err=", err)
			} else if request.Count == 0 {
				err = keepAlive(key, ttl)
				if err != nil {
					fmt.Print("keepAlive err=", err)
				}
			}
		}
		<-timeTicker.C
	}()
	return nil
}

func keepAlive(key string, ttl int64) error {
	leaseGrantRes, err := etcdClient.Grant(context.Background(), ttl)
	if err != nil {
		fmt.Print("etcdClient.Grant 失败err=", err)
		return err
	}
	_, err = etcdClient.Put(context.Background(), key, settings.AppSetting.EtcdSetting.ServerAddress, clientv3.WithLease(leaseGrantRes.ID))
	if err != nil {
		fmt.Print("etcdClient.Put 失败err=", err)
		return err
	}
	_, err = etcdClient.KeepAlive(context.Background(), leaseGrantRes.ID)
	if err != nil {
		fmt.Print("etcdClient.KeepAlive 失败err=", err)
		return err
	}
	return nil
}

func UnRegisterETCD() {
	if etcdClient != nil {
		key := common_tool.GetEtcdFmtKey(settings.AppSetting.EtcdSetting.Schema, settings.AppSetting.ServiceName, settings.AppSetting.EtcdSetting.ServerAddress)
		_, err := etcdClient.Delete(context.Background(), key)
		if err != nil {
			fmt.Print("etcdClient.Delete err=", err)
		}
	}
}

func GetEtcdClient() *clientv3.Client {
	if etcdClient != nil {
		return etcdClient
	}
	return nil
}
