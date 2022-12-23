package gorm

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func Init(conf ConfigDB) (*gorm.DB, *sql.DB) {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local", conf.User, conf.Pass,
		conf.Host, conf.DatabaseName, conf.Charset)

	driverConfig := mysql.New(mysql.Config{
		DSN:                       dsn,   // data source name
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	})

	gormConfig := &gorm.Config{
		PrepareStmt: true,
		Logger:      newLogger,
	}

	provider, err := gorm.Open(driverConfig, gormConfig)
	if err != nil {
		log.Fatalf("model not connected : %s", err.Error())
	}

	sqlDB, err := provider.DB()
	if err != nil {
		log.Fatalf("model not connected : %s", err.Error())
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("failed checking connection to database, err: %s", err.Error())
	}

	sqlDB.SetMaxIdleConns(conf.PoolMaxIdleConn)
	sqlDB.SetMaxOpenConns(conf.PoolMaxOpenConn)
	sqlDB.SetConnMaxLifetime(conf.PoolMaxConnLifetime)

	log.Println("GormDB : Successfully Connected to Database")

	if conf.DebugLog {
		provider = provider.Debug()
	}

	return provider, sqlDB
}
