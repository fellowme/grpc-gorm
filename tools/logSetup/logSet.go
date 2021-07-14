package logSetup

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"grpc-gorm/tools/settings"
	"os"
)

//　初始化zaplogger日志库
var MyLogger *zap.Logger

func InitLogger() *zap.Logger {
	hook := lumberjack.Logger{
		Filename:   settings.AppSetting.ZapLoggerSetting.LogPath, // 日志文件路径
		MaxSize:    128,                                          // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,                                           // 日志文件最多保存多少个备份
		MaxAge:     7,                                            // 文件最多保存多少天
		Compress:   true,                                         // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line_num",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 短路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zapcore.Level(settings.AppSetting.ZapLoggerSetting.LevelInt))

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.WithCaller(settings.AppSetting.ZapLoggerSetting.ZapCallerFlag)
	// 开启文件及行号
	development := zap.Development()

	// 设置初始化字段
	filed := zap.Fields(zap.String("serviceName", settings.AppSetting.ServiceName))
	// 构造日志
	MyLogger = zap.New(core, caller, development, filed)

	MyLogger.Info(settings.AppSetting.ServiceName + "_log 初始化成功")
	return MyLogger
}
