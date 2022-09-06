package constraints

import "gitlab.altex.ro/ams/go_ddwcs/common/dto/stocks"

type ReadWarehouseStoreInterface interface {
	GetWarehouses(map[string]string) []stocks.Warehouse
}
