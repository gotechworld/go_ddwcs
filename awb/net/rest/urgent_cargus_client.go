package rest

import (
	"errors"
	"gitlab.altex.ro/ams/go_ddwcs/awb/service/mapper"
	"gitlab.altex.ro/ams/go_ddwcs/common/net/rest"
	lib "gitlab.altex.ro/plug/go_atx_lib"
	logger "gitlab.altex.ro/plug/go_logger"
	"os"
)

type UrgentCargusClient struct {
	rest.HttpReadOnlyClient
	responseMapper mapper.Mapper
}

// path to get the awb information
const cargusAwbTracePath = "api/NoAuth/GetAwbTrace"

/**
 * InitAldSoapClient
 * @desc - init a new client by validating the provided host and populating the credentials for a future request
 */
func InitUrgentCargusClient(responseMapper mapper.Mapper, logger logger.Logger) (*UrgentCargusClient, error) {
	client := &UrgentCargusClient{responseMapper: responseMapper}
	client.SetLogger(logger)

	err := client.SetUrl(lib.GetContainer().GetConfig().SystemVars["urgent_cargus_host"])
	if err != nil {
		client.GetLogger().Error("[UrgentCargusClient]: " + err.Error())
		return nil, errors.New("invalid urgent base url")
	}

	client.SetHeaders(map[string]string{
		"Accept":                    "application/json",
		"Ocp-Apim-Subscription-Key": os.Getenv("CARGUS_SUBSCRIPTION_KEY"),
	})

	return client, nil
}

/**
 * GetAwbInformation
 * @param awbCode String
 * @return []byte
 * @desc - reads the awb info from external service
 */
func (urgentClient *UrgentCargusClient) GetAwbInformations(awbCode string) *mapper.CargusAwbInformationResponse {
	bytes, err := urgentClient.Get(cargusAwbTracePath, map[string]string{"barCode": awbCode})
	if err != nil {
		urgentClient.GetLogger().Info("[UrgentCargusClient]: " + err.Error())
		return nil
	}

	if urgentClient.responseMapper == nil {
		urgentClient.GetLogger().Error("[UrgentCargusClient]: You need to use struct instance method")
		return nil
	}

	mapped, ok := urgentClient.responseMapper.MapAwbInfo(bytes, awbCode).(mapper.CargusAwbInformationResponse)
	if !ok {
		urgentClient.GetLogger().Error("[UrgentCargusClient]: could not convert response to [CargusAwbInformationResponse]")
		return nil
	}

	return &mapped
}
