package header

import (
	"context"
	"github.com/tuya/tuya-connector-go/connector/constant"
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
	m[constant.Header_ContentType] = constant.ContentType_JSON
	m[constant.Header_SignMethod] = constant.SignMethod_HMAC
	m[constant.Header_DevChannel] = constant.Dev_Channel
	m[constant.Header_DevLang] = constant.Dev_Lang
	m[constant.Header_ClientID] = env.Config.GetAccessID()
	nonce := utils.GetUUID()
	m[constant.Header_Nonce] = nonce
	var token, err = token.Handler.GetToken(ctx)
	if err != nil {
		logger.Log.Errorf("[GetHeader] get token err: %s", err.Error())
		return nil
	}
	m[constant.Header_AccessToken] = token

	ts := utils.IntToStr(utils.Microstamp())
	m[constant.Header_TimeStamp] = ts

	ctx = context.WithValue(ctx, constant.TOKEN, token)
	ctx = context.WithValue(ctx, constant.TS, ts)
	ctx = context.WithValue(ctx, constant.NONCE, nonce)
	signStr := sign.Handler.GetSign(ctx)
	m[constant.Header_Sign] = signStr
	return m
}
