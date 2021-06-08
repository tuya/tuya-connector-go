package connector

import (
	"context"
	"github.com/tuya/tuya-connector-go/connector/env"
	"github.com/tuya/tuya-connector-go/connector/error_proc"
	"github.com/tuya/tuya-connector-go/connector/header"
	"github.com/tuya/tuya-connector-go/connector/httplib"
	"github.com/tuya/tuya-connector-go/connector/logger"
	"github.com/tuya/tuya-connector-go/connector/message"
	"github.com/tuya/tuya-connector-go/connector/sign"
	"github.com/tuya/tuya-connector-go/connector/token"
)

type ParamFunc func(v *httplib.ProxyHttp)

// init env config
// init handler
func InitWithOptions(opts ...env.OptionFunc) {
	env.Config = env.NewEnv()
	for _, v := range opts {
		v(env.Config)
	}
	env.Config.Init()
	//if the sign handle is nil, use default handle
	if sign.Handler == nil {
		sign.Handler = sign.NewSignWrapper()
	}
	// check log whether create
	if logger.Log == nil {
		logger.Log = logger.NewDefaultLogger(env.Config.GetAppName(), env.Config.DebugMode())
	}
	if token.Handler == nil {
		token.Handler = token.NewTokenWrapper()
	}
	if header.Handler == nil {
		header.Handler = header.NewHeaderWrapper()
	}
	if message.Handler == nil {
		message.Handler = message.NewEventMsgWrapper()
	}
	logger.Log.Info("### iot core init success ###")
}

// make request api
func MakeRequest(ctx context.Context, params ...ParamFunc) error {
	defer func() {
		if v := recover(); v != nil {
			logger.Log.Errorf("unknown error, err=%+v", v)
		}
	}()
	ph := httplib.NewProxyHttp()
	for _, p := range params {
		p(ph)
	}
	// set header
	if ph.GetProxyHeader() == nil {
		ph.SetHeader(header.Handler.GetHeader(ctx))
	}
	//get req
	return ph.DoRequest(ctx)
}

func WithHeader(h map[string]string) ParamFunc {
	return func(v *httplib.ProxyHttp) {
		v.SetHeader(h)
	}
}

func WithMethod(method string) ParamFunc {
	return func(v *httplib.ProxyHttp) {
		v.SetMethod(method)
	}
}

func WithAPIUri(uri string) ParamFunc {
	return func(v *httplib.ProxyHttp) {
		v.SetAPIUri(env.Config.GetApiHost() + uri)
	}
}

func WithPayload(body []byte) ParamFunc {
	return func(v *httplib.ProxyHttp) {
		v.SetPayload(body)
	}
}

func WithResp(res interface{}) ParamFunc {
	return func(v *httplib.ProxyHttp) {
		v.SetResp(res)
	}
}

func WithErrProc(code int, f error_proc.IError) ParamFunc {
	return func(v *httplib.ProxyHttp) {
		v.SetErrProc(code, f)
	}
}
