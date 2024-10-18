package db

import (
	"errors"
	"github.com/kulisi/openapi/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GenerateMysqlDb(m conf.Mysql) (*gorm.DB, error) {
	//m := global.ApiConfig.Mysql
	if m.Dbname == "" {
		return nil, errors.New("dbname is empty")
	}
	config := mysql.Config{
		DSN:               m.Dsn(),
		DefaultStringSize: 255,
	}

	if db, err := gorm.Open(mysql.New(config), Config(m.GeneralDB)); err != nil {
		return nil, err
	} else {
		db.InstanceSet("db:table_options", "ENGINE="+m.Engine)
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db, nil
	}
}
