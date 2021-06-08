package token

import (
	"context"
	"fmt"
	"github.com/tuya/tuya-connector-go/connector/env"
	"github.com/tuya/tuya-connector-go/connector/error_proc"
	"github.com/tuya/tuya-connector-go/connector/httplib"
	"github.com/tuya/tuya-connector-go/connector/logger"
	"github.com/tuya/tuya-connector-go/connector/sign"
	"github.com/tuya/tuya-connector-go/connector/utils"
	"net/http"
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
	th.SetErrProc(error_proc.TOKEN_EXPIRED, &tokenError{})
	ts := utils.IntToStr(utils.Microstamp())
	ctx = context.WithValue(ctx, sign.TOKEN, "")
	ctx = context.WithValue(ctx, sign.TS, ts)
	signStr := sign.Handler.GetSign(ctx)
	th.SetHeader(map[string]string{
		"Content-Type": "application/json",
		"sign_method":  "HMAC-SHA256",
		"client_id":    env.Config.GetAccessID(),
		"t":            ts,
		"sign":         signStr,
	})
	err := th.DoRequest(ctx)
	if err != nil {
		logger.Log.Infof("[commonReqToken] req failed, req:%v, err:%v", th.GetReqHandler(), err)
		return nil, err
	}
	t.setToken(resp.Result.AccessToken, resp.Result.RefreshToken, resp.Result.ExpireTime)
	return resp, err
}

type tokenError struct {
}

func (t *tokenError) Process(ctx context.Context, code int, msg string) {
	if code == error_proc.TOKEN_EXPIRED {
		_, _ = Handler.GetRefreshToken(ctx)
	}
}