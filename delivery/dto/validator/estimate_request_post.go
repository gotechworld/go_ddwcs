package validator

import (
	"fmt"
	"gitlab.altex.ro/ams/go_ddwcs/common"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/dto"
)

type EstimateRequestValidator struct{}

// NewEstimateRequestValidator - constructor
func NewEstimateRequestValidator() *EstimateRequestValidator {
	return &EstimateRequestValidator{}
}

// Set WebsiteCode with default value and Validate "products" and "attribute_set_id" keys
func (e EstimateRequestValidator) Validate(request *dto.EstimateRequestPost) (isValid bool, message string) {
	// tried with github.com/go-playground/validator/v10 v10.3.0 but for some reason attribute_set_id fails
	// @todo
	//validate := validator.New()
	//err := validate.Struct(request)
	//if err != nil {
	//	return true, err.Error()
	//}

	if !common.IsValidWebsiteCode(request.WebsiteCode) {
		request.WebsiteCode = common.ALTEX_WEBSITE_CODE
	}

	if nil == request.Products {
		return false, "'products' key is missing"
	}

	for sku, productInfo := range request.Products {
		if 0 == productInfo.AttributeSetId {
			return true, fmt.Sprintf("'attribute_set_id' key is missing for %s", sku)
		}
	}

	return true, ""
}
