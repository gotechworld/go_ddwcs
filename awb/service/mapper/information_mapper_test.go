package mapper_test

import (
	"gitlab.altex.ro/ams/go_ddwcs/awb/service/mapper"
	"testing"
)

var testCargusPayload = `[
    {
        "Code": "ALXR10722458",
        "Type": "Colet",
        "MeasuredWeight": 12.20,
        "VolumetricWeight": 0.0,
        "ConfirmationName": null,
        "Observation": "LD-000201336 - Comenzi: CV-050090338 ",
        "ResponseCode": "",
        "Event": [
            {
                "Date": "2020-06-06T08:37:48.127",
                "EventId": 234,
                "Description": "PickUp Colectat (info in sistem inexistent) ",
                "LocalityName": "BUCURESTI"
            },
            {
                "Date": "2020-06-06T09:20:41.773",
                "EventId": 234,
                "Description": "PickUp Colectat (info in sistem inexistent) ",
                "LocalityName": "BUCURESTI"
            }
        ]
    },
    {
        "Code": "ALXR10722458002",
        "Type": "Colet",
        "MeasuredWeight": 12.20,
        "VolumetricWeight": 0.0,
        "ConfirmationName": null,
        "Observation": "LD-000201336 - Comenzi: CV-050090338 ",
        "ResponseCode": "",
        "Event": [
            {
                "Date": "2020-06-06T08:37:49.887",
                "EventId": 234,
                "Description": "PickUp Colectat (info in sistem inexistent) ",
                "LocalityName": "BUCURESTI"
            },
            {
                "Date": "2020-06-06T09:20:43.44",
                "EventId": 234,
                "Description": "PickUp Colectat (info in sistem inexistent) ",
                "LocalityName": "BUCURESTI"
            }
        ]
    },
    {
        "Code": "ALXR10722458003",
        "Type": "Colet",
        "MeasuredWeight": 12.20,
        "VolumetricWeight": 0.0,
        "ConfirmationName": null,
        "Observation": "LD-000201336 - Comenzi: CV-050090338 ",
        "ResponseCode": "",
        "Event": [
            {
                "Date": "2020-06-06T08:37:52.403",
                "EventId": 234,
                "Description": "PickUp Colectat (info in sistem inexistent) ",
                "LocalityName": "BUCURESTI"
            },
            {
                "Date": "2020-06-06T09:20:51.507",
                "EventId": 234,
                "Description": "PickUp Colectat (info in sistem inexistent) ",
                "LocalityName": "BUCURESTI"
            }
        ]
    },
    {
        "Code": "ALXR10722458004",
        "Type": "Colet",
        "MeasuredWeight": 12.20,
        "VolumetricWeight": 0.0,
        "ConfirmationName": null,
        "Observation": "LD-000201336 - Comenzi: CV-050090338 ",
        "ResponseCode": "",
        "Event": [
            {
                "Date": "2020-06-06T08:37:54.54",
                "EventId": 234,
                "Description": "PickUp Colectat (info in sistem inexistent) ",
                "LocalityName": "BUCURESTI"
            },
            {
                "Date": "2020-06-06T09:20:58.747",
                "EventId": 234,
                "Description": "PickUp Colectat (info in sistem inexistent) ",
                "LocalityName": "BUCURESTI"
            }
        ]
    },
    {
        "Code": "ALXR10722458005",
        "Type": "Colet",
        "MeasuredWeight": 0.40,
        "VolumetricWeight": 0.0,
        "ConfirmationName": null,
        "Observation": "LD-000201336 - Comenzi: CV-050090338 ",
        "ResponseCode": "",
        "Event": [
            {
                "Date": "2020-06-06T08:37:21.983",
                "EventId": 234,
                "Description": "PickUp Colectat (info in sistem inexistent) ",
                "LocalityName": "BUCURESTI"
            },
            {
                "Date": "2020-06-06T09:20:45.173",
                "EventId": 234,
                "Description": "PickUp Colectat (info in sistem inexistent) ",
                "LocalityName": "BUCURESTI"
            }
        ]
    },
    {
        "Code": "ALXR10722458006",
        "Type": "Colet",
        "MeasuredWeight": 12.20,
        "VolumetricWeight": 0.0,
        "ConfirmationName": null,
        "Observation": "LD-000201336 - Comenzi: CV-050090338 ",
        "ResponseCode": "",
        "Event": [
            {
                "Date": "2020-06-06T08:37:26.46",
                "EventId": 234,
                "Description": "PickUp Colectat (info in sistem inexistent) ",
                "LocalityName": "BUCURESTI"
            },
            {
                "Date": "2020-06-06T09:20:36.18",
                "EventId": 234,
                "Description": "PickUp Colectat (info in sistem inexistent) ",
                "LocalityName": "BUCURESTI"
            }
        ]
    }
]`

func TestCargusResponseMapper_MapAwbInfo(t *testing.T) {
	t.Log("Empty body will return empty CargusAwbInformationResponse")
	cmapper := mapper.CargusResponseMapper{}
	info := cmapper.MapAwbInfo([]byte("[]"), "")

	val, ok := info.(mapper.CargusAwbInformationResponse)
	if !ok {
		t.Log("Cannot convert response to CargusAwbInformationResponse")
		t.FailNow()
	}

	if val.Code != "" && len(val.Event) > 0 {
		t.Log("Wrong information provided")
		t.FailNow()
	}
}

func TestCargusResponseMapper_MapAwbInfoWithListResponse(t *testing.T) {
	t.Log("Awb details list body will return proper CargusAwbInformationResponse")
	cmapper := mapper.CargusResponseMapper{}
	awbCode := "ALXR10722458"
	info := cmapper.MapAwbInfo([]byte(testCargusPayload), awbCode)

	val, ok := info.(mapper.CargusAwbInformationResponse)
	if !ok {
		t.Log("Cannot convert response to CargusAwbInformationResponse")
		t.FailNow()
	}

	if val.Code != awbCode {
		t.Logf("Expected %s, Given %s", awbCode, val.Code)
		t.FailNow()
	}
}
