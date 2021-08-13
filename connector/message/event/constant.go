package event

const (
	// device online
	ONLINE = "online"
	// device offline
	OFFLINE = "offline"
	// update device name
	NAME_UPDATE = "nameUpdate"
	// modify device function name
	DP_NAME_UPDATE = "dpNameUpdate"
	// unbind device
	DELETE = "delete"
	// bind device
	BIND_USER = "bindUser"
	// device upgrade status
	UPGRADE_STATUS = "upgradeStatus"
	// device share
	SHARE = "share"
	// device signal
	DEVICE_SIGNAL = "deviceSignal"
	// user register
	USER_REGISTER = "userRegister"
	// user data update
	USER_UPDATE = "userUpdate"
	// user logout
	USER_DELETE = "userDelete"
	// create home
	HOME_CREATE = "homeCreate"
	// update home
	HOME_UPDATE = "homeUpdate"
	// delete home
	HOME_DELETE = "homeDelete"
	// create room
	ROOM_CREATE = "roomCreate"
	// delete room
	ROOM_DELETE = "roomDelete"
	// update room name
	ROOM_NAME_UPDATE = "roomNameUpdate"
	// room sort
	ROOM_SORT = "roomSort"
	// device control command
	DEVICE_DP_COMMAND = "deviceDpCommand"
	// report device status data
	STATUS_REPORT              = "statusReport"
	AUTOMATION_EXTERNAL_ACTION = "automationExternalAction"
	SCENE_EXECUTE              = "sceneExecute"
)

const (
	ONLINE_MESSAGE                     = "OnlineMessage"
	OFFLINE_MESSAGE                    = "OfflineMessage"
	NAME_UPDATE_MESSAGE                = "NameUpdateMessage"
	DP_NAME_UPDATE_MESSAGE             = "DpNameUpdateMessage"
	DELETE_MESSAGE                     = "DeleteMessage"
	BIND_USER_MESSAGE                  = "BindUserMessage"
	UPGRADE_STATUS_MESSAGE             = "UpgradeStatusMessage"
	SHARE_MESSAGE                      = "ShareMessage"
	DEVICE_SIGNAL_MESSAGE              = "DeviceSignalMessage"
	USER_REGISTER_MESSAGE              = "UserRegisterMessage"
	USER_UPDATE_MESSAGE                = "UserUpdateMessage"
	USER_DELETE_MESSAGE                = "UserDeleteMessage"
	HOME_CREATE_MESSAGE                = "HomeCreateMessage"
	HOME_UPDATE_MESSAGE                = "HomeUpdateMessage"
	HOME_DELETE_MESSAGE                = "HomeDeleteMessage"
	ROOM_CREATE_MESSAGE                = "RoomCreateMessage"
	ROOM_DELETE_MESSAGE                = "RoomDeleteMessage"
	ROOM_NAME_UPDATE_MESSAGE           = "RoomNameUpdateMessage"
	ROOM_SORT_MESSAGE                  = "RoomSortMessage"
	DEVICE_DP_COMMAND_MESSAGE          = "DeviceDpCommandMessage"
	STATUS_REPORT_MESSAGE              = "StatusReportMessage"
	AUTOMATION_EXTERNAL_ACTION_MESSAGE = "AutomationExternalActionMessage"
	SCENE_EXECUTE_MESSAGE              = "SceneExecuteMessage"
)

const (
	PROTOCOL_STATUS = 4
	PROTOCOL_DEVICE = 20
)
