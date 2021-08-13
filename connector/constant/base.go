package constant

const (
	REQ_INFO          = "requestInfo"
	Signature_Headers = "Signature-Headers"
)

const (
	TOKEN = "token"
	TS    = "ts"
	NONCE = "nonce"
)

const (
	ExeCount = "exeCnt"
)

const (
	Header_ContentType = "Content-Type"
	Header_SignMethod  = "sign_method"
	Header_ClientID    = "client_id"
	Header_TimeStamp   = "t"
	Header_Sign        = "sign"
	Header_Nonce       = "nonce"
	Header_AccessToken = "access_token"
	Header_DevChannel  = "Dev_channel"
	Header_DevLang     = "Dev_lang"

	ContentType_JSON = "application/json"
	SignMethod_HMAC  = "HMAC-SHA256"
	Dev_Channel      = "SaaSFramework"
	Dev_Lang         = "golang"
)

const (
	TUYA_HEADER     = "tuya_header"
	TUYA_SIGN       = "tuya_sign"
	TUYA_LOG        = "tuya_log"
	TUYA_TOKEN      = "tuya_token"
	TUYA_MESSAGE    = "tuya_message"
	TUYA_ERROR_PROC = "tuya_error_proc"
)

const (
	// system error,please contact the admin
	SYSTEM_ERROR = 500
	// data not exist
	DATA_NOT_EXIST = 1000
	// secret invalid
	SECRET_INVALID = 1001
	// token is expired
	TOKEN_EXPIRED = 1010
)
