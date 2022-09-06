package rest

import (
	"gitlab.altex.ro/ams/go_ddwcs/common/net/rest"
	lib "gitlab.altex.ro/plug/go_atx_lib"
	httpLocal "gitlab.altex.ro/plug/go_http"
	logger "gitlab.altex.ro/plug/go_logger"
	"errors"
	"fmt"
	"os"
)

// Global Config Client
type GlobalConfigClient struct {
	rest.HttpReadOnlyClient
}

type GlobalConfigResponse struct {
	httpLocal.StandardResponse
}

// GlobalConfigClient - constructor
// @desc - init a new GC client by validating the provided host and populating the credentials for a future request
func NewGlobalConfigClient(logger logger.Logger) (*GlobalConfigClient, error) {
	client := &GlobalConfigClient{}
	client.SetLogger(logger)

	baseUrl := lib.GetContainer().GetConfig().SystemVars["global_config_host"]
	err := client.SetUrl(baseUrl)
	if err != nil {
		errMsg := fmt.Sprintf("[ERROR] [GlobalConfigAPI] invalid base url [%s]", baseUrl)
		client.GetLogger().Error(errMsg)
		return nil, errors.New(errMsg)
	}

	client.SetHeaders(map[string]string{
		"Accept":         "application/json",
		"X-Service-Name": os.Getenv("DDWCS_USER"),
		"X-Service-Key":  os.Getenv("DDWCS_KEY"),
	})

	return client, nil
}

