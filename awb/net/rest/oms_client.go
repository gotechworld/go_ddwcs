package rest

import (
	"errors"
	"gitlab.altex.ro/ams/go_ddwcs/common/net/rest"
	lib "gitlab.altex.ro/plug/go_atx_lib"
	httpLocal "gitlab.altex.ro/plug/go_http"
	logger "gitlab.altex.ro/plug/go_logger"
	"os"
)

/**
 * extends HttpReadOnlyClient
 */
type OmsClient struct {
	rest.HttpReadOnlyClient
}

type OmsResponse struct {
	httpLocal.StandardResponse
	Messages []string
}

/**
 * InitOmsClient
 * @desc - init a new oms client by validating the provided host and populating the credentials for a future request
 */
func InitOmsClient(logger logger.Logger) (*OmsClient, error) {
	client := &OmsClient{}
	client.SetLogger(logger)

	err := client.SetUrl(lib.GetContainer().GetConfig().SystemVars["oms_host"])
	if err != nil {
		client.GetLogger().Error("[OmsClient]: " + err.Error())
		return nil, errors.New("invalid oms base url")
	}

	client.SetHeaders(map[string]string{
		"Accept":         "application/json",
		"X-Service-Name": os.Getenv("DDWCS_USER"),
		"X-Service-Key":  os.Getenv("DDWCS_KEY"),
	})

	return client, nil
}
