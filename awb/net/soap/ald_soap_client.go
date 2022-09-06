package soap

import (
	"encoding/xml"
	"gitlab.altex.ro/ams/go_ddwcs/awb/service/mapper"
	lib "gitlab.altex.ro/plug/go_atx_lib"
	logger "gitlab.altex.ro/plug/go_logger"
	"os"
)

// Request -  getAwbInformations template
type GetAwbInformations struct {
	XMLName xml.Name `xml:"GetAwbInformations"`
	Xmlns   string   `xml:"xmlns,attr"`
	Session string   `xml:"session"`
	ListAwb ListAwb  `xml:"listAwb"`
}

type ListAwb struct {
	Ns3    string   `xml:"xmlns:ns3,attr"`
	String []string `xml:"ns3:string"`
}

// Request - Login details
type Login struct {
	XMLName  xml.Name `xml:"Login"`
	Xmlns    string   `xml:"xmlns,attr"`
	User     string   `xml:"user"`
	Password string   `xml:"password"`
}

const (
	namespace         = "http://tempuri.org/"
	loginAction       = "http://tempuri.org/IInformationsService/Login"
	readAwbInfoAction = "http://tempuri.org/IInformationsService/GetAwbInformations"
)

/**
 * extends HttpReadOnlyClient
 */
type AldSoapClient struct {
	logger         logger.Logger
	responseMapper mapper.Mapper
}

/**
 * InitAldSoapClient
 * @desc - init a new oms client by validating the provided host and populating the credentials for a future request
 */
func InitAldSoapClient(logger logger.Logger, responseMapper mapper.Mapper) *AldSoapClient {
	client := &AldSoapClient{
		logger:         logger,
		responseMapper: responseMapper,
	}

	return client
}

/**
 * Login
 * @return string
 * @desc - login and retrieve an authorization token
 */
func (aldClient *AldSoapClient) Login() string {
	soapClient := NewSoapClient(lib.GetContainer().GetConfig().SystemVars["ald_host"])
	if soapClient == nil {
		aldClient.logger.Error("[AldClient]: could not instantiate soap client")
		return ""
	}

	soapClient.SetToHeader(lib.GetContainer().GetConfig().SystemVars["ald_host"])
	soapClient.SetActionHeader(loginAction)
	soapClient.SetBodyValue(Login{
		Xmlns:    namespace,
		User:     os.Getenv("ALD_USERNAME"),
		Password: os.Getenv("ALD_PASSWORD"),
	})

	response, err := soapClient.ExecuteRequest()
	if err != nil {
		aldClient.logger.Error("[AldClient]:" + err.Error())
		return ""
	}

	if aldClient.responseMapper == nil {
		aldClient.logger.Println("[AldClient]: You need to use struct instance method")
		return ""
	}

	return aldClient.responseMapper.MapLoginToken(response)
}

/**
 * GetAwbInformation
 * @param awbCode String
 * @return AwbInfoResponseEnvelope
 * @desc - reads the awb info from external service
 */
func (aldClient *AldSoapClient) GetAwbInformation(awbCode, token string) *mapper.AldAwbInformationResponse {
	soapClient := NewSoapClient(lib.GetContainer().GetConfig().SystemVars["ald_host"])
	if soapClient == nil {
		aldClient.logger.Error("[AldClient]: could not instantiate soap client")
		return nil
	}

	soapClient.SetToHeader(lib.GetContainer().GetConfig().SystemVars["ald_host"])
	soapClient.SetActionHeader(readAwbInfoAction)
	soapClient.SetBodyValue(GetAwbInformations{
		Xmlns:   namespace,
		Session: token,
		ListAwb: ListAwb{
			Ns3:    "http://schemas.microsoft.com/2003/10/Serialization/Arrays",
			String: []string{awbCode},
		},
	})

	response, err := soapClient.ExecuteRequest()
	if err != nil {
		aldClient.logger.Error("[AldClient]:" + err.Error())
		return nil
	}

	if aldClient.responseMapper == nil {
		aldClient.logger.Println("[AldClient]: You need to use struct instance method")
		return nil
	}

	mapped := aldClient.responseMapper.MapAwbInfo(response, awbCode).(mapper.AldAwbInformationResponse)
	return &mapped

}
