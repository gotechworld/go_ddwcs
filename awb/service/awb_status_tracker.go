package service

import (
	"gitlab.altex.ro/ams/go_ddwcs/awb/dto"
	"gitlab.altex.ro/ams/go_ddwcs/awb/service/adapter"
	"gitlab.altex.ro/ams/go_ddwcs/awb/service/strategy"
	logger "gitlab.altex.ro/plug/go_logger"
	"strings"
)

const (
	CARGUS = "cargus"
	ALD    = "altex curier"
)

type AwbStatusTracker struct {
	Logger logger.Logger
}

/**
 * GetAwbStatus
 * @param awb dto.AwbInfo
 * @return string
 * @desc - reads the awb info from external service
 */
func (ast *AwbStatusTracker) GetAwbStatus(awb *dto.AwbInfo) *adapter.AwbStatusDetails {
	var queryStrategy strategy.AwbInfoStrategy

	switch strings.ToLower(awb.AwbCarrier) {
	case CARGUS:
		queryStrategy = strategy.GetNewUrgentCargusAwbStrategy(awb, ast.Logger)
	case ALD:
		queryStrategy = strategy.GetNewAldAwbStrategy(awb, ast.Logger)
	}

	if queryStrategy == nil {
		ast.Logger.Printf("[AwbStatusTracker]: no strategy found for [%s]", awb.AwbCarrier)
		return nil
	}

	return queryStrategy.GetAwbInfo()
}
