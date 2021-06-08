package httplib

const (
	// cn
	URL_CN = "https://openapi.tuyacn.com"
	MSG_CN = "pulsar+ssl://mqe.tuyacn.com:7285/"
	// us
	URL_US = "https://openapi.tuyaus.com"
	MSG_US = "pulsar+ssl://mqe.tuyaus.com:7285/"
	// eu
	URL_EU = "https://openapi.tuyaeu.com"
	MSG_EU = "pulsar+ssl://mqe.tuyaeu.com:7285/"
	// in
	URL_IN = "https://openapi.tuyain.com"
	MSG_IN = "pulsar+ssl://mqe.tuyain.com:7285/"
)

type response struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	T       int64  `json:"t"`
}
