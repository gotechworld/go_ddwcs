package controller

import (
	"encoding/json"
	"gitlab.altex.ro/ams/go_ddwcs/common"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/dto"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/dto/validator"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/service"
	lib "gitlab.altex.ro/plug/go_atx_lib"
	httpLocal "gitlab.altex.ro/plug/go_http"
	"io/ioutil"
	"net/http"
)

type EstimateController struct{}

func init() {
	lib.GetContainer().GetControllerRegistry().Register(NewEstimateController())
}

// EstimateController - constructor
func NewEstimateController() *EstimateController {
	return &EstimateController{}
}

// EstimateAction GET /v1/delivery/estimate/
// @request params "sku[]", "region_code" and "website_code"
// @return httpLocal.Response (atx_lib http.Response)
// @desc   Controller action used to estimate delivery time for a product list
// @deprecated - use POST /v1/delivery/estimate/
func (ac *EstimateController) EstimateAction(r *http.Request) httpLocal.Response {
	queryString := r.URL.Query()

	// get multiple values
	sku := queryString["sku[]"]
	regionCode := queryString.Get("region_code")
	websiteCode := queryString.Get("website_code")
	if !common.IsValidWebsiteCode(websiteCode) {
		websiteCode = common.ALTEX_WEBSITE_CODE
	}

	deliveryEstimator := service.NewDeliveryEstimator(lib.GetContainer().GetLogger(), websiteCode)

	statusCode := http.StatusOK
	response := new(dto.EstimateResponseGet)
	response.EstimatedDeliveryDays, statusCode = deliveryEstimator.EstimateDelivery(sku, regionCode)

	return httpLocal.NewResponse(response, statusCode)
}

// EstimateAction POST /v1/delivery/estimate/
// @request body dto.Request
// @return httpLocal.Response (atx_lib http.Response)
// @desc   Controller action used to estimate delivery time for a product list
func (ac *EstimateController) EstimatePostAction(r *http.Request) httpLocal.Response {

	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return httpLocal.NewErrorResponse("Read body error", http.StatusBadRequest)
	}

	// Unmarshal
	var requestData dto.EstimateRequestPost
	err = json.Unmarshal(b, &requestData)
	if err != nil {
		return httpLocal.NewErrorResponse("Invalid/unaccepted JSON", http.StatusBadRequest)
	}

	isValid, message := validator.NewEstimateRequestValidator().Validate(&requestData)
	if !isValid {
		return httpLocal.NewErrorResponse(message, http.StatusBadRequest)
	}

	deliveryEstimator := service.NewDeliveryPostEstimator(lib.GetContainer().GetLogger(), &requestData)
	response, statusCode := deliveryEstimator.EstimateDelivery()

	return httpLocal.NewResponse(response, statusCode)
}
