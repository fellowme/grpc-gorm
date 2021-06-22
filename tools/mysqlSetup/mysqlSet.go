package mysqlSetup

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"grpc/middleware"
	"grpc/tools/settings"
	"log"
	"time"
)

var mysqlDb *gorm.DB

func SetUp() *gorm.DB {
	var err error
	mysqlDb, err = gorm.Open(settings.AppSetting.MysqlSetting.DbDriverName, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		settings.AppSetting.MysqlSetting.DbUserName,
		settings.AppSetting.MysqlSetting.DbPassword,
		fmt.Sprintf("%s:%d", settings.AppSetting.MysqlSetting.DbHost, settings.AppSetting.MysqlSetting.DbPort),
		settings.AppSetting.MysqlSetting.DbName))

	if err != nil {
		fmt.Println("err:", err)
		log.Fatalf("mysql.Setup err: %v", err)
	}
	mysqlDb.DB().SetMaxIdleConns(100)
	mysqlDb.DB().SetMaxOpenConns(200)
	mysqlDb.Callback().Create().Replace("gorm:before_create", updateTimeStampForCreateCallback)
	mysqlDb.Callback().Update().Replace("gorm:before_update", updateTimeStampForUpdateCallback)
	mysqlDb.LogMode(true)
	return mysqlDb
}

func CloseDB() {
	if mysqlDb != nil {
		if err := mysqlDb.Close(); err != nil {

		}

	}
}

func updateTimeStampForCreateCallback(scope *gorm.Scope) {

	if !scope.HasError() {
		nowTime := time.Now().Unix()
		ok := scope.HasColumn("create_time")
		if !ok {
			if err := scope.SetColumn("create_time", nowTime); err != nil {
				middleware.MyLogger.Error(fmt.Sprintf("%v", err))
			}
		}

		if ok := scope.HasColumn("update_time"); !ok {
			scope.SetColumn("update_time", nowTime)
		}
	}
}

// 注册更新钩子在持久化之前
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	nowTime := time.Now()
	if ok := scope.HasColumn("update_time"); !ok {
		scope.SetColumn("update_time", nowTime)
	}
}
