package util

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//在性能很好但不是很关键的上下文中，使用SugaredLogger,支持结构化和printf风格的日志记录
//在每一微秒和每一次内存分配都很重要的上下文中，使用Logger,它甚至比SugaredLogger更快，内存分配次数也更少，但它只支持强类型的结构化日志记录。
var Logger *zap.SugaredLogger

func InitLogger() error {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	//添加调用函数信息
	logger := zap.New(core, zap.AddCaller)
	Logger = logger.Sugar()
	return nil
}

//获取编码器
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	//修改时间编码器
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	//使用大写字母记录日志级别
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

//指定日志将写到哪里去，
func getLogWriter() zapcore.WriteSyncer {
	//使用Lumberjack进行日志切割归档
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./test.log",
		MaxSize:    1,     //在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: 5,     //保留旧文件的最大个数
		MaxAge:     30,    //保留旧文件的最大天数
		Compress:   false, //是否压缩/归档旧文件
	}
	return zapcore.AddSync(lumberJackLogger)
}
