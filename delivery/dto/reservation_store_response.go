package dto

import "gitlab.altex.ro/ams/go_ddwcs/common/dto/cms"

/**
{
  "stores": [
    {
      "name": "Billa",
      "erp_id": 208
    },
    {
      "name": "Plaza Romania",
      "erp_id": 233
    }
  ]
}
*/

type StoreContainer struct {
	Stores []cms.Store `json:"stores"`
}

func NewStoreContainer() *StoreContainer {
	return &StoreContainer{
		Stores: []cms.Store{},
	}
}

/**
 * @param store Store
 * return reference of StoreContainer (fluent setter)
 */
func (s *StoreContainer) AddStore(store cms.Store) *StoreContainer {
	s.Stores = append(s.Stores, store)

	return s
}
