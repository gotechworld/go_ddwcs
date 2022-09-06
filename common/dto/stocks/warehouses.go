package stocks

import "gitlab.altex.ro/ams/go_ddwcs/common/dto/ecom_base"

// Possible statuses for Warehouse primary field
const WH_IS_PRIMARY = 1
const WH_IS_NOT_PRIMARY = 0

// Possible statuses for Warehouse hasDisplayProducts
const WH_HAS_DISPLAY_PRODUCTS = 1
const WH_HAS_NO_DISPLAY_PRODUCTS = 0

// Possible statuses for Warehouse status
const WH_IS_ACTIVE = "1"
const WH_IS_NOT_ACTIVE = "0"

//  Possible options for Warehouse isZonal
const WH_IS_ZONAL = 1
const WH_IS_NOT_ZONAL = 0

// Possible statuses for store availability
const WH_STORE_IS_AVAILABLE = 1
const WH_STORE_IS_NOT_AVAILABLE = 0

// Possible statuses for autoconfirm
const WH_AUTOCONFIRM_ENABLED = 1
const WH_AUTOCONFIRM_DISABLED = 0

// Possible statuses for cart synchronisation
const WH_SYNC_CART_ENABLED = 1
const WH_SYNC_CART_DISABLED = 0
const WH_THRESHOLD_DEFAULT = 0

type WarehousesList struct {
	Response
	WarehousesData Warehouses `json:"data"`
}

type Warehouses struct {
	ecom_base.Metadata `json:"metadata"`
	Warehouses         []Warehouse `json:"warehouses"`
}

type Warehouse struct {
	Id           int    `json:"id"` //erp id
	CostCenterId int    `json:"cost_center_id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	CountryCode  string `json:"country_code"`
	RegionCode   string `json:"region_code"`
	RegionName   string `json:"region_name"`
	IsZonal      bool   `json:"is_zonal"`
	IsStore      int    `json:"is_store"`
	// Array of region codes for zonal warehouses
	ConfirmRegions     []string `json:"confirm_regions"`
	SyncCartEnabled    int      `json:"sync_cart_enabled"`
	AutoconfirmEnabled int      `json:"autoconfirm_enabled"`
	// Maximum quantity that can be auto confirmed from a warehouse
	AutoconfirmThreshold int `json:"autoconfirm_threshold"`
	// Possible values: 0 - is not primary warehouse, 1 - is primary warehouse
	AutoconfirmIsPrimary int `json:"autoconfirm_is_primary"`
	// Possible values: 0 - has no product to display, 1 - has products to display
	HasDisplayProducts int `json:"has_display_products"`
	// Warehouse is available for ("Verifica stoc in magazin")
	StoreAvailabilityEnabled int `json:"store_availability_enabled"`
	Status                   int `json:"status"`
	Priority                 int `json:"priority"`
	// Period in days eg. "1-2"
	DeliveryEstimate string    `json:"delivery_estimate"`
	DefaultCourier   string    `json:"default_courier"`
	Couriers         []Courier `json:"couriers"`
}

type Courier struct {
	Name          string `json:"name"`
	AttributeSets []int  `json:"attribute_sets"`
	Status        int    `json:"status"`
}
