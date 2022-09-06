package envelope

import "encoding/xml"

// Response
type BodyLessResponseEnvelope struct {
	XMLName xml.Name       `xml:"Envelope"`
	S       string         `xml:"s,attr"`
	A       string         `xml:"a,attr"`
	Header  ResponseHeader `xml:"Header"`
}

type ResponseHeader struct {
	Action    MustUnderstandField `xml:"Action"`
	RelatesTo string              `xml:"RelatesTo"`
}

// Request
type RequestEnvelope struct {
	XMLName xml.Name      `xml:"env:Envelope"`
	Env     string        `xml:"xmlns:env,attr"`
	Header  RequestHeader `xml:"env:Header"`
	Body    Body          `xml:"env:Body"`
}

type RequestHeader struct {
	Wsa       string              `xml:"xmlns:wsa,attr"`
	To        string              `xml:"wsa:To"`
	Action    MustUnderstandField `xml:"wsa:Action"`
	MessageID string              `xml:"wsa:MessageID"`
}

//common
type Body struct {
	Text      string `xml:",chardata"`
	BodyValue interface{}
}

type MustUnderstandField struct {
	Text           string `xml:",chardata"`
	MustUnderstand string `xml:"mustUnderstand,attr"`
}
