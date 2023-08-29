package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/callmegema/tuya-connector-go/connector"
	"github.com/callmegema/tuya-connector-go/connector/constant"
	"github.com/callmegema/tuya-connector-go/connector/env/extension"
	"github.com/callmegema/tuya-connector-go/connector/logger"
	"github.com/callmegema/tuya-connector-go/example/messaging"
	"github.com/callmegema/tuya-connector-go/example/router"
)

func main() {
	// custom init config
	/*connector.InitWithOptions(env.WithApiHost(httplib.URL_CN),
	env.WithMsgHost(httplib.MSG_CN),
	env.WithAccessID(""),
	env.WithAccessKey(""),
	env.WithAppName(""),
	env.WithDebugMode(true))*/

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
			extension.GetMessage(constant.TUYA_MESSAGE).Stop()
			logger.Log.Infof("receive sig:%v, shutdown the http server...", c.String())
			return
		}
	}
}
