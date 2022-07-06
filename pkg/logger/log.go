package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//在性能很好但不是很关键的上下文中，使用SugaredLogger,支持结构化和printf风格的日志记录
//在每一微秒和每一次内存分配都很重要的上下文中，使用Logger,它甚至比SugaredLogger更快，内存分配次数也更少，但它只支持强类型的结构化日志记录。
var logger *zap.SugaredLogger

func Init() error {
	//两个判断日志等级的interface
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	infoWriteSyncer := getWriter("./log/info.log")
	errorWriteSyncer := getWriter("./log/error.log")
	encoder := getEncoder()
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, infoWriteSyncer, infoLevel),
		zapcore.NewCore(encoder, errorWriteSyncer, errorLevel),
	)
	//添加调用函数信息
	log := zap.New(core, zap.AddCaller())
	logger = log.Sugar()
	logger.Info("logger init success")
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
func getWriter(filename string) zapcore.WriteSyncer {
	//使用Lumberjack进行日志切割归档
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    1,     //在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: 5,     //保留旧文件的最大个数
		MaxAge:     30,    //保留旧文件的最大天数
		Compress:   false, //是否压缩/归档旧文件
	}
	return zapcore.AddSync(lumberJackLogger)
}

//定义外部可直接访问的函数（首字母大写）
func Sync() {
	logger.Sync()
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

func DPanic(args ...interface{}) {
	logger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	logger.DPanicf(template, args...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	logger.Panicf(template, args...)
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	logger.Fatalf(template, args...)
}
