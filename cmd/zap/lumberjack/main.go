package main

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger

/*
 * 使用 lumberjack 实现根据等级分文件写日志,及文件切割
 */
func main() {
	_, err := os.Stat("logs")
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir("logs", os.ModePerm)
			if err != nil {
				fmt.Printf("mkdir failed![%v]\n", err)
			}
		}
	}

	// 获取编码器
	encoderconfig := zap.NewDevelopmentEncoderConfig()
	encoderconfig.EncodeTime = zapcore.ISO8601TimeEncoder        //指定时间格式
	encoderconfig.EncodeLevel = zapcore.CapitalColorLevelEncoder //按级别显示不同颜色，不需要的话取值zapcore.CapitalLevelEncoder就可以了
	// encoderconfig.EncodeCaller = zapcore.FullCallerEncoder //显示完整文件路径
	encoder := zapcore.NewConsoleEncoder(encoderconfig)

	// 日志级别
	errorLevel := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //error级别
		return lev >= zap.ErrorLevel
	})
	infoLevel := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //info和debug级别,debug级别是最低的
		return lev < zap.ErrorLevel && lev >= zap.InfoLevel
	})

	// info 文件 writeSyncer
	infoWS := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/info.log", //日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    1,               //文件大小限制,单位MB
		MaxBackups: 5,               //最大保留日志文件数量
		MaxAge:     30,              //日志文件保留天数
		// Compress:   false,           //是否压缩处理
	})
	infoCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoWS, zapcore.AddSync(os.Stdout)), infoLevel)

	// error 文件 writeSyncer
	errorWS := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/error.log",
		MaxSize:    30,
		MaxBackups: 10,
		MaxAge:     30,
	})
	errorCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorWS, zapcore.AddSync(os.Stderr)), errorLevel)
	logger = zap.New(zapcore.NewTee(infoCore, errorCore), zap.AddCaller())
	sugar := logger.Sugar()
	sugar.Infow("failed to fetch URL",
		"url", "http://example.com",
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("failed to fetch URL: %s", "http://example.com")
	sugar.Error("error!!!")
	sugar.Debug("debug...")
	fmt.Println("down")
}
