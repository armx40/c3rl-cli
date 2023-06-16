package c3rliotroxy

const (
	APIv2CODEOK                       = 0
	APIv2CODEFAIL                     = 1
	APIv2INVALIDGTYPE                 = 2
	APIv2AUTHENTICATIONFAILED         = 3
	APIv2NOAUTHENTICATIONTOKENPRESENT = 4
)

type general_payload_v2_t struct {
	Payload interface{} `json:"payload" validate:"required"`
	Status  string      `json:"status" validate:"required"`
	Code    int         `json:"code" validate:"required"`
}

type expose_response_payload_t struct {
	Status             string `json:"status" validate:"required"`
	Code               int    `json:"code" validate:"required"`
	ExposedPorts       map[uint16]uint16
	ExposedDomainPorts map[uint16]string
}

type request_response_request_packet_t struct {
	GType    string      `json:"g"`
	Data     interface{} `json:"d"`
	Sequence uint32      `json:"s"`
}
