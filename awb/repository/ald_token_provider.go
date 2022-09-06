package repository

import (
	"gitlab.altex.ro/ams/go_ddwcs/awb/net/soap"
	"gitlab.altex.ro/ams/go_ddwcs/common/data_store"
	"time"
)

const cachedAldTokenKey = "ald_token_key"
const cachedAldTokenTtl = time.Minute * 15 // 15 minutes of cache

type AldTokenProvider struct {
	cache     *data_store.CacheStore
	aldClient *soap.AldSoapClient
}

/**
 * NewAldTokenProvider - Constructor
 * @desc - Create a new AldTokenProvider instance
 */
func NewAldTokenProvider(cache *data_store.CacheStore, aldClient *soap.AldSoapClient) *AldTokenProvider {
	return &AldTokenProvider{
		cache:     cache,
		aldClient: aldClient,
	}
}

/**
 * GetToken
 * @return string
 * @desc - try to get the token from cache, if it doesn't exists then call AldSoapClient Login and save the result in cache for 15 minutes
 */
func (provider AldTokenProvider) GetToken() string {
	if provider.aldClient == nil && provider.cache == nil {
		return ""
	}

	token := provider.cache.Get(cachedAldTokenKey)
	if token != "" {
		return token
	}

	token = provider.aldClient.Login()
	provider.cache.Save(data_store.Param{Key: cachedAldTokenKey, Value: token, Ttl: cachedAldTokenTtl})

	return token
}
