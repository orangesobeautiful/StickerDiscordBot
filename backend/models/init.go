package models

import (
	"database/sql"
	stdLog "log"
	"os"
	"time"

	"backend/config"
	"backend/pkg/log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db    *gorm.DB
	sqlDB *sql.DB
)

func Init(cfg *config.CfgInfo) error {
	var err error

	logLevel := logger.Silent
	if cfg.Debug {
		logLevel = logger.Info
	}

	gormLogger := logger.New(
		stdLog.New(os.Stdout, "\r\n", stdLog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logLevel,    // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)

	db, err = gorm.Open(mysql.Open(cfg.Database.DSN),
		&gorm.Config{
			Logger: gormLogger,
		})
	if err != nil {
		log.Errorf("gorm.Open failed, err=%s", err)
		return err
	}

	sqlDB, err = db.DB()
	if err != nil {
		log.Errorf("db.DB() failed, err=%s", err)
		return err
	}

	err = autoMigrate()
	if err != nil {
		return err
	}

	err = initSetting()
	if err != nil {
		return err
	}

	return nil
}

func Close() {
	_ = sqlDB.Close()
}

func autoMigrate() error {
	autoMigrateList := map[string]any{
		"sticker": Sticker{},
		"botinfo": BotInfo{},
	}

	var err error
	for key, table := range autoMigrateList {
		err = db.AutoMigrate(table)
		if err != nil {
			log.Errorf(key+" AutoMigrate failed, err=", err)
			return err
		}
	}

	return nil
}
