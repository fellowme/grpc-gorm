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
	DialTimeout   int
}

type mysqlConfig struct {
	DbUserName      string
	DbPassword      string
	DbHost          string
	DbPort          int
	DbName          string
	DbDriverName    string
	MaxOpenConnects int
	MaxIdleConnects int
	LogMode         bool
}

type appConfig struct {
	ServiceName   string
	ServiceHost   string
	Port          int
	ZapCallerFlag bool
	Weight        string
	LogPath       string
	LevelInt      int
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
