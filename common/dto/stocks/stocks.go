package stocks

import "gitlab.altex.ro/ams/go_ddwcs/common/dto/ecom_base"

// struct that maps the stock details from StocksAPI GET "/api/v1.0/stocks" response
/*
* "data": [
*   "<product_SKU>": [
*      "<seller_id>": [
*          "inventory": [
*              "status": <status>,
*              "quantity": <quantity>
*          ],
*
*          "warehouse_inventory": [
*              "<warehouse_id>": <warehouse_quantity>,
*              "<warehouse_id>": <warehouse_quantity>,
*              "<warehouse_id>": <warehouse_quantity>,
*              "<supplier_warehouse_id>": ["<supplier_id>": <quantity>, ...]
*          ]
*      ]
*   ], ...
* ]
 */

const STATUS_OUT_STOCK = 0
const STATUS_IN_STOCK = 1
const STATUS_IN_SUPPLIER_STOCK = 2
const STATUS_PREORDER = 3
const STATUS_EOL = 4

//{
//	"messages":[],
//	"status":"success",
//	"data":{
//		"metadata":{"total_items":2,"total_pages":1,"items_per_page":100,"current_page":1},
//		"stocks":{
//			"ACC01075":{
//				"1":{
//					"inventory":{"status":1,"available_in_store":1,"quantity":871},
//					"warehouse_inventory":{"16":1,"101":19}
//				}
//			},
//			"AIOMRR12ROA":{
//				"1":{
//					"inventory":{"status":3,"available_in_store":0,"quantity":9999},
//					"warehouse_inventory":{"3000":50000, "97": ["SUPPLIER_CODE_1":20, "SUPPLIER_CODE_2":70]}
//				}
//			}
//		}//stocks
//	} //data
//}
type StockList struct {
	Response
	StocksData StockData `json:"data"`
}

type StockData struct {
	Metadata ecom_base.Metadata     `json:"metadata"`
	Stocks   map[string]SellerStock `json:"stocks"`
}

type SellerStock = map[int]Stock

type Stock struct {
	Inventory          Inventory              `json:"inventory"`
	WarehouseInventory map[string]interface{} `json:"warehouse_inventory"`
}

type Inventory struct {
	Status           int `json:"status"`
	AvailableInStore int `json:"available_in_store"`
	Quantity         int `json:"quantity"`
}
