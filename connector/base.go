package connector

import (
	"context"
	"github.com/tuya/tuya-connector-go/connector/constant"
	"github.com/tuya/tuya-connector-go/connector/env"
	"github.com/tuya/tuya-connector-go/connector/error_proc"
	"github.com/tuya/tuya-connector-go/connector/header"
	"github.com/tuya/tuya-connector-go/connector/httplib"
	"github.com/tuya/tuya-connector-go/connector/logger"
	"github.com/tuya/tuya-connector-go/connector/message"
	"github.com/tuya/tuya-connector-go/connector/sign"
	"github.com/tuya/tuya-connector-go/connector/token"
	"net/http"
)

type ParamFunc func(*httplib.ProxyHttp)

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
func makeRequest(ctx context.Context, params ...ParamFunc) error {
	defer func() {
		if v := recover(); v != nil {
			logger.Log.Errorf("unknown error, err=%+v", v)
		}
	}()
	ph := httplib.NewProxyHttp()
	for _, p := range params {
		p(ph)
	}
	ctx = context.WithValue(ctx, constant.REQ_INFO, ph.GetReqHandler())
	// set header
	ph.SetHeader(header.Handler.GetHeader(ctx))
	//get req
	return ph.DoRequest(ctx)
}

// GET
func MakeGetRequest(ctx context.Context, params ...ParamFunc) error {
	params = append(params, withMethod(http.MethodGet))
	return makeRequest(ctx, params...)
}

// POST
func MakePostRequest(ctx context.Context, params ...ParamFunc) error {
	params = append(params, withMethod(http.MethodPost))
	return makeRequest(ctx, params...)
}

// PUT
func MakePutRequest(ctx context.Context, params ...ParamFunc) error {
	params = append(params, withMethod(http.MethodPut))
	return makeRequest(ctx, params...)
}

// DELETE
func MakeDeleteRequest(ctx context.Context, params ...ParamFunc) error {
	params = append(params, withMethod(http.MethodDelete))
	return makeRequest(ctx, params...)
}

// PATCH
func MakePatchRequest(ctx context.Context, params ...ParamFunc) error {
	params = append(params, withMethod(http.MethodPatch))
	return makeRequest(ctx, params...)
}

// HEAD
func MakeHeadRequest(ctx context.Context, params ...ParamFunc) error {
	params = append(params, withMethod(http.MethodHead))
	return makeRequest(ctx, params...)
}

func WithHeader(h map[string]string) ParamFunc {
	return func(v *httplib.ProxyHttp) {
		v.SetHeader(h)
	}
}

func withMethod(method string) ParamFunc {
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
