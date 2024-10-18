package db

import (
	"fmt"
	"github.com/kulisi/openapi/conf"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

// Config
// @author: [kulisi](https://github.com/kulisi)
// @function: Config
func Config(general conf.GeneralDB) *gorm.Config {

	return &gorm.Config{
		SkipDefaultTransaction: false,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   general.Prefix,
			SingularTable: general.Singular,
		},
		Logger: logger.New(NewWriter(general, log.New(os.Stdout, "\r\n", log.LstdFlags)), logger.Config{
			SlowThreshold: 200 * time.Second,
			LogLevel:      general.LogLevel(),
			Colorful:      true,
		}),
		DisableForeignKeyConstraintWhenMigrating: true,
	}

}

type Writer struct {
	config conf.GeneralDB
	writer logger.Writer
}

func NewWriter(config conf.GeneralDB, writer logger.Writer) *Writer {
	return &Writer{config: config, writer: writer}
}

func (w *Writer) Printf(msg string, data ...any) {
	if w.config.LogZap {
		switch w.config.LogLevel() {
		case logger.Silent:
			zap.L().Debug(fmt.Sprintf(msg, data))
		case logger.Info:
			zap.L().Info(fmt.Sprintf(msg, data))
		case logger.Warn:
			zap.L().Warn(fmt.Sprintf(msg, data))
		case logger.Error:
			zap.L().Error(fmt.Sprintf(msg, data))
		default:
			zap.L().Info(fmt.Sprintf(msg, data))
		}
	} else {
		w.writer.Printf(msg, data)
	}
}
