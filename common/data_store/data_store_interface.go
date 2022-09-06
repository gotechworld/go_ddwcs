package data_store

import "gitlab.altex.ro/ams/go_ddwcs/awb/dto"

type DataStoreReadOnly interface {
	GetList(params map[string]string) []dto.AwbInfo
}

type DataStore interface {
	DataStoreReadOnly
	Save(string, []dto.AwbInfo)
}
