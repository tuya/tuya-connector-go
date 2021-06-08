package model

type DeviceModel struct {
	UUID   string `json:"uuid"`
	UID    string `json:"uid"`
	Name   string `json:"name"`
	IP     string `json:"ip"`
	Sub    bool   `json:"sub"`
	Model  string `json:"model"`
	Status []struct {
		Code  string      `json:"code"`
		Value interface{} `json:"value"`
	} `json:"status"`
	Category    string `json:"category"`
	Online      bool   `json:"online"`
	ID          string `json:"id"`
	TimeZone    string `json:"time_zone"`
	LocalKey    string `json:"local_key"`
	UpdateTime  int    `json:"update_time"`
	ActiveTime  int    `json:"active_time"`
	OwnerID     string `json:"owner_id"`
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
}

type GetDeviceResponse struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Success bool        `json:"success"`
	Result  DeviceModel `json:"result"`
	T       int64       `json:"t"`
}

type PostDeviceCmdResponse struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	Result  bool   `json:"result"`
	T       int64  `json:"t"`
}
