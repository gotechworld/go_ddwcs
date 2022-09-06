package dto

/*
{
    "message": "", -- used for error messages
    "data": {
      "estimated_delivery": {
        "SMTG973FZGD": {
          "bundle_id": "MM-0000081-7", -- string, optional
          "estimated_delivery_days": "2-3", -- string, mandatory period in days
          "deliver_from_region": "Bucuresti", -- string, mandatory WH region_name
          "carrier": "Urgent Curier" -- string, mandatory
        },
        "DOCJBLGO2BLK": {
          "bundle_id": "MM-0000081-7",
          "estimated_delivery_days": "5-7",
          "deliver_from_region": "Bucuresti",
          "carrier": "Altex"
        },
        "LAP81D2009BRM": {
          "estimated_delivery_days": "1-2",
          "deliver_from_region": "Bucuresti",
          "carrier": "Altex"
        }
      },
      "generic": {
        "estimated_delivery": "1 - 4", -- string, mandatory computed as "min(estimated_delivery_days[][start]) - max(estimated_delivery_days[][end])"
        "estimated_delivery_dates": [ -- computed as
          "2021-02-06", -- start from current day YYYY-MM-DD (2021-02-05) + estimated_delivery[start] in this case "1" because we have "1 - 4"
          "2021-02-07",
          "2021-02-08",
          "2021-02-09"  -- end from current day YYYY-MM-DD + estimated_delivery[end] in this case "4" ("1 - 4")
        ],
        "additional_delivery_dates": [
          "2021-02-10", -- 14 dates starting from current day YYYY-MM-DD + estimated_delivery[end] + 1
          "2021-02-12", -- skipping sundays and national holidays (national holidays are fetched from GlobalConfig)
          ...
        ],
		"ship_to_store": 0, -- int, mandatory, 1 - if all products have enough stocks in Central WHs (Primary & enabled), 0 - otherwise
      }
    }
}
*/
//type EstimateResponsePost struct {
//	Message string `json:"message"`
//	Data    Data   `json:"data"`
//}

type EstimateResponsePost struct {
	EstimatedDelivery EstimatedDelivery `json:"estimated_delivery"`
	Generic           GenericDelivery   `json:"generic"`
}

type EstimatedDelivery = map[string]EstimatedProductInfo

type EstimatedProductInfo struct {
	BundleId              string `json:"bundle_id"`
	EstimatedDeliveryDays string `json:"estimated_delivery_days"`
	DeliverFromRegion     string `json:"deliver_from_region"`
	Carrier               string `json:"carrier"`
}

type GenericDelivery struct {
	EstimatedDelivery       string   `json:"estimated_delivery"`
	EstimatedDeliveryDates  []string `json:"estimated_delivery_dates"`
	AdditionalDeliveryDates []string `json:"additional_delivery_dates"`
	ShipToStore             int      `json:"ship_to_store"`
}
