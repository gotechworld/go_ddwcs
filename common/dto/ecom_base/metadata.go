package ecom_base

type Metadata struct {
	TotalItems   int `json:"total_items"`
	TotalPages   int `json:"total_pages"`
	ItemsPerPage int `json:"items_per_page"`
	CurrentPage  int `json:"current_page"`
}
