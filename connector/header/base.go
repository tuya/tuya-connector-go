package header

import (
	"context"
	"github.com/tuya/tuya-connector-go/connector/env"
	"github.com/tuya/tuya-connector-go/connector/logger"
	"github.com/tuya/tuya-connector-go/connector/sign"
	"github.com/tuya/tuya-connector-go/connector/token"
	"github.com/tuya/tuya-connector-go/connector/utils"
)

// header interface
type IHeader interface {
	GetHeader(ctx context.Context) map[string]string
}

type headerWrapper struct {
	ak    string
	sk    string
	token string
	ts    string
}

var Handler IHeader

func NewHeaderWrapper() IHeader {
	if env.Config.GetHeaderHandler() != nil {
		return env.Config.GetHeaderHandler().(IHeader)
	}
	return &headerWrapper{}
}

func (t *headerWrapper) GetHeader(ctx context.Context) map[string]string {
	m := make(map[string]string)
	m["Content-Type"] = "application/json"
	m["sign_method"] = "HMAC-SHA256"
	m["client_id"] = env.Config.GetAccessID()

	var token, err = token.Handler.GetToken(ctx)
	if err != nil {
		logger.Log.Errorf("[GetHeader] get token err: %s", err.Error())
		return nil
	}
	m["access_token"] = token

	ts := utils.IntToStr(utils.Microstamp())
	m["t"] = ts

	ctx = context.WithValue(ctx, sign.TOKEN, token)
	ctx = context.WithValue(ctx, sign.TS, ts)
	signStr := sign.Handler.GetSign(ctx)
	m["sign"] = signStr
	return m
}
