package dto

/**
{
  "website_code": "altex",
  "products": [
    {
      "sku": "sku001",
      "qty": 2
    },
    {
      "sku": "sku002",
      "qty": 1
    }
  ]
}
*/
type ReservationStorePayload struct {
	WebsiteCode string           `json:"website_code"`
	Products    []ProductDetails `json:"products"`
}

type ProductDetails struct {
	Sku string `json:"sku"`
	Qty int    `json:"qty"`
}
