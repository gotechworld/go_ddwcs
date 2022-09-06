package strategy

import (
	"gitlab.altex.ro/ams/go_ddwcs/awb/data_store"
	"gitlab.altex.ro/ams/go_ddwcs/awb/dto"
	"gitlab.altex.ro/ams/go_ddwcs/awb/net/soap"
	"gitlab.altex.ro/ams/go_ddwcs/awb/repository"
	"gitlab.altex.ro/ams/go_ddwcs/awb/service/adapter"
	"gitlab.altex.ro/ams/go_ddwcs/awb/service/mapper"
	logger "gitlab.altex.ro/plug/go_logger"
)

type AldAwbStrategy struct {
	awb           *dto.AwbInfo
	client        *soap.AldSoapClient
	logger        logger.Logger
	tokenProvider *repository.AldTokenProvider
	adapter       adapter.Adapter
}

/**
 * GetNewAldAwbStrategy
 * @param awb dto.AwbInfo pointer
 * @param logger log.Logger pointer
 * @return AldAwbStrategy pointer
 * @desc - creates an AldAwbStrategy, initializing its dependencies
 */
func GetNewAldAwbStrategy(awb *dto.AwbInfo, logger logger.Logger) AwbInfoStrategy {
	strategy := new(AldAwbStrategy)
	strategy.logger = logger
	strategy.awb = awb
	strategy.client = soap.InitAldSoapClient(logger, &mapper.AldResponseMapper{})
	strategy.tokenProvider = repository.NewAldTokenProvider(data_store.CreateCacheStoreInstance(), strategy.client)
	strategy.adapter = new(adapter.AldAwbInfoAdapter)

	return strategy
}

/**
 * GetAwbInfo
 * @return adapter.AwbStatusDetails pointer | nil
 * @desc - reads the awb info from external service and returns an uniform type of data
 */
func (strategy *AldAwbStrategy) GetAwbInfo() *adapter.AwbStatusDetails {
	if strategy.awb == nil || strategy.tokenProvider == nil || strategy.client == nil { // ensure that the constructor was used
		strategy.logger.Error("[UrgentCargusAwbStrategy]: missing dependencies, you should use construct method")
		return nil
	}

	status := strategy.adapter.InitAwbStatusDetails(strategy.awb)
	token := strategy.tokenProvider.GetToken()
	if token == "" {
		strategy.logger.Println("[AldAwbStrategy]: Could not login")
		return status
	}

	response := strategy.client.GetAwbInformation(strategy.awb.AwbNumber, token)
	if response == nil ||
		response.Body.XMLName.Space == "" ||
		0 == len(response.Body.GetAwbInformationsResponse.GetAwbInformationsResult.Informations) {
		strategy.logger.Println("[AldAwbStrategy]: could not find any status for " + strategy.awb.AwbNumber)
		return status
	}

	return strategy.adapter.Transform(response.Body.GetAwbInformationsResponse.GetAwbInformationsResult.Informations, status)
}
