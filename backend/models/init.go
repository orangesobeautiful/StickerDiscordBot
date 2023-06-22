package models

import (
	"database/sql"

	"backend/config"
	"backend/pkg/log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var sqlDB *sql.DB

func Init() error {
	var err error

	var cfg = config.GetCfg()

	db, err = gorm.Open(mysql.Open(cfg.Database.DSN), &gorm.Config{})
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
