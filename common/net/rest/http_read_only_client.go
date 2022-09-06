package rest

import (
	"errors"
	"gitlab.altex.ro/ams/go_ddwcs/common/net"
	"gitlab.altex.ro/plug/go_atx_lib/logger"
	httpLocal "gitlab.altex.ro/plug/go_http"
	"gitlab.altex.ro/plug/go_http/builder"
	"net/http"
	"net/url"
)

type HttpReadOnlyClient struct {
	net.HttpClientBaseSetup
	logger.LoggerAware
}

/**
 * Get
 * @param path string
 * @param queryParams map[string]string
 * @return []byte, error
 * @desc - make a GET request to a third party, but not before builds the entire url adding query string params
 */
func (readOnlyClient *HttpReadOnlyClient) Get(path string, queryParams map[string]string) ([]byte, error) {
	baseUrl := readOnlyClient.GetUrl()
	if readOnlyClient.GetUrl() == nil {
		readOnlyClient.GetLogger().Error("[HttpReadOnlyClient]: you must provide a valid URL")
		return nil, errors.New("you must provide a valid URL")
	}

	baseUrl.Path = path
	if len(queryParams) > 0 {
		params := url.Values{}
		for key, value := range queryParams {
			params.Add(key, value)
		}

		baseUrl.RawQuery = params.Encode()
	}

	request := builder.NewRequestBuilder().Url(baseUrl.String()).Method(http.MethodGet).Headers(*readOnlyClient.GetHeaders()).Build()
	client, err := httpLocal.NewClientRequest(request)

	if err != nil {
		readOnlyClient.GetLogger().Error("[HttpReadOnlyClient] [Build request]: " + err.Error())
		return nil, err
	}

	return client.Call()
}
