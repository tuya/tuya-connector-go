package header

import (
	"context"
	"fmt"

	"github.com/callmegema/tuya-connector-go/connector/constant"
	"github.com/callmegema/tuya-connector-go/connector/env"
	"github.com/callmegema/tuya-connector-go/connector/env/extension"
	"github.com/callmegema/tuya-connector-go/connector/logger"
	"github.com/callmegema/tuya-connector-go/connector/utils"
)

func init() {
	extension.SetHeader(constant.TUYA_HEADER, newHeaderInstance)
	fmt.Println("init header extension......")
}

func newHeaderInstance() extension.IHeader {
	return NewHeaderWrapper()
}

type headerWrapper struct {
}

func NewHeaderWrapper() extension.IHeader {
	return &headerWrapper{}
}

func (t *headerWrapper) Do(ctx context.Context) map[string]string {
	m := make(map[string]string)
	m[constant.Header_ContentType] = constant.ContentType_JSON
	m[constant.Header_SignMethod] = constant.SignMethod_HMAC
	m[constant.Header_DevChannel] = constant.Dev_Channel
	m[constant.Header_DevLang] = constant.Dev_Lang
	m[constant.Header_ClientID] = env.Config.GetAccessID()
	nonce := utils.GetUUID()
	m[constant.Header_Nonce] = nonce
	var token, err = extension.GetToken(constant.TUYA_TOKEN).Do(ctx)
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
	signStr := extension.GetSign(constant.TUYA_SIGN).Sign(ctx)
	m[constant.Header_Sign] = signStr
	return m
}
