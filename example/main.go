package main

import (
	"github.com/tuya/tuya-connector-go/connector"
	"github.com/tuya/tuya-connector-go/connector/logger"
	"github.com/tuya/tuya-connector-go/connector/message"
	"github.com/tuya/tuya-connector-go/example/messaging"
	"github.com/tuya/tuya-connector-go/example/router"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// custom init config
	/*connector.InitWithOptions(env.WithApiHost(httplib.URL_CN),
	env.WithMsgHost(httplib.MSG_CN),
	env.WithAccessID(""),
	env.WithAccessKey(""),
	env.WithAppName(""),
	env.WithDebugMode(true),
	env.WithLogWrapper(nil),
	env.WithHeaderWrapper(nil),
	env.WithTokenWrapper(nil),
	env.WithSignWrapper(nil),
	env.WithEventMsgWrapper(nil))*/

	// default init config
	connector.InitWithOptions()

	go messaging.Listener()

	r := router.NewGinEngin()
	go r.Run("0.0.0.0:2021")
	watitSignal()
}

func watitSignal() {
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	for {
		select {
		case c := <-quitCh:
			message.Handler.Stop()
			logger.Log.Infof("receive sig:%v, shutdown the http server...", c.String())
			return
		}
	}
}
