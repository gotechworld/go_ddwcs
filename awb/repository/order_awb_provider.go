package repository

import (
	"gitlab.altex.ro/ams/go_ddwcs/awb/data_store"
	"gitlab.altex.ro/ams/go_ddwcs/awb/dto"
	"gitlab.altex.ro/ams/go_ddwcs/awb/util"
	commonDataStore "gitlab.altex.ro/ams/go_ddwcs/common/data_store"
)

type OrderAwbsRepository struct {
	orderIncrementId string
	customer         struct {
		value     string
		valueType string
	}
	storePool   []commonDataStore.DataStoreReadOnly
	awbInfoList []dto.AwbInfo
}

/**
 * NewOrderAwbsRepo - Constructor
 * @desc - Create a new awb info awbInfoList provider - need to fail fast in case of empty params
 */
func NewOrderAwbsRepo(orderIncrementId, customer string) *OrderAwbsRepository {
	// add increment id
	provider := &OrderAwbsRepository{
		orderIncrementId: orderIncrementId,
	}

	//add customer details
	provider.customer = struct {
		value     string
		valueType string
	}{value: customer, valueType: provider.extractCustomerInfoType(customer)}

	// create a data store pool
	provider.storePool = data_store.GetDataStoreList()

	// init the list of available awbs
	provider.awbInfoList = provider.initList()

	return provider
}

/**
 * initList - Private
 * @return []dto.AwbInfo
 * @desc - Full fill the list of awbs for future purpose
 */
func (provider *OrderAwbsRepository) initList() []dto.AwbInfo {
	params := map[string]string{
		"increment_id": provider.orderIncrementId,
		"customer_" + provider.customer.valueType: provider.customer.value,
	}

	var list []dto.AwbInfo
	if len(provider.storePool) == 0 {
		return list
	}

	// iterate over the data store and get the appropriate result
	for _, store := range provider.storePool {
		list = store.GetList(params)
		if len(list) > 0 {
			break
		}
	}

	// add increment id to the info
	for idx := range list {
		list[idx].OrderIncrementId = provider.orderIncrementId
	}

	return list
}

/**
 * GetList
 * @return []dto.AwbInfo
 * @desc - Get a awbInfoList of awb info provided by OMS using order increment id and customer filters
 */
func (provider *OrderAwbsRepository) GetList() []dto.AwbInfo {
	return provider.awbInfoList
}

/**
 * Get a awbInfoList of awb info provided by OMS using order increment id and customer filters
 */
func (provider *OrderAwbsRepository) Get(awbNumber string) *dto.AwbInfo {
	for _, awbDetails := range provider.awbInfoList {
		if awbDetails.AwbNumber == awbNumber {
			return &awbDetails
		}
	}

	return nil
}

/**
 * extractCustomerInfoType - Private
 * @return string
 * @desc - extracts the customer value type [id, email]
 */
func (provider *OrderAwbsRepository) extractCustomerInfoType(customer string) string {
	customerInfoType := "id"
	if util.IsValidEmail(customer) {
		customerInfoType = "email"
	}

	return customerInfoType
}
