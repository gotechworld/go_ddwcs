package global_config

import (
	"gitlab.altex.ro/ams/go_ddwcs/common/dto/global_config"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/net/rest"
	"gitlab.altex.ro/plug/go_atx_lib/logger"
	"encoding/json"
)

const GC_API_CONFIG_ENDPOINT = "/config"

type ExternalConfigStore struct {
	globalConfigClient *rest.GlobalConfigClient
	logger.LoggerAware
}

// ExternalConfigStore - constructor
func NewExternalConfigStore(globalConfigClient *rest.GlobalConfigClient) *ExternalConfigStore {
	return &ExternalConfigStore{globalConfigClient: globalConfigClient}
}

// Get Configuration List
// @param params map[string]string
// @return global_config.ConfigGet
// @desc - get all Configuration from Global Config API
//
func (ecs *ExternalConfigStore) GetConfig(params map[string]string) global_config.ConfigGet {
	var config global_config.ConfigGet
	if ecs.globalConfigClient == nil {
		ecs.GetLogger().Printf("[ERROR] [GlobalConfigAPI]: 'globalConfigClient' is undefined")
		return config
	}

	responseBytes, err := ecs.globalConfigClient.Get(GC_API_CONFIG_ENDPOINT, params)
	if err != nil {
		ecs.GetLogger().Printf(
			"[ERROR] [GlobalConfigAPI] [%s]: %+v",
			ecs.globalConfigClient.GetUrl(),
			err,
		)
		return config
	}

	var configResponse global_config.ConfigGet
	err = json.Unmarshal(responseBytes, &configResponse)
	if err != nil {
		ecs.GetLogger().Printf(
			"[ERROR] [GlobalConfigAPI] [%s] on Unmarshal: %s. Response unmarshalled: %s",
			ecs.globalConfigClient.GetUrl(),
			err,
			responseBytes,
		)
		return config
	}

	ecs.GetLogger().Printf(
		"[GlobalConfigAPI] [%s]",
		ecs.globalConfigClient.GetUrl(),
	)

	return configResponse
}
