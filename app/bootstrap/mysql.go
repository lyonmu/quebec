package bootstrap

import (
	"fmt"
	"time"

	commonconfig "github.com/lyonmu/quebec/app/common/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitializeDB(config *commonconfig.MySqlConfig) (*gorm.DB, error) {

	// 构建 DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"+
		"&timeout=%ds"+ // 连接超时时间
		"&readTimeout=%ds"+ // 读取超时时间
		"&writeTimeout=%ds", // 写入超时时间
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.TimeOut, // timeout
		config.TimeOut, // readTimeout
		config.TimeOut, // writeTimeout
	)

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %v", err)
	}

	// 获取底层的 sqlDB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)                  // 设置空闲连接池中的最大连接数
	sqlDB.SetMaxOpenConns(100)                 // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Hour)        // 设置连接的最大生命周期
	sqlDB.SetConnMaxIdleTime(time.Minute * 10) // 设置空闲连接的最大生命周期

	return db, nil
}
