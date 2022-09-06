package strategy

import "gitlab.altex.ro/ams/go_ddwcs/awb/service/adapter"

type AwbInfoStrategy interface {
	GetAwbInfo() *adapter.AwbStatusDetails
}
