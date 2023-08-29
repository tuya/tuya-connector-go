package token

import (
	"context"
	"fmt"
	"net/http"

	"github.com/callmegema/tuya-connector-go/connector/constant"
	"github.com/callmegema/tuya-connector-go/connector/env"
	"github.com/callmegema/tuya-connector-go/connector/env/extension"
	"github.com/callmegema/tuya-connector-go/connector/httplib"
	"github.com/callmegema/tuya-connector-go/connector/logger"
	"github.com/callmegema/tuya-connector-go/connector/utils"
)

const (
	GET_TOKEN_REQ_URI   = "/v1.0/token?grant_type=1"
	GET_REFRESH_REQ_URI = "/v1.0/token/%s"
)

type tokenAPIResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	T       int64  `json:"t"`
	Result  struct {
		ExpireTime   int    `json:"expire_time"`
		UID          string `json:"uid"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	} `json:"result"`
}

// get token
func (t *token) fromAPIGetToken(ctx context.Context) (*tokenAPIResponse, error) {
	return t.commonReqToken(ctx, GET_TOKEN_REQ_URI)
}

// refresh token
func (t *token) fromAPIRefreshToken(ctx context.Context) (*tokenAPIResponse, error) {
	return t.commonReqToken(ctx, fmt.Sprintf(GET_REFRESH_REQ_URI, t.reToken))
}

func (t *token) commonReqToken(ctx context.Context, uri string) (*tokenAPIResponse, error) {
	resp := &tokenAPIResponse{}
	th := httplib.NewProxyHttp()
	th.SetMethod(http.MethodGet)
	th.SetAPIUri(env.Config.GetApiHost() + uri)
	th.SetResp(resp)
	ts := utils.IntToStr(utils.Microstamp())
	nonce := utils.GetUUID()
	ctx = context.WithValue(ctx, constant.REQ_INFO, th.GetReqHandler())
	ctx = context.WithValue(ctx, constant.TOKEN, "")
	ctx = context.WithValue(ctx, constant.TS, ts)
	ctx = context.WithValue(ctx, constant.NONCE, nonce)
	signStr := extension.GetSign(constant.TUYA_SIGN).Sign(ctx)
	th.SetHeader(map[string]string{
		constant.Header_ContentType: constant.ContentType_JSON,
		constant.Header_SignMethod:  constant.SignMethod_HMAC,
		constant.Header_ClientID:    env.Config.GetAccessID(),
		constant.Header_TimeStamp:   ts,
		constant.Header_Sign:        signStr,
		constant.Header_Nonce:       nonce,
	})
	err := th.DoRequest(ctx)
	if err != nil {
		logger.Log.Infof("[commonReqToken] req failed, req:%v, err:%v", th.GetReqHandler(), err)
		return nil, err
	}
	t.setToken(resp.Result.AccessToken, resp.Result.RefreshToken, resp.Result.ExpireTime)
	return resp, err
}
