package messaging

import (
	"github.com/tuya/tuya-connector-go/connector/logger"
	"github.com/tuya/tuya-connector-go/connector/message"
	"github.com/tuya/tuya-connector-go/connector/message/event"
)

func Listener() {
	message.Handler.InitMessageClient()

	message.Handler.SubEventMessage(func(m *event.NameUpdateMessage) {
		logger.Log.Info("=========== name update： ==========")
		logger.Log.Info(m)
	})

	message.Handler.SubEventMessage(func(m *event.StatusReportMessage) {
		logger.Log.Info("=========== report data： ==========")
		for _, v := range m.Status {
			logger.Log.Infof(v.Code, v.Value)
		}
	})
}
