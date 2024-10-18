package openapi

import (
	"context"
	"errors"
	"fmt"
	"github.com/kardianos/service"
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

func NewDefaultOpenApiService(handler http.Handler) (*OpenApi, service.Service, error) {
	api, err := NewDefaultOpenApi("config", "yaml", ".")
	if err != nil {
		return nil, nil, err
	}
	webServiceConfig := &WebServer{
		cfg: &service.Config{
			Name:        api.OpenApiConfig.Service.Name,
			DisplayName: api.OpenApiConfig.Service.DisplayName,
			Description: api.OpenApiConfig.Service.Description,
		},
		startFunc: func(server *http.Server) {
			server = &http.Server{Addr: fmt.Sprintf(":%s", api.OpenApiConfig.Gin.Addr), Handler: handler}
			if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		},
		stopFunc: func(server *http.Server) error {
			ctx, cancel := context.WithTimeout(context.Background(), api.OpenApiConfig.Gin.WaitFor*time.Second)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				return err
			}
			select {
			case <-ctx.Done():
			}
			return nil
		},
	}
	svc, err := service.New(webServiceConfig, webServiceConfig.cfg)
	if err != nil {
		return nil, nil, err
	}
	return api, svc, nil
}
