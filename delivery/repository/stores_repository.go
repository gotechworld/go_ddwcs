package repository

import (
	commonCms "gitlab.altex.ro/ams/go_ddwcs/common/dto/cms"
	"gitlab.altex.ro/ams/go_ddwcs/common/util"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/data_store/cms"
	"time"
)

const cachedCmsStoresKey = "cms_stores:"
const cachedCmsStoresTtl = time.Minute * 15 // 15 minutes of cache

type StoreRepository struct {
	cache     *cms.CacheCmsStore
	cmsClient *cms.ExternalCmsStore
	website   string
}

/**
 * NewAldTokenProvider - Constructor
 * @desc - Create a new AldTokenProvider instance
 */
func NewStoreRepository(cache *cms.CacheCmsStore, cmsClient *cms.ExternalCmsStore, website string) *StoreRepository {
	return &StoreRepository{
		cache:     cache,
		cmsClient: cmsClient,
		website:   website,
	}
}

/**
 * GetStoresBy
 * @return string
 * @desc - try to get the token from cache, if it doesn't exists then call AldSoapClient Login and save the result in cache for 15 minutes
 */
func (sr *StoreRepository) GetStoresBy(filters map[string]string) []commonCms.Store {
	if sr.cmsClient == nil && sr.cache == nil || sr.website == "" {
		return []commonCms.Store{}
	}

	cacheKey := sr.composeCacheKey(filters)
	if cachedStores := sr.cache.Get(cacheKey); len(cachedStores) > 0 {
		return cachedStores
	}

	stores := sr.cmsClient.GetStores(filters)
	sr.cache.Save(cacheKey, stores, cachedCmsStoresTtl)

	return stores
}

/**
 * @desc - compose cache key using filters applied for a possible http request
 *
 * @param filters map[string]string
 * @return string
 */
func (sr *StoreRepository) composeCacheKey(filters map[string]string) string {
	cacheKey := cachedCmsStoresKey + sr.website
	if len(filters) > 0 {
		cacheKey += ":" + util.Implode(filters)
	}

	return cacheKey
}
