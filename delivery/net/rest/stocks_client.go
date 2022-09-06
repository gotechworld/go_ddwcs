package rest

import (
	"errors"
	"fmt"
	"gitlab.altex.ro/ams/go_ddwcs/common/net/rest"
	lib "gitlab.altex.ro/plug/go_atx_lib"
	httpLocal "gitlab.altex.ro/plug/go_http"
	logger "gitlab.altex.ro/plug/go_logger"
	"os"
)

// extends HttpReadOnlyClient
type StocksClient struct {
	rest.HttpReadOnlyClient
}

type StocksResponse struct {
	httpLocal.StandardResponse
}

// StocksClient - constructor
// @desc - init a new stocks client by validating the provided host and populating the credentials for a future request
func NewStocksClient(websiteCode string, logger logger.Logger) (*StocksClient, error) {
	client := &StocksClient{}
	client.SetLogger(logger)

	baseUrl := lib.GetContainer().GetConfig().SystemVars["stocks_host_"+websiteCode]
	err := client.SetUrl(baseUrl)
	if err != nil {
		errMsg := fmt.Sprintf("[ERROR] [StocksAPI] invalid base url [%s] for website [%s]", baseUrl, websiteCode)
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
