package openapi

import (
	"context"
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/kulisi/openapi/conf"
	"github.com/kulisi/openapi/db"
	"github.com/kulisi/openapi/logger"
	"github.com/kulisi/openapi/util"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"
)

type OpenApi struct {
	_ConfigDriver   *viper.Viper
	_LoggerDriver   *zap.Logger
	_DataBaseDriver *gorm.DB
	OpenApiConfig   *conf.Config

	ConfigChangeFunc func(event fsnotify.Event)

	_Handler *gin.Engine
}

func NewDefaultOpenApi(filename, filetype string, path ...string) (*OpenApi, error) {
	openApi := &OpenApi{_ConfigDriver: viper.New()}
	// 限定配置文件的文件名
	openApi._ConfigDriver.SetConfigName(filename)
	// 限定配置文件的文件扩展名
	openApi._ConfigDriver.SetConfigType(filetype)
	// 添加配置文件检索路径
	for _, s := range path {
		openApi._ConfigDriver.AddConfigPath(s)
	}
	// 读取配置文件
	if err := openApi._ConfigDriver.ReadInConfig(); err != nil {
		return nil, err
	}
	// 反序列化配置文件
	if err := openApi._ConfigDriver.Unmarshal(&openApi.OpenApiConfig); err != nil {
		return nil, err
	}
	/*
		配置日志实例
	*/
	// 根据配置判断是否启用日志
	if openApi.UseLogger() {
		// 转换日志保存目录为绝对路径
		absPath, err := filepath.Abs(openApi.OpenApiConfig.Zap.Director)
		if err != nil {
			return nil, err
		}
		// 创建日志存放目录
		if ok, _ := util.PathExists(absPath); !ok {
			err := os.MkdirAll(absPath, os.ModePerm)
			if err != nil {
				return nil, err
			}
		}
		cores := make([]zapcore.Core, 0, len(openApi.OpenApiConfig.Zap.Levels()))
		for i := 0; i < len(openApi.OpenApiConfig.Zap.Levels()); i++ {
			core := logger.NewZapCore(&logger.ZapCoreConfig{
				Level:        openApi.OpenApiConfig.Zap.Levels()[i],
				Encoder:      openApi.OpenApiConfig.Zap.Encoder(),
				Director:     absPath,
				RetentionDay: openApi.OpenApiConfig.Zap.RetentionDay,
				LogInConsole: openApi.OpenApiConfig.Zap.LogInConsole,
			})
			cores = append(cores, core)
		}
		openApi._LoggerDriver = zap.New(zapcore.NewTee(cores...))
		if openApi.OpenApiConfig.Zap.ShowLine {
			openApi._LoggerDriver.WithOptions(zap.AddCaller())
		}
	}
	/*
		配置数据库实例
	*/
	switch strings.ToLower(openApi.OpenApiConfig.Gorm.Use) {
	case "mssql", "sqlserver":
		sqlserverDb, err := db.GenerateSqlserverDb(openApi.OpenApiConfig.Gorm.Mssql)
		if err != nil {
			return nil, err
		}
		openApi._DataBaseDriver = sqlserverDb
	case "mysql":
		mysqlDb, err := db.GenerateMysqlDb(openApi.OpenApiConfig.Gorm.Mysql)
		if err != nil {
			return nil, err
		}
		openApi._DataBaseDriver = mysqlDb
	default:
		openApi._DataBaseDriver = nil
	}
	return openApi, nil
}

func (openApi *OpenApi) UseLogger() bool {
	return openApi.OpenApiConfig.Zap.Use
}

// RunOpenApi 服务
func (openApi *OpenApi) RunOpenApi() (err error) {
	if !openApi.OpenApiConfig.Gin.Use || openApi.OpenApiConfig.Gin.Addr == "" {
		return errors.New("gin is not configured")
	}
	// 判断 web handler 是否为 nil
	if openApi._Handler == nil {
		return errors.New("openApi handler is nil")
	}
	webService := &http.Server{Addr: fmt.Sprintf(":%s", openApi.OpenApiConfig.Gin.Addr), Handler: openApi._Handler}
	go func() {
		if err = webService.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), openApi.OpenApiConfig.Gin.WaitFor*time.Second)
	defer cancel()
	if err = webService.Shutdown(ctx); err != nil {
		return
	}
	return
}

// DebugLog 输出 Debug 级别的日志
func (openApi *OpenApi) DebugLog(msg string, fields ...zap.Field) {
	if openApi._LoggerDriver != nil {
		openApi._LoggerDriver.Debug(msg, fields...)
	} else {
		log.Println(msg)
	}
}

// InfoLog 输出 Info 级别的日志
func (openApi *OpenApi) InfoLog(msg string, fields ...zap.Field) {
	if openApi._LoggerDriver != nil {
		openApi._LoggerDriver.Info(msg, fields...)
	} else {
		log.Println(msg)
	}
}

// WarnLog 输出 Warn 级别的日志
func (openApi *OpenApi) WarnLog(msg string, fields ...zap.Field) {
	if openApi._LoggerDriver != nil {
		openApi._LoggerDriver.Warn(msg, fields...)
	} else {
		log.Println(msg)
	}
}

// ErrorLog 输出 Error 级别的日志
func (openApi *OpenApi) ErrorLog(msg string, fields ...zap.Field) {
	if openApi._LoggerDriver != nil {
		openApi._LoggerDriver.Error(msg, fields...)
	} else {
		log.Println(msg)
	}
}

// DbAutoMigrate run auto migration for given models
func (openApi *OpenApi) DbAutoMigrate(dst ...interface{}) error {
	if openApi._DataBaseDriver != nil {
		return openApi._DataBaseDriver.AutoMigrate(dst...)
	}
	return errors.New("db is not configured")
}

// Db 获取数据库实例
func (openApi *OpenApi) Db() *gorm.DB {
	return openApi._DataBaseDriver
}

// DbExec 执行数据库操作
func (openApi *OpenApi) DbExec(sql string, values ...interface{}) *gorm.DB {
	return openApi._DataBaseDriver.Exec(sql, values...)
}

// DbRaw 查询数据信息
func (openApi *OpenApi) DbRaw(sql string, values ...interface{}) *gorm.DB {
	return openApi._DataBaseDriver.Raw(sql, values...)
}

// SetWebHandler 设置 Web 控制器
func (openApi *OpenApi) SetWebHandler(Handler *gin.Engine) {
	openApi._Handler = Handler
}

// SetDefaultWebHandler 设置 测试的 Web 控制器
func (openApi *OpenApi) SetDefaultWebHandler() {
	gin.SetMode(gin.ReleaseMode)
	handler := gin.Default()
	handler.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})
	handler.GET("/time", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"time": time.Now().Unix(),
		})
	})
	openApi._Handler = handler
}
