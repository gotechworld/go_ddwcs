package mapper

import (
	"encoding/xml"
	"gitlab.altex.ro/ams/go_ddwcs/awb/net/soap/envelope"
)

// AwbInfo
type AldAwbInformationResponse struct {
	envelope.BodyLessResponseEnvelope
	Body struct {
		XMLName                    xml.Name `xml:"Body"`
		GetAwbInformationsResponse GetAwbInformationsResponse
	}
}

type GetAwbInformationsResponse struct {
	XmlName                  xml.Name `xml:"GetAwbInformationsResponse"`
	Text                     string   `xml:",chardata"`
	Xmlns                    string   `xml:"xmlns,attr"`
	GetAwbInformationsResult struct {
		B            string         `xml:"b,attr"`
		I            string         `xml:"i,attr"`
		Informations []Informations `xml:"Informations"`
	} `xml:"GetAwbInformationsResult"`
}

type Informations struct {
	XMLName   xml.Name `xml:"Informations"`
	Text      string   `xml:",chardata"`
	Awb       string   `xml:"Awb"`
	AwbChild  string   `xml:"AwbChild"`
	CodFiscal struct {
		Text string `xml:",chardata"`
		Nil  string `xml:"nil,attr"`
	} `xml:"CodFiscal"`
	Curier     string `xml:"Curier"`
	DataStatus string `xml:"DataStatus"`
	StatusAwb  string `xml:"StatusAwb"`
	UserName   string `xml:"UserName"`
}

//login
type LoginResponseEnvelope struct {
	envelope.BodyLessResponseEnvelope
	Body struct {
		LoginResponse LoginResponse `xml:"LoginResponse"`
	} `xml:"Body"`
}

type LoginResponse struct {
	Xmlns       string `xml:"xmlns,attr"`
	LoginResult string `xml:"LoginResult"`
}

type AldResponseMapper struct{}

/**
 * MapAwbInfo
 * @param data []byte
 * @return interface|AldAwbInformationResponse
 * @desc - maps the provided xml
 */
func (mapper *AldResponseMapper) MapAwbInfo(data []byte, awbCode string) interface{} {
	var response AldAwbInformationResponse
	err := xml.Unmarshal(data, &response)
	if err != nil {
		return response
	}

	return response
}

/**
 * MapLoginToken
 * @param data []byte
 * @return string
 * @desc - maps the login result
 */
func (mapper *AldResponseMapper) MapLoginToken(data []byte) string {
	var response LoginResponseEnvelope
	err := xml.Unmarshal(data, &response)
	if err != nil {
		return ""
	}

	return response.Body.LoginResponse.LoginResult
}
