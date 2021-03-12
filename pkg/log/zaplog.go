package log

import (
	"fmt"
	"io"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLog(logPath, errPath, logLevel string) *zap.Logger {
	al := zap.NewAtomicLevel()
	al.UnmarshalText([]byte(logLevel))
	// config := zapcore.EncoderConfig{
	// MessageKey:   "msg",  //结构化（json）输出：msg的key
	// LevelKey:     "level",//结构化（json）输出：日志级别的key（INFO，WARN，ERROR等）
	// TimeKey:      "ts",   //结构化（json）输出：时间的key（INFO，WARN，ERROR等）
	// CallerKey:    "file", //结构化（json）输出：打印日志的文件对应的Key
	// EncodeLevel:  zapcore.CapitalLevelEncoder, //将日志级别转换成大写（INFO，WARN，ERROR等）
	// EncodeCaller: zapcore.ShortCallerEncoder, //采用短文件路径编码输出（test/main.go:14 ）
	// EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	//     enc.AppendString(t.Format("2006-01-02 15:04:05"))
	// },//输出的时间格式
	// EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
	//     enc.AppendInt64(int64(d) / 1000000)
	// },//
	// }
	//自定义日志级别：自定义Info级别
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel && lvl >= al.Level()
	})
	//自定义日志级别：自定义Error级别
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel && lvl >= al.Level()
	})
	// 获取io.Writer的实现
	infoWriter := getWriter(logPath)
	errorWriter := getWriter(errPath)
	// 实现多个输出
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()), zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()), zapcore.AddSync(errorWriter), errorLevel),
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), infoLevel), //同时将日志输出到控制台，NewJSONEncoder 是结构化输出
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), zapcore.AddSync(os.Stdout), errorLevel),                             //同时将日志输出到控制台，NewJSONEncoder 是结构化输出
		// zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.AddSync(infoWriter), infoLevel), //将info及以下写入logPath，NewConsoleEncoder 是非结构化输出
		// zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.AddSync(warnWriter), warnLevel),//warn及以上写入errPath
		// zapcore.NewCore(zapcore.NewJSONEncoder(config), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), logLevel),//同时将日志输出到控制台，NewJSONEncoder 是结构化输出
	)
	// logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel))
	// sugarLogger = logger.Sugar()
	return zap.New(core)
	// sugar = zap.New(core).Sugar()
}

func getWriter(filename string) io.Writer {
	hook, err := rotatelogs.New(
		filename+".%Y%m%d%H",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*2),    // 保存30天
		rotatelogs.WithRotationTime(time.Hour*4), //切割频率 24小时
	)
	if err != nil {
		fmt.Println("日志启动异常")
		panic(err)
	}
	return hook
}
