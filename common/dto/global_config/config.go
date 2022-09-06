package global_config

//{
//	"messages": [],
//	"status": "success",
//	"data": [
//		{
//		"scope": 1,
//		"platform": "stock",
//		"code": "some_code",
//		"value": "some_code value",
//		"backend_type": "encrypted",
//		"label": "descriere",
//		"validation": [
//				{
//					"rule": "required",
//					"error": "Acest camp este obligatoriu"
//				}
//			]
//		}
//	]
//}
type ConfigGet struct {
	Response
	ConfigGetData []ConfigGetData `json:"data"`
}

type ConfigGetData struct {
	Value       string `json:"value"`
	BackendType string `json:"backend_type"`
	Label       string `json:"label"`
}
