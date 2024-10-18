package db

import (
	"errors"
	"github.com/kulisi/openapi/conf"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func GenerateSqlserverDb(m conf.Mssql) (*gorm.DB, error) {
	//m := ApiConfig.Mssql
	if m.Dbname == "" {
		return nil, errors.New("dbname is empty")
	}
	config := sqlserver.Config{
		DSN:               m.Dsn(),
		DefaultStringSize: 255,
	}

	if db, err := gorm.Open(sqlserver.New(config), Config(m.GeneralDB)); err != nil {
		return nil, err
	} else {
		db.InstanceSet("db:table_options", "ENGINE="+m.Engine)
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db, nil
	}
}
