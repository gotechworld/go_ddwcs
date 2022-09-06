package adapter

import (
	"gitlab.altex.ro/ams/go_ddwcs/awb/service/mapper"
	"gitlab.altex.ro/ams/go_ddwcs/awb/util"
	utilLib "gitlab.altex.ro/plug/go_atx_lib/util"
)

/**
 * ALD possible statuses
 *
 * Incarcat la curier
 * Incarcare finalizata
 * Livrat
 * Nelivrat
 * Retur in depozit
 * Buy back
 * In depozit
 * Neplanificat
 * Receptionat in depozit
 * Sortare
 * Reinitializare
 * Manifest nepreluat
 * Reprogramat
 * Receptionat retur
 * Returnat la expeditor
 * ReincarcareAwb
 * Anulat
 * Transfer
 * Constituire cursa
 * Receptie depozit destinatie
 * Colectat de curier
 * Inventar
 * Stationare depozit
 * Rutare
 * Stationare HUB
 */

var ignoredEvents = []string{
	"Reinitializare",
	"Buy back",
	"Rutare",
	"Constituire cursa",
}

var aldEventMap = map[string]int{
	"Manifest nepreluat": StatusAwbGenerated,
	"ReincarcareAwb":     StatusAwbGenerated,
	"Neplanificat":       StatusAwbGenerated,

	"Transfer":           StatusPickedUp,
	"Colectat de curier": StatusPickedUp,
	"Receptionat retur":  StatusPickedUp,

	"Receptionat in depozit":      StatusInCourierWarehouse,
	"Receptie depozit destinatie": StatusInCourierWarehouse,
	"In depozit":                  StatusInCourierWarehouse,
	"Stationare depozit":          StatusInCourierWarehouse,
	"Stationare HUB":              StatusInCourierWarehouse,
	"Retur in depozit":            StatusInCourierWarehouse,
	"Sortare":                     StatusInCourierWarehouse,
	"Inventar":                    StatusInCourierWarehouse,

	"Incarcat la curier":   StatusInTransit,
	"Incarcare finalizata": StatusInTransit,
	"Reprogramat":          StatusInTransit,
	"Livrat":               StatusDelivered,

	"Nelivrat":              StatusUndelivered,
	"Anulat":                StatusUndelivered,
	"Returnat la expeditor": StatusUndelivered,
}

type AldAwbInfoAdapter struct {
	BaseAdapter
}

/**
 * Transform
 * @param infoData interface{}|[]mapper.Informations
 * @param awbStatusDetails AwbStatusDetails pointer
 * @return AwbStatusDetails pointer
 * @desc - transforms the awb info from ald into a common structure
 */
func (adapter *AldAwbInfoAdapter) Transform(infoData interface{}, awbStatusDetails *AwbStatusDetails) *AwbStatusDetails {
	info := infoData.([]mapper.Informations)
	if len(info) == 0 {
		return awbStatusDetails
	}

	for _, history := range info {
		// ignore some events
		if utilLib.InArray(ignoredEvents, history.StatusAwb) {
			continue
		}

		if aldEventMap[history.StatusAwb] > 0 {
			awbStatusDetails.Status = aldEventMap[history.StatusAwb]
		}

		awbStatusDetails.History = append(awbStatusDetails.History, StatusHistory{
			Status: history.StatusAwb,
			Date:   util.StringToDateTimeFormat(history.DataStatus),
		})
	}

	return awbStatusDetails
}
