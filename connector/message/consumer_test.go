package message

import (
	"fmt"
	"github.com/tuya/tuya-connector-go/connector/env"
	"github.com/tuya/tuya-connector-go/connector/logger"
	"github.com/tuya/tuya-connector-go/connector/message/event"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	fmt.Println("init....")
	env.Config = env.NewEnv()
	env.Config.Init()
	if Handler == nil {
		Handler = NewEventMsgWrapper()
	}
	if logger.Log == nil {
		logger.Log = logger.NewDefaultLogger(env.Config.GetAppName(), env.Config.DebugMode())
	}
	fmt.Println("### iot core init success ###")
	m.Run()
}

func TestEventMsg(t *testing.T) {
	Handler.InitMessageClient()
	Handler.SubEventMessage(func(m *event.NameUpdateMessage) {
		logger.Log.Info("=========== name update： ==========")
		logger.Log.Info(m)
	})
	Handler.SubEventMessage(func(m *event.StatusReportMessage) {
		logger.Log.Info("=========== report data： ==========")
		for _, v := range m.Status {
			logger.Log.Infof(v.Code, v.Value)
		}
	})

	time.Sleep(20 * time.Second)
	Handler.Stop()
	t.Log("end.....")
}
