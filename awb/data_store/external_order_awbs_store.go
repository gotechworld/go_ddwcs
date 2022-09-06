package data_store

import (
	"encoding/json"
	"gitlab.altex.ro/ams/go_ddwcs/awb/dto"
	"gitlab.altex.ro/ams/go_ddwcs/awb/net/rest"
	"gitlab.altex.ro/plug/go_atx_lib/logger"
)

const PATH_TO_AWB_LIST = "v1.0/awb"

type ExternalOrderAwbsStore struct {
	logger.LoggerAware
	omsClient *rest.OmsClient
}

func GetExternalStore(omsClient *rest.OmsClient) *ExternalOrderAwbsStore {
	return &ExternalOrderAwbsStore{omsClient: omsClient}
}

/**
 * GetList
 * @param params map[string]string
 * @return []dto.AwbInfo
 * @desc - get the awbs list from external service
 * the params map will contain the following keys: increment_id, customer_id, customer_email
 */
func (externalAwbsStore *ExternalOrderAwbsStore) GetList(params map[string]string) []dto.AwbInfo {
	var emptyList []dto.AwbInfo
	if externalAwbsStore.omsClient == nil {
		return emptyList
	}

	response, err := externalAwbsStore.omsClient.Get(PATH_TO_AWB_LIST, params)
	if err != nil {
		externalAwbsStore.GetLogger().Warning(err.Error())
		return emptyList
	}

	var list rest.OmsResponse
	var data *dto.OrderAwbList

	list.Data = &data
	err = json.Unmarshal(response, &list)
	if err != nil {
		externalAwbsStore.GetLogger().Error(err.Error())
		return emptyList
	}

	//todo should cache the result
	return data.Awbs
}
