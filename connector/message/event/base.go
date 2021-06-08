package event

import "encoding/json"

type baseMessage struct {
	BizCode    string   `json:"bizCode"`
	BizData    bizDataM `json:"bizData"`
	DevID      string   `json:"devId"`
	ProductKey string   `json:"productKey"`
	Ts         int64    `json:"ts"`
	Uuid       string   `json:"uuid"`
}

type bizDataM struct {
	DevID string `json:"devId"`
	Uid   string `json:"uid"`
	Name  string `json:"name"`
}

func (m baseMessage) toString() string {
	b, _ := json.Marshal(m)
	return string(b)
}

type NameUpdateMessage struct {
	baseMessage
}

type AutomationExternalActionMessage struct {
	baseMessage
	AutomationId string `json:"automationId"`
}

type BindUserMessage struct {
	baseMessage
}

type DeleteMessage struct {
	baseMessage
}

type DeviceDpCommandMessage struct {
	baseMessage
	Command []commandItem `json:"command"`
}

type commandItem struct {
	DpID  int64 `json:"dpId"`
	Value bool  `json:"value"`
}

type DeviceSignalMessage struct {
	baseMessage
	ReportData []deviceSignalItem `json:"reportData"`
}

type deviceSignalItem struct {
	Memory int64 `json:"memory"`
	Rssi   int   `json:"rssi"`
	T      int64 `json:"t"`
}

type DpNameUpdateMessage struct {
	baseMessage
}

type HomeCreateMessage struct {
	baseMessage
}

type HomeDeleteMessage struct {
	baseMessage
}

type HomeUpdateMessage struct {
	baseMessage
}

type OfflineMessage struct {
	baseMessage
}

type OnlineMessage struct {
	baseMessage
}

type RoomCreateMessage struct {
	baseMessage
}

type RoomDeleteMessage struct {
	baseMessage
}

type RoomNameUodateMessage struct {
	baseMessage
}

type RoomSortMessage struct {
	baseMessage
}

type SceneExecuteMessage struct {
	baseMessage
	Gid     int64              `json:"gid"`
	Uid     string             `json:"uid"`
	Actions []sceneExecuteItem `json:"actions"`
}

type sceneExecuteItem struct {
	EntityID    string `json:"entityId"`
	ExecStatus  int    `json:"execStatus"`
	ExecuteTime int64  `json:"executeTime"`
	ID          string `json:"id"`
}

type ShareMessage struct {
	baseMessage
}

type StatusReportMessage struct {
	baseMessage
	DataId string             `json:"dataId"`
	Status []statusReportItem `json:"status"`
}

type statusReportItem struct {
	Code  string      `json:"code"`
	T     int64       `json:"t"`
	Value interface{} `json:"value"`
}

type UpgradeStatusMessage struct {
	baseMessage
}

type UserDeleteMessage struct {
	baseMessage
}

type UserRegisterMessage struct {
	baseMessage
}

type UserUpdateMessage struct {
	baseMessage
}

func GetMessageNameByType(v string) string {
	switch v {
	case ONLINE:
		return ONLINE_MESSAGE
	case OFFLINE:
		return OFFLINE_MESSAGE
	case NAME_UPDATE:
		return NAME_UPDATE_MESSAGE
	case DP_NAME_UPDATE:
		return DP_NAME_UPDATE_MESSAGE
	case DELETE:
		return DELETE_MESSAGE
	case BIND_USER:
		return BIND_USER_MESSAGE
	case UPGRADE_STATUS:
		return UPGRADE_STATUS_MESSAGE
	case SHARE:
		return SHARE_MESSAGE
	case DEVICE_SIGNAL:
		return DEVICE_SIGNAL_MESSAGE
	case USER_REGISTER:
		return USER_REGISTER_MESSAGE
	case USER_UPDATE:
		return USER_UPDATE_MESSAGE
	case USER_DELETE:
		return USER_DELETE_MESSAGE
	case HOME_CREATE:
		return HOME_CREATE_MESSAGE
	case HOME_UPDATE:
		return HOME_UPDATE_MESSAGE
	case HOME_DELETE:
		return HOME_DELETE_MESSAGE
	case ROOM_CREATE:
		return ROOM_CREATE_MESSAGE
	case ROOM_DELETE:
		return ROOM_DELETE_MESSAGE
	case ROOM_NAME_UPDATE:
		return ROOM_NAME_UPDATE_MESSAGE
	case ROOM_SORT:
		return ROOM_SORT_MESSAGE
	case DEVICE_DP_COMMAND:
		return DEVICE_DP_COMMAND_MESSAGE
	case STATUS_REPORT:
		return STATUS_REPORT_MESSAGE
	case AUTOMATION_EXTERNAL_ACTION:
		return AUTOMATION_EXTERNAL_ACTION_MESSAGE
	case SCENE_EXECUTE:
		return SCENE_EXECUTE_MESSAGE
	}
	return ""
}
