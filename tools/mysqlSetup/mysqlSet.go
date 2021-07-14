package mysqlSetup

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"grpc-gorm/tools/logSetup"
	"grpc-gorm/tools/settings"
)

var db *gorm.DB

func SetUp() *gorm.DB {
	msg := fmt.Sprintf("mysql.Setup 开始链接数据库")
	logSetup.MyLogger.Info(msg)
	var err error
	db, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		settings.AppSetting.MysqlSetting.DbUserName,
		settings.AppSetting.MysqlSetting.DbPassword,
		fmt.Sprintf("%s:%d", settings.AppSetting.MysqlSetting.DbHost, settings.AppSetting.MysqlSetting.DbPort),
		settings.AppSetting.MysqlSetting.DbName)), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		msg := fmt.Sprintf("mysql.Setup err= %v", err)
		logSetup.MyLogger.Error(msg)
	}
	mysqlDB, err := db.DB()
	if err != nil {
		msg := fmt.Sprintf("mysql.Setup err= %v", err)
		logSetup.MyLogger.Error(msg)
	}
	mysqlDB.SetMaxIdleConns(settings.AppSetting.MysqlSetting.MaxIdleConnects)
	mysqlDB.SetMaxOpenConns(settings.AppSetting.MysqlSetting.MaxOpenConnects)
	msg = fmt.Sprintf("mysql.Setup 链接数据库成功")
	logSetup.MyLogger.Info(msg)
	return db
}

func CloseDB() {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			logSetup.MyLogger.Error(fmt.Sprintf("mysql db.DB() 关闭错误！err=%v", err))
			return
		}
		if err := sqlDB.Close(); err != nil {
			logSetup.MyLogger.Error(fmt.Sprintf("mysql sqlDB.Close() 关闭错误！err=%v", err))
		}
	}
}
