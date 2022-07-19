/*
 * @Author: lqc
 * @Date: 2021-11-18 09:59:38
 * @Description: 使用zap作为日志管理
 */

package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ashin-l/go-demo/pkg/logger"
	gLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// Colors
const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

type gormLog struct {
	gLogger.Config
	// Writer
	// Config
	// infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

func defaultLog(logLevel string) gLogger.Interface {

	config := gLogger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  gLogger.Info,
		IgnoreRecordNotFoundError: false,
		Colorful:                  true,
	}

	//Silent,Error,Warn,Info
	switch logLevel {
	case "silent":
		config.LogLevel = gLogger.Silent
	case "error":
		config.LogLevel = gLogger.Error
	case "warn":
		config.LogLevel = gLogger.Warn
	}

	var (
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	if config.Colorful {
		traceStr = Green + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
		traceWarnStr = Green + "%s " + Yellow + "%s\n" + Reset + RedBold + "[%.3fms] " + Yellow + "[rows:%v]" + Magenta + " %s" + Reset
		traceErrStr = RedBold + "%s " + MagentaBold + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
	}

	return &gormLog{
		Config:       config,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

// LogMode log mode
func (l *gormLog) LogMode(level gLogger.LogLevel) gLogger.Interface {
	newgormLog := *l
	newgormLog.LogLevel = level
	return &newgormLog
}

// Info print info
func (l gormLog) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gLogger.Info {
		// l.Printf(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
		logger.Logger().Info(msg, append([]interface{}{utils.FileWithLineNum()}, data...))
	}
}

// Warn print warn messages
func (l gormLog) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gLogger.Warn {
		// l.Printf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
		logger.Logger().Warn(msg, append([]interface{}{utils.FileWithLineNum()}, data...))
	}
}

// Error print error messages
func (l gormLog) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gLogger.Error {
		// l.Printf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
		logger.Logger().Error(msg, append([]interface{}{utils.FileWithLineNum()}, data...))
	}
}

// Trace print sql message
func (l gormLog) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gLogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gLogger.Error && (!errors.Is(err, gLogger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			// l.Printf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			logger.Logger().Errorf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			// l.Printf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			logger.Logger().Infof(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gLogger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			// l.Printf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			logger.Logger().Warnf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			logger.Logger().Infof(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.LogLevel == gLogger.Info:
		sql, rows := fc()
		if rows == -1 {
			logger.Logger().Infof(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			logger.Logger().Infof(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
