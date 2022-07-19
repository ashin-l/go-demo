package logger

import (
	"os"

	"github.com/ashin-l/go-demo/pkg/option"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var defaultLogger *zap.SugaredLogger

func Init(opt *option.Options) {
	logLevel := zap.InfoLevel
	if opt.Level == "debug" {
		logLevel = zap.DebugLevel
	}
	// 获取编码器
	// encoderconfig := zap.NewDevelopmentEncoderConfig()
	encoderconfig := zap.NewProductionEncoderConfig()
	encoderconfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05") //指定时间格式
	encoderconfig.EncodeLevel = zapcore.CapitalColorLevelEncoder                  //按级别显示不同颜色，不需要的话取值zapcore.CapitalLevelEncoder就可以了
	encoderconfig.EncodeCaller = zapcore.FullCallerEncoder                        //显示完整文件路径
	encoder := zapcore.NewConsoleEncoder(encoderconfig)

	// 日志级别
	infoLevel := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //info和debug级别,debug级别是最低的
		return lev < zap.ErrorLevel && lev >= logLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //error级别
		return lev >= zap.ErrorLevel
	})

	// info 文件 writeSyncer
	infoWS := zapcore.AddSync(&lumberjack.Logger{
		Filename:   opt.Path + "/info.log", //日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    10,                     //文件大小限制,单位MB
		MaxBackups: 100,                    //最大保留日志文件数量
		MaxAge:     30,                     //日志文件保留天数
		// Compress:   false,           //是否压缩处理
	})
	infoCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoWS, zapcore.AddSync(os.Stdout)), infoLevel)

	// error 文件 writeSyncer
	errorWS := zapcore.AddSync(&lumberjack.Logger{
		Filename:   opt.Path + "/error.log",
		MaxSize:    30,
		MaxBackups: 10,
		MaxAge:     30,
	})
	errorCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorWS, zapcore.AddSync(os.Stderr)), errorLevel)
	defaultLogger = zap.New(zapcore.NewTee(infoCore, errorCore), zap.AddCaller()).Sugar()
	defaultLogger.Info("log 初始化成功")
}

func Logger() *zap.SugaredLogger {
	return defaultLogger
}
