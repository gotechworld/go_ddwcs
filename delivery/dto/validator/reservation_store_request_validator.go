package validator

import (
	"gitlab.altex.ro/ams/go_ddwcs/common"
	"gitlab.altex.ro/ams/go_ddwcs/common/net/request"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/dto"
	"net/http"
)

const (
	minProductsToEvaluate = 1
	maxProductsToEvaluate = 100
)

type PayloadValidator struct{}

func (v PayloadValidator) Validate(payload interface{}) (bool, error) {
	reservationPayload, ok := payload.(*dto.ReservationStorePayload)
	if !ok || len(reservationPayload.Products) < minProductsToEvaluate {
		return false, request.NewRequestError(http.StatusBadRequest, "Invalid data!")
	}

	if !common.IsValidWebsiteCode(reservationPayload.WebsiteCode) {
		return false, request.NewRequestError(http.StatusBadRequest, "Invalid website code!")
	}

	if len(reservationPayload.Products) >= maxProductsToEvaluate {
		return false, request.NewRequestError(http.StatusBadRequest, "To much products!")
	}

	for _, item := range reservationPayload.Products {
		if item.Sku == "" || item.Qty < 1 {
			return false, request.NewRequestError(http.StatusBadRequest, "Invalid products!")
		}
	}

	return true, nil
}
