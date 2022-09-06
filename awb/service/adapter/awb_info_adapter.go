package adapter

import "gitlab.altex.ro/ams/go_ddwcs/awb/dto"

//Altex Statuses
const (
	StatusNotFound           = 0
	StatusAwbGenerated       = 1
	StatusPickedUp           = 2
	StatusInCourierWarehouse = 3
	StatusInTransit          = 4
	StatusDelivered          = 5
	StatusUndelivered        = 6
)

// Status to description mapping
var StatusDescription = map[int]string{
	StatusNotFound:           "Informatie inexistenta",
	StatusAwbGenerated:       "Awb Generat",
	StatusPickedUp:           "Produse predate curierului",
	StatusInCourierWarehouse: "Produse in depozitul curierului",
	StatusInTransit:          "Produse in curs de livrare",
	StatusDelivered:          "Produse livrate",
	StatusUndelivered:        "Produse nelivrate",
}

// general Awb details struct
type AwbStatusDetails struct {
	IncrementId string
	Awb         string
	Courier     string
	Status      int
	History     []StatusHistory
}

type StatusHistory struct {
	Status string
	Date   string
}

type Adapter interface {
	InitAwbStatusDetails(*dto.AwbInfo) *AwbStatusDetails
	Transform(interface{}, *AwbStatusDetails) *AwbStatusDetails
}

// Generic adapter
type BaseAdapter struct{}

/**
 * initAwbStatusDetails - private
 * @param awb dto.AwbInfo
 * @return adapter.AwbStatusDetails pointer
 * @desc - basic info of the awb status
 */
func (adapter BaseAdapter) InitAwbStatusDetails(awb *dto.AwbInfo) *AwbStatusDetails {
	status := new(AwbStatusDetails)
	status.IncrementId = awb.OrderIncrementId
	status.Awb = awb.AwbNumber
	status.Courier = awb.AwbCarrier
	status.Status = StatusAwbGenerated

	return status
}
