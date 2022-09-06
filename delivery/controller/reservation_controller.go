package controller

import (
	"gitlab.altex.ro/ams/go_ddwcs/common/net/request"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/dto"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/dto/validator"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/service"
	lib "gitlab.altex.ro/plug/go_atx_lib"
	httpLocal "gitlab.altex.ro/plug/go_http"
	"net/http"
)

type ReservationController struct{}

func init() {
	lib.GetContainer().GetControllerRegistry().Register(NewReservationController())
}

// ReservationController - constructor
func NewReservationController() *ReservationController {
	return &ReservationController{}
}

/**
 * ReservationStoreAction
 * @param r *http.Request
 * @return httpLocal.Response (atx_lib http.Response)
 * @desc - Controller action used check available stores for a certain shopping cart configuration (product-quantity)
 */
func (ac *ReservationController) ReservationStoreAction(r *http.Request) httpLocal.Response {
	var requestData dto.ReservationStorePayload
	logger := lib.GetContainer().GetLogger()
	requestHandler := request.NewRequestHandler(&request.RequestNormalizer{}, &validator.PayloadValidator{}, logger)

	if err := requestHandler.Normalize(r.Body, &requestData).Validate().Result(); err != nil {
		return httpLocal.NewErrorResponse(err.Error(), err.GetStatus())
	}

	stores := service.NewReservationStoreService(requestData.WebsiteCode, logger).GetAvailableStores(&requestData)

	return httpLocal.NewResponse(stores, http.StatusOK)
}
