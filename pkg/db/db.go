/*
 * @Author: lqc
 * @Date: 2021-11-16 13:17:58
 * @Description: 数据初始化
 */

package db

import (
	"fmt"
	"time"

	"github.com/ashin-l/go-demo/pkg/logger"
	"github.com/ashin-l/go-demo/pkg/option"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func Init(opt *option.Options) {
	var dsn string
	var err error
	ds := opt.DataSource
	switch ds.Driver {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", ds.UserName, ds.PassWord, ds.Host, ds.DbName)
	case "postgres":
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Shanghai", ds.Host, ds.Port, ds.UserName, ds.PassWord, ds.DbName, ds.Sslmode)
	default:
		err = fmt.Errorf("wrong sql type")
		logger.Logger().Fatal(err)
	}

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: defaultLog(opt.DataSource.LogLevel)})
	if err != nil {
		logger.Logger().Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Logger().Fatal(err)
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	logger.Logger().Info("数据库连接成功")
}

// Db returns a database connection.
func Db() *gorm.DB {
	return db
}
