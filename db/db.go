package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	DB_USER    = "user"
	DB_PASS    = "password"
	DB_HOST    = "127.0.0.1"
	DB_PORT    = "13306"
	DB_NAME    = "gqlgentodos"
	DB_OPTIONS = "charset=utf8mb4&parseTime=True&loc=Local"
)

type SQLHandler struct {
	*gorm.DB
}

// DB接続をする
func Open() (*SQLHandler, error) {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		DB_USER, DB_PASS,
		DB_HOST, DB_PORT,
		DB_NAME, DB_OPTIONS,
	)

	conf := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // SQLログを出力する
	}

	db, err := gorm.Open(mysql.Open(dsn), conf)
	return &SQLHandler{DB: db}, err
}
