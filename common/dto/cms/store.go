package cms

type StoreList struct {
	Items []Store `json:"items"`
}

type Store struct {
	Name  string `json:"name"`
	ErpID int    `json:"erp_id"`
}
