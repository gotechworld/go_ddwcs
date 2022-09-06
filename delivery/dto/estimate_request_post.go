package dto

/*
 website_code: "altex", -- optional, string ("altex"/"mediagalaxy") default "altex"
 delivery_region: "B", -- optional, string - customer shipping region
  product_sku_1: {
    bundle_id: "MM-0000081-7", -- optional, string
    is_preorder: 0, -- optional, int (0/1), default 0
    attribute_set_id: 4, -- mandatory, int
	qty: 2 -- optional, int default 1
  },
  product_sku_2: {
  ...
*/
type EstimateRequestPost struct {
	WebsiteCode    string           `json:"website_code"`
	DeliveryRegion string           `json:"delivery_region"`
	Products       EstimateProducts `json:"products" validate:"required"`
}

type EstimateProducts = map[string]ProductInfo

type ProductInfo struct {
	BundleId       string `json:"bundle_id"`
	IsPreOrder     int    `json:"is_preorder"`
	AttributeSetId int    `json:"attribute_set_id" validate:"required"`
	Qty            int    `json:"qty"`
}
