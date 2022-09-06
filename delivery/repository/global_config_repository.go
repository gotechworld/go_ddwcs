package repository

import (
	"gitlab.altex.ro/ams/go_ddwcs/common"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/data_store"
	"strconv"
	"strings"
)

type GlobalConfigRepository struct {
}

// GlobalConfigRepository - constructor
func NewGlobalConfigRepository() *GlobalConfigRepository {
	return &GlobalConfigRepository{}
}

// Get National Holidays from Global config "ddwcs/national_holidays" key
// @return *map[string]string - {"YYYY-MM-DD":"YYYY-MM-DD", ...}
//
func (gcr GlobalConfigRepository) GetNationalHolidays(websiteCode string) *map[string]string {
	params := map[string]string{
		"code": "ddwcs/national_holidays",
		"platform": "DDWCS_API",
		"scope": strconv.Itoa(common.GetWebsiteId(websiteCode)),
	}
	data := data_store.CreateDecoratedGlobalConfigExternalStore().GetConfig(params)
	nationalHolidays := map[string]string{}
	if len(data.ConfigGetData) != 0 {
		gcNationalHolidays := strings.Split(data.ConfigGetData[0].Value, ",")
		for i := range gcNationalHolidays {
			nationalHolidays[strings.TrimSpace(gcNationalHolidays[i])] = strings.TrimSpace(gcNationalHolidays[i])
		}
	}

	return &nationalHolidays
}
