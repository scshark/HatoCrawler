/**
 * @Author: scshark
 * @Description:
 * @File:  db
 * @Date: 12/9/22 3:06 PM
 */
package conf

import (
	. "SecCrawler/config"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
	"sync"
	"time"
)

var (
	db    *gorm.DB
	once  sync.Once
	Redis *redis.Client
)
var Ctx = context.Background()

func MustGormDb() *gorm.DB {
	once.Do(func() {
		var err error
		if db, err = newDBEngine(); err != nil {
			logrus.Fatalf("new gorm db failed: %s", err)
		}
	})
	return db
}

/**
 * 初始化 gorm db Engine
 * @author: sc.shark
 * @date: Created in 2022-12-09 18:01:19
 * @Description:
 * @return db
 * @return err
 */
func newDBEngine() (db *gorm.DB, err error) {

	// 初始化 DB日志配置
	newLogger := logger.New(
		logrus.StandardLogger(), // io writer（日志输出的目标，前缀和日志包含的内容）
		logger.Config{
			SlowThreshold:             time.Second,  // 慢查询
			IgnoreRecordNotFoundError: true,         // 忽略ErrRecordNotFound（记录未找到）错误
			LogLevel:                  logger.Error, // 日志级别
			Colorful:                  false,        // 禁用彩色打印
		})

	gormConfig := &gorm.Config{
		Logger: newLogger, // 日志
		NamingStrategy: schema.NamingStrategy{ // 覆盖默认的NamingStrategy来更改命名约定
			TablePrefix:   Cfg.Database.TablePrefix, // 表前缀
			SingularTable: true,                     // 使用小写表名, table for `User` would be `user` with this option enabled
		},
	}

	// 注册连接池 DBResolver 为 GORM 提供了多个数据库支持，
	dbPlugin := dbresolver.Register(dbresolver.Config{}). // xxx多个连接
								SetConnMaxIdleTime(time.Hour).              // 设置连接最大空闲时间
								SetConnMaxLifetime(24 * time.Hour).         // 设置连接最大生命周期
								SetMaxIdleConns(Cfg.Database.MaxIdleConns). // 设置最大空闲连接数
								SetMaxOpenConns(Cfg.Database.MaxOpenConns)  // 设置最大打开连接数

	// Open initialize db ， 建立mysql连接，
	if db, err = gorm.Open(mysql.Open(Dsn()), gormConfig); err == nil {
		// Use use plugin
		_ = db.Use(dbPlugin)
	}
	return db, err
}

/**
 * 返回mysql连接完整信息
 * @author: sc.shark
 * @date: Created in 2022-12-09 18:02:15
 * @Description:
 * @return string
 */
func Dsn() string {

	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
		Cfg.Database.UserName,
		Cfg.Database.Password,
		Cfg.Database.Host,
		Cfg.Database.DBName,
		Cfg.Database.Charset,
		Cfg.Database.ParseTime,
	)

}

func SetupRedisEngine() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     Cfg.Redis.Host,
		Password: Cfg.Redis.Password,
		DB:       Cfg.Redis.DB,
	})
}
