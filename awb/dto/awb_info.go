package dto

// struct that maps the awb details from oms response
type AwbInfo struct {
	OrderIncrementId string
	OnlineSaleId     string `json:"online_sale_id"`
	AwbNumber        string `json:"awb_number"`
	AwbCarrier       string `json:"awb_carrier"`
	AwbStatus        string `json:"awb_status"`
}

type OrderAwbList struct {
	Awbs []AwbInfo `json:"awbs"`
}
