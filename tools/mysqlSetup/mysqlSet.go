package mysqlSetup

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"grpc/middleware"
	"grpc/tools/settings"
)

var mysqlDb *gorm.DB

func SetUp() *gorm.DB {
	msg := fmt.Sprintf("mysql.Setup 开始链接数据库")
	middleware.MyLogger.Info(msg)
	var err error
	mysqlDb, err = gorm.Open(settings.AppSetting.MysqlSetting.DbDriverName, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		settings.AppSetting.MysqlSetting.DbUserName,
		settings.AppSetting.MysqlSetting.DbPassword,
		fmt.Sprintf("%s:%d", settings.AppSetting.MysqlSetting.DbHost, settings.AppSetting.MysqlSetting.DbPort),
		settings.AppSetting.MysqlSetting.DbName))

	if err != nil {
		msg := fmt.Sprintf("mysql.Setup err= %v", err)
		middleware.MyLogger.Error(msg)
	}
	mysqlDb.DB().SetMaxIdleConns(settings.AppSetting.MysqlSetting.MaxIdleConnects)
	mysqlDb.DB().SetMaxOpenConns(settings.AppSetting.MysqlSetting.MaxOpenConnects)
	mysqlDb.LogMode(settings.AppSetting.MysqlSetting.LogMode)
	msg = fmt.Sprintf("mysql.Setup 链接数据库成功")
	middleware.MyLogger.Info(msg)
	return mysqlDb
}

func CloseDB() {
	if mysqlDb != nil {
		if err := mysqlDb.Close(); err != nil {
			middleware.MyLogger.Error("mysql 关闭错误！")
		}
	}
}
