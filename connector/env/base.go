package env

import (
	"os"
	"sync/atomic"
)

const (
	TUYA_API_HOST  = "TUYA_API_HOST"
	TUYA_MSG_HOST  = "TUYA_MSG_HOST"
	TUYA_ACCESSID  = "TUYA_ACCESSID"
	TUYA_ACCESSKEY = "TUYA_ACCESSKEY"
)

var (
	Config *environment
)

type environment struct {
	done uint32
	// message subscription
	// msg host
	msgHost string
	// api host
	apiHost string
	// api ak
	accessID string
	// api sk
	accessKey string
	appName   string
	// if true, debug mod
	isDebug bool
}

func NewEnv() *environment {
	return &environment{
		isDebug: true,
	}
}

func (env *environment) Init() *environment {
	if atomic.LoadUint32(&env.done) == 1 {
		return env
	}
	defer atomic.StoreUint32(&env.done, 1)
	if env.appName == "" {
		env.appName = "tysdk"
	}
	//if the token handle is nil, use default handle
	if env.apiHost == "" {
		env.apiHost = os.Getenv(TUYA_API_HOST)
		if env.apiHost == "" {
			panic("no set api host")
		}
	}
	if env.msgHost == "" {
		env.msgHost = os.Getenv(TUYA_MSG_HOST)
		if env.msgHost == "" {
			panic("no set msg host")
		}
	}
	if env.accessID == "" {
		env.accessID = os.Getenv(TUYA_ACCESSID)
		if env.accessID == "" {
			panic("no set access id")
		}
	}
	if env.accessKey == "" {
		env.accessKey = os.Getenv(TUYA_ACCESSKEY)
		if env.accessKey == "" {
			panic("no set access key")
		}
	}
	return env
}

func (env *environment) GetApiHost() string {
	return env.apiHost
}

func (env *environment) GetMsgHost() string {
	return env.msgHost
}

func (env *environment) GetAccessID() string {
	return env.accessID
}

func (env *environment) GetAccessKey() string {
	return env.accessKey
}

func (env *environment) GetAppName() string {
	return env.appName
}

func (env *environment) DebugMode() bool {
	return env.isDebug
}
