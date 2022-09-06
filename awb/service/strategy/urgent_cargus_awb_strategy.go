package strategy

import (
	"gitlab.altex.ro/ams/go_ddwcs/awb/dto"
	"gitlab.altex.ro/ams/go_ddwcs/awb/net/rest"
	"gitlab.altex.ro/ams/go_ddwcs/awb/service/adapter"
	"gitlab.altex.ro/ams/go_ddwcs/awb/service/mapper"
	logger "gitlab.altex.ro/plug/go_logger"
)

type UrgentCargusAwbStrategy struct {
	awb     *dto.AwbInfo
	client  *rest.UrgentCargusClient
	logger  logger.Logger
	adapter adapter.Adapter
}

/**
 * GetNewUrgentCargusAwbStrategy - Construct
 * @param awb dto.AwbInfo pointer
 * @param logger log.Logger pointer
 * @return UrgentCargusAwbStrategy pointer
 * @desc - creates an urgentCargusAwbStrategy, initializing its dependencies
 */
func GetNewUrgentCargusAwbStrategy(awb *dto.AwbInfo, logger logger.Logger) AwbInfoStrategy {
	strategy := new(UrgentCargusAwbStrategy)
	strategy.logger = logger
	strategy.awb = awb

	cargus, err := rest.InitUrgentCargusClient(&mapper.CargusResponseMapper{}, logger)
	if err != nil {
		logger.Error("[UrgentCargusAwbStrategy]: " + err.Error())
		return strategy
	}

	strategy.client = cargus
	strategy.adapter = new(adapter.UrgentCargusAwbInfoAdapter)

	return strategy
}

/**
 * GetAwbInfo
 * @return adapter.AwbStatusDetails pointer | nil
 * @desc - reads the awb info from external service and returns an uniform type of data
 */
func (strategy *UrgentCargusAwbStrategy) GetAwbInfo() *adapter.AwbStatusDetails {
	if strategy.awb == nil || strategy.adapter == nil {
		strategy.logger.Error("[UrgentCargusAwbStrategy]: missing dependencies, you should use construct method")
		return nil
	}

	status := strategy.adapter.InitAwbStatusDetails(strategy.awb)
	if strategy.client == nil {
		strategy.logger.Error("[UrgentCargusAwbStrategy]: missing client")
		return status
	}

	response := strategy.client.GetAwbInformations(strategy.awb.AwbNumber)
	if response == nil {
		strategy.logger.Println("[UrgentCargusAwbStrategy]: could not find any status for " + strategy.awb.AwbNumber)
		return status
	}

	return strategy.adapter.Transform(response, status)
}
