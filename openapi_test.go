package openapi

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	api, err := NewDefaultOpenApi("config_debug", "yaml", ".")
	if err != nil {
		fmt.Println(err)
	}
	api.SetDefaultWebHandler()
	err = api.RunOpenApi()
	if err != nil {
		fmt.Println(err)
	}
}

func TestNewDefaultOpenApiService(t *testing.T) {
	//api, err := NewDefaultOpenApi("config_debug", "yaml", ".")
}
