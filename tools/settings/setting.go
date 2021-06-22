package settings

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type etcdConfig struct {
	ServerAddress string
	Schema        string
	TTL           int64
}

type mysqlConfig struct {
	DbUserName   string
	DbPassword   string
	DbHost       string
	DbPort       int
	DbName       string
	DbDriverName string
}

type appConfig struct {
	ServiceName   string
	ServiceHost   string
	Port          int
	ZapCallerFlag bool
	EtcdSetting   etcdConfig
	MysqlSetting  mysqlConfig
}

var AppSetting appConfig

func SettingSetUp() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("read config failed: %v", err)
	}
	if err := viper.Unmarshal(&AppSetting); err != nil {
		fmt.Print("viper.Unmarshal err=", err)
	}

}
