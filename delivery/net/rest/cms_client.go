package rest

import (
	"errors"
	"fmt"
	"gitlab.altex.ro/ams/go_ddwcs/common/net/rest"
	lib "gitlab.altex.ro/plug/go_atx_lib"
	httpLocal "gitlab.altex.ro/plug/go_http"
	"gitlab.altex.ro/plug/go_http/builder"
	logger "gitlab.altex.ro/plug/go_logger"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const defaultFilterKey = "filter"

// extends HttpReadOnlyClient
type CMSClient struct {
	rest.HttpReadOnlyClient
}

type CMSResponse struct {
	httpLocal.StandardResponse
}

/**
 * CMSClient - constructor
 * @desc - init a new cms client
 */

func NewCMSClient(websiteCode string, logger logger.Logger) (*CMSClient, error) {
	client := &CMSClient{}
	client.SetLogger(logger)

	baseUrl := lib.GetContainer().GetConfig().SystemVars["cms_host_"+websiteCode]
	err := client.SetUrl(baseUrl)
	if err != nil {
		errMsg := fmt.Sprintf("[ERROR] [CMS_API] invalid base url [%s] for website [%s]", baseUrl, websiteCode)
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

/**
 * @override
 *
 * @param path string
 * @param queryParams map[string]string
 * @return []byte, error
 * @desc - make a GET request to cms-api
 */
func (c *CMSClient) Get(path string, queryParams map[string]string) ([]byte, error) {
	baseUrl := c.GetUrl()
	if c.GetUrl() == nil {
		c.GetLogger().Error("[CMSClient]: you must provide a valid URL")
		return nil, errors.New("you must provide a valid URL")
	}

	baseUrl.Path = path
	if len(queryParams) > 0 {
		params := url.Values{}
		for key, value := range queryParams {
			if strings.Contains(key, defaultFilterKey) {
				params.Add(defaultFilterKey, value)
			} else {
				params.Add(key, value)
			}
		}

		baseUrl.RawQuery = params.Encode()
	}

	request := builder.NewRequestBuilder().Url(baseUrl.String()).Method(http.MethodGet).Headers(*c.GetHeaders()).Build()
	client, err := httpLocal.NewClientRequest(request)

	if err != nil {
		c.GetLogger().Error("[CMSClient] [Build request]: " + err.Error())
		return nil, err
	}

	return client.Call()
}
