package soap

import (
	"encoding/xml"
	"gitlab.altex.ro/ams/go_ddwcs/awb/net/soap/envelope"
	"gitlab.altex.ro/ams/go_ddwcs/common/net"
	httpLocal "gitlab.altex.ro/plug/go_http"
	"gitlab.altex.ro/plug/go_http/builder"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type SoapClient struct {
	net.HttpClientBaseSetup
	request   *envelope.RequestEnvelope
	MessageID string
}

/**
 * NewSoapClient
 * @param url String
 * @return SoapClient pointer
 * @desc - init details for a future soap request
 */
func NewSoapClient(url string) *SoapClient {
	client := &SoapClient{}
	err := client.SetUrl(url)
	if err != nil {
		return nil
	}

	client.SetHeaders(map[string]string{
		"Content-Type": "application/soap+xml",
		"Accept":       "application/xml",
	})

	client.MessageID = strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Int())
	client.request = &envelope.RequestEnvelope{
		Env: "http://www.w3.org/2003/05/soap-envelope",
		Header: envelope.RequestHeader{
			Wsa:       "http://www.w3.org/2005/08/addressing",
			Action:    envelope.MustUnderstandField{MustUnderstand: "true"},
			MessageID: client.MessageID,
		},
	}

	return client
}

/**
 * SetActionHeader
 * @param action string
 * @return void
 * @desc - set the action soap header
 */
func (soap *SoapClient) SetActionHeader(action string) {
	soap.request.Header.Action.Text = action
}

/**
 * SetToHeader
 * @param to string
 * @return void
 * @desc - set the to soap header
 */
func (soap *SoapClient) SetToHeader(to string) {
	soap.request.Header.To = to
}

/**
 * SetBodyValue
 * @param to string
 * @return void
 * @desc - set the to soap body details
 */
func (soap *SoapClient) SetBodyValue(bodyValue interface{}) {
	soap.request.Body = envelope.Body{
		BodyValue: bodyValue,
	}
}

/**
 * buildRequestBody
 * @return []byte, error
 * @desc - format the data for the next request
 */
func (soap *SoapClient) buildRequestBody() ([]byte, error) {
	return xml.Marshal(soap.request)
}

/**
 * ExecuteRequest
 * @return []byte, error
 * @desc - executes a http request
 */
func (soap *SoapClient) ExecuteRequest() ([]byte, error) {
	req, _ := soap.buildRequestBody()

	request := builder.NewRequestBuilder().Url(soap.GetUrl().String()).Method(http.MethodPost).Headers(*soap.GetHeaders()).WithBody(req).Build()
	client, err := httpLocal.NewClientRequest(request)

	if err != nil {
		return []byte{}, err
	}

	return client.Call()
}
