package constraints

import (
	"gitlab.altex.ro/ams/go_ddwcs/common/dto/global_config"
)

type ReadGlobalConfigStoreInterface interface {
	GetConfig(params map[string]string) global_config.ConfigGet
}
