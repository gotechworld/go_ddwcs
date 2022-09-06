package net

import "net/url"

type HttpClientBaseSetup struct {
	baseUrl *url.URL
	headers map[string]string
}

/**
 * GetUrl
 * @return url.URL pointer
 */
func (setupClient *HttpClientBaseSetup) GetUrl() *url.URL {
	return setupClient.baseUrl
}

/**
 * SetUrl
 * @param url string
 * @return error
 * @desc - set the url and parse it in order to validate the integrity
 */
func (setupClient *HttpClientBaseSetup) SetUrl(urlStr string) error {
	baseUrl, err := url.Parse(urlStr)
	setupClient.baseUrl = baseUrl

	return err
}

/**
 * SetHeaders
 * @param headers map[string]string
 * @return void
 */
func (setupClient *HttpClientBaseSetup) SetHeaders(headers map[string]string) {
	setupClient.headers = headers
}

/**
 * GetHeaders
 * @return map[string]string
 */
func (setupClient *HttpClientBaseSetup) GetHeaders() *map[string]string {
	return &setupClient.headers
}
