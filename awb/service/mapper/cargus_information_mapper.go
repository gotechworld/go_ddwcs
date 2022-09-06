package mapper

import (
	"encoding/json"
)

type CargusAwbInformationResponse struct {
	Code             string  `json:"Code"`
	Type             string  `json:"Type"`
	MeasuredWeight   float32 `json:"MeasuredWeight"`
	VolumetricWeight float32 `json:"VolumetricWeight"`
	ConfirmationName string  `json:"ConfirmationName"`
	Observation      string  `json:"Observation"`
	ResponseCode     string  `json:"ResponseCode"`
	Event            []Event `json:"Event"`
}

type Event struct {
	Date         string `json:"Date"`
	EventID      int    `json:"EventId"`
	Description  string `json:"Description"`
	LocalityName string `json:"LocalityName"`
}

type CargusResponseMapper struct{}

/**
 * MapAwbInfo
 * @param data []byte
 * @return interface|CargusAwbInformationResponse
 * @desc - maps the awb info from a provided json
 */
func (mapper *CargusResponseMapper) MapAwbInfo(data []byte, awbCode string) interface{} {
	var responseList []CargusAwbInformationResponse
	err := json.Unmarshal(data, &responseList)
	if err != nil || len(responseList) == 0 {
		return CargusAwbInformationResponse{}
	}

	//filter the responseList
	for _, response := range responseList {
		if response.Code == awbCode {
			return response
		}
	}

	return responseList[0]
}

/**
 * MapLoginToken
 * @param data []byte
 * @return string
 * @desc - maps the login result
 */
func (mapper *CargusResponseMapper) MapLoginToken(data []byte) string {
	return string(data)
}
