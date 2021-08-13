package message

import (
	"encoding/base64"
	"encoding/json"
	"github.com/tuya/tuya-connector-go/connector/logger"
	"github.com/tuya/tuya-connector-go/connector/message/event"
	"github.com/tuya/tuya-connector-go/connector/utils"
	"reflect"
)

type messageFunc func([]byte)

func (c *client) receiveMsg() messageFunc {
	return func(msg []byte) {
		defer func() {
			if v := recover(); v != nil {
				logger.Log.Errorf("distribute event message failed, err=%+v", v)
			}
		}()
		m := map[string]interface{}{}
		err := json.Unmarshal(msg, &m)
		if err != nil {
			logger.Log.Errorf("json unmarshal failed, err=%s", err.Error())
			return
		}
		protocol := int64(m["protocol"].(float64))
		bs := m["data"].(string)
		de, err := base64.StdEncoding.DecodeString(bs)
		if err != nil {
			logger.Log.Errorf("base64 decode failed, err=%s", err.Error())
			return
		}
		deData := utils.AesEcbDecrypt(de, []byte(c.cfg.accessKey[8:24]))
		err = json.Unmarshal(deData, &m)
		if err != nil {
			logger.Log.Errorf("de json unmarshal failed, err=%s", err.Error())
			return
		}
		if protocol == event.PROTOCOL_STATUS {
			c.mu.RLock()
			f, ok := c.eventSubPool[event.STATUS_REPORT_MESSAGE]
			c.mu.RUnlock()
			if !ok {
				return
			}
			fv := reflect.ValueOf(f)
			if fv.Kind() != reflect.Func {
				return
			}
			m := &event.StatusReportMessage{}
			err := json.Unmarshal(deData, &m)
			if err != nil {
				logger.Log.Errorf("protocol %d json unmarshal failed, err=%s", protocol, err.Error())
			}
			params := []reflect.Value{reflect.ValueOf(m)}
			fv.Call(params)
			return
		} else if protocol == event.PROTOCOL_DEVICE {
			if code, ok := m["bizCode"]; ok {
				c.mu.RLock()
				f, ok := c.eventSubPool[event.GetMessageNameByType(code.(string))]
				c.mu.RUnlock()
				if !ok {
					return
				}
				c.switchCode(code.(string), deData, f)
			}
		} else {
			logger.Log.Warnf("please contact tuya technical support, protocol=%d", protocol)
		}
		return
	}
}

func (c *client) switchCode(code string, deData []byte, f interface{}) {
	fv := reflect.ValueOf(f)
	if fv.Kind() != reflect.Func {
		return
	}
	var m interface{}
	params := make([]reflect.Value, 1)
	switch code {
	case event.ONLINE:
		m = &event.OnlineMessage{}
	case event.OFFLINE:
		m = &event.OfflineMessage{}
	case event.NAME_UPDATE:
		m = &event.NameUpdateMessage{}
	case event.DP_NAME_UPDATE:
		m = &event.DpNameUpdateMessage{}
	case event.DELETE:
		m = &event.DeleteMessage{}
	case event.BIND_USER:
		m = &event.BindUserMessage{}
	case event.UPGRADE_STATUS:
		m = &event.UpgradeStatusMessage{}
	case event.SHARE:
		m = &event.ShareMessage{}
	case event.DEVICE_SIGNAL:
		m = &event.DeviceSignalMessage{}
	case event.USER_REGISTER:
		m = &event.UserRegisterMessage{}
	case event.USER_UPDATE:
		m = &event.UserUpdateMessage{}
	case event.USER_DELETE:
		m = &event.UserDeleteMessage{}
	case event.HOME_CREATE:
		m = &event.HomeCreateMessage{}
	case event.HOME_UPDATE:
		m = &event.HomeUpdateMessage{}
	case event.HOME_DELETE:
		m = &event.HomeDeleteMessage{}
	case event.ROOM_CREATE:
		m = &event.RoomCreateMessage{}
	case event.ROOM_DELETE:
		m = &event.RoomDeleteMessage{}
	case event.ROOM_NAME_UPDATE:
		m = &event.RoomNameUodateMessage{}
	case event.ROOM_SORT:
		m = &event.RoomSortMessage{}
	case event.DEVICE_DP_COMMAND:
		m = &event.DeviceDpCommandMessage{}
	case event.STATUS_REPORT:
		m = &event.StatusReportMessage{}
	case event.AUTOMATION_EXTERNAL_ACTION:
		m = &event.AutomationExternalActionMessage{}
	case event.SCENE_EXECUTE:
		m = &event.SceneExecuteMessage{}
	default:
		params[0] = reflect.ValueOf(string(deData))
		fv.Call(params)
		return
	}
	err := json.Unmarshal(deData, &m)
	if err != nil {
		logger.Log.Errorf("%s json unmarshal failed, err=%s", code, err.Error())
	}
	params[0] = reflect.ValueOf(m)
	fv.Call(params)
}
