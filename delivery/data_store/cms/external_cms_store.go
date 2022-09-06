package cms

import (
	"encoding/json"
	"gitlab.altex.ro/ams/go_ddwcs/common/dto/cms"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/net/rest"
	"gitlab.altex.ro/plug/go_atx_lib/logger"
)

const cmsStoresPath = "/stores"

type ExternalCmsStore struct {
	cmsClient *rest.CMSClient
	logger.LoggerAware
}

/**
 * ExternalStocksStore - constructor
 */
func NewExternalStocksStore(cmsClient *rest.CMSClient) *ExternalCmsStore {
	return &ExternalCmsStore{cmsClient: cmsClient}
}

/**
 * @param params map[string]string
 * @return cms.StoreList
 *
 * @desc - get the available stores from CMS API
 */
func (ess *ExternalCmsStore) GetStores(params map[string]string) []cms.Store {
	var emptyList []cms.Store

	if nil == ess.cmsClient {
		ess.GetLogger().Error("[CMS_API]: 'cmsClient' is undefined")
		return emptyList
	}

	responseBytes, err := ess.cmsClient.Get(cmsStoresPath, params)
	if err != nil {
		ess.GetLogger().Printf("[ERROR] [CMS_API] [%s] : %+v", ess.cmsClient.GetUrl(), err)
		return emptyList
	}

	var stockResponse cms.StoreList
	err = json.Unmarshal(responseBytes, &stockResponse)
	if err != nil {
		ess.GetLogger().Printf("[ERROR] [CMS_API] Unmarshal: %s. Response unmarshalled: %s", err, responseBytes)
		return emptyList
	}

	return stockResponse.Items
}
