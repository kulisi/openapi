package openapi

import (
	"context"
	"errors"
	"fmt"
	"github.com/kardianos/service"
	"log"
	"net/http"
	"time"
)

// WebServer 实现接口
type WebServer struct {
	cfg *service.Config
	web *http.Server

	startFunc func(server *http.Server)
	stopFunc  func(server *http.Server) error
}

// Start 实现接口函数
func (w WebServer) Start(svc service.Service) (err error) {
	go w.startFunc(w.web)
	return nil
}

// Stop 实现接口函数
func (w WebServer) Stop(svc service.Service) error {
	return w.stopFunc(w.web)
}

func NewDefaultOpenApiServiceByOpenApi(api *OpenApi) (service.Service, error) {
	if api._Handler == nil {
		return nil, errors.New("openapi handler is nil")
	}
	webServiceConfig := &WebServer{
		cfg: &service.Config{
			Name:        api.OpenApiConfig.Service.Name,
			DisplayName: api.OpenApiConfig.Service.DisplayName,
			Description: api.OpenApiConfig.Service.Description,
		},
		startFunc: func(server *http.Server) {
			server = &http.Server{Addr: fmt.Sprintf(":%s", api.OpenApiConfig.Gin.Addr), Handler: api._Handler}
			if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		},
		stopFunc: func(server *http.Server) error {
			ctx, cancel := context.WithTimeout(context.Background(), api.OpenApiConfig.Gin.WaitFor*time.Second)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				log.Println("shutdown error: ", err)
			}
			select {
			case <-ctx.Done():
				log.Println("timeout of 5 second")
			}
			log.Println("server exiting")
			return nil
		},
	}
	svc, err := service.New(webServiceConfig, webServiceConfig.cfg)
	if err != nil {
		return nil, err
	}
	return svc, nil
}
