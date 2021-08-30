package httplib

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/tuya/tuya-connector-go/connector/constant"
	"github.com/tuya/tuya-connector-go/connector/env/extension"
	"github.com/tuya/tuya-connector-go/connector/logger"
	"io/ioutil"
	"net/http"
	"net/url"
)

type ProxyHttp struct {
	header  map[string]string
	method  string
	apiUri  string
	payload []byte
	req     *http.Request
	resp    interface{}
	//mu        *sync.RWMutex
	errMap map[int]extension.IError
}

func NewProxyHttp() *ProxyHttp {
	return &ProxyHttp{
		header: make(map[string]string),
		req: &http.Request{
			Header: make(http.Header),
		},
		errMap: make(map[int]extension.IError),
	}
}

func (t *ProxyHttp) SetHeader(h map[string]string) {
	for k, v := range h {
		t.req.Header.Set(k, v)
		t.header[k] = v
	}
}

func (t *ProxyHttp) SetMethod(v string) {
	t.method = v
	t.req.Method = v
}

func (t *ProxyHttp) SetAPIUri(v string) {
	t.apiUri = v
	u, err := url.Parse(v)
	if err != nil {
		logger.Log.Errorf("[SetAPIUri] set uri err: %s", err.Error())
	}
	t.req.URL = u
}

func (t *ProxyHttp) SetPayload(v []byte) {
	t.payload = v
	t.req.Body = ioutil.NopCloser(bytes.NewBuffer(v))
}

func (t *ProxyHttp) GetPayload() []byte {
	return t.payload
}

func (t *ProxyHttp) SetResp(v interface{}) {
	t.resp = v
}

func (t *ProxyHttp) SetErrProc(code int, v extension.IError) {
	t.errMap[code] = v
}

func (t *ProxyHttp) GetProxyHeader() map[string]string {
	return t.header
}

func (t *ProxyHttp) GetReqHandler() *http.Request {
	return t.req
}

func (t *ProxyHttp) DoRequest(ctx context.Context) error {
	var err error
	resp, err := http.DefaultClient.Do(t.req)
	if err != nil {
		logger.Log.Errorf("[ProxyHttp] do req failed err:%v, req:%v", err.Error(), t.req)
		return err
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Errorf("[ProxyHttp] do req failed err:%v, req:%v", err.Error(), t.req)
		return err
	}
	var rst response
	err = json.Unmarshal(bs, &rst)
	if err != nil {
		logger.Log.Errorf("[ProxyHttp] do req failed err:%v, req:%v, resp:%v", err.Error(), t.req, string(bs))
		return err
	}
	if !rst.Success {
		logger.Log.Errorf("[ProxyHttp] do req failed req:%v, resp:%v", t.req, string(bs))
		if f, ok := t.errMap[rst.Code]; ok {
			// avoid loop
			exeCnt := ctx.Value(constant.ExeCount)
			if exeCnt != nil && exeCnt.(int) > 0 {
				return errors.New(rst.Msg)
			}
			ctx = context.WithValue(ctx, constant.ExeCount, 1)
			f.Process(ctx, rst.Code, rst.Msg)
			if rst.Code == constant.TOKEN_EXPIRED {
				return nil
			}
			return errors.New(rst.Msg)
		}
	}
	err = json.Unmarshal(bs, &t.resp)
	if err != nil {
		logger.Log.Errorf("[ProxyHttp] do req failed err:%v, req:%v, resp:%v", err.Error(), t.req, string(bs))
		return err
	}
	logger.Log.Infof("[ProxyHttp] success req:%v, resp:%+v", t.req, t.resp)
	return nil
}
