package sign

import (
	"bytes"
	"context"
	"github.com/tuya/tuya-connector-go/connector/constant"
	"github.com/tuya/tuya-connector-go/connector/env"
	"github.com/tuya/tuya-connector-go/connector/logger"
	"github.com/tuya/tuya-connector-go/connector/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

// sign interface
// implemented this interface and supports custom signatures
type ISign interface {
	GetSign(ctx context.Context) string
}

type signWrapper struct {
	ak           string
	sk           string
	token        string
	ts           string
	nonce        string
	stringToSign string
}

var Handler ISign

func NewSignWrapper() ISign {
	if env.Config.GetSignHandler() != nil {
		return env.Config.GetSignHandler().(ISign)
	}
	return &signWrapper{
		ak: env.Config.GetAccessID(),
		sk: env.Config.GetAccessKey(),
	}
}

// No need to pass the token parameter when getting the token
func (t *signWrapper) GetSign(ctx context.Context) string {
	t.token = ctx.Value(constant.TOKEN).(string)
	t.ts = ctx.Value(constant.TS).(string)
	t.nonce = ctx.Value(constant.NONCE).(string)
	t.stringToSign = t.calStringToSign(ctx)
	sign := utils.HS256Sign(t.sk, t.ak+t.token+t.ts+t.nonce+t.stringToSign)
	return strings.ToUpper(sign)
}

func (t *signWrapper) calStringToSign(ctx context.Context) string {
	req, ok := ctx.Value(constant.REQ_INFO).(*http.Request)
	if !ok {
		return ""
	}
	contentSha256 := ""
	if req.Body != nil {
		buf, _ := ioutil.ReadAll(req.Body)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
		contentSha256 = utils.GetSha256(buf)
	} else {
		contentSha256 = utils.GetSha256([]byte(""))
	}

	headers := ""
	signHeaderKeys := req.Header.Get(constant.Signature_Headers)
	if signHeaderKeys != "" {
		keys := strings.Split(signHeaderKeys, ":")
		for _, key := range keys {
			headers += key + ":" + req.Header.Get(key) + "\n"
		}
	}

	uri := req.URL.Path
	keys := make([]string, 0, 10)
	form, err := url.ParseQuery(req.URL.RawQuery)
	if err == nil {
		for key, _ := range form {
			keys = append(keys, key)
		}
	}
	if len(keys) > 0 {
		uri += "?"
		sort.Strings(keys)
		for _, keyName := range keys {
			value := form.Get(keyName)
			uri += keyName + "=" + value + "&"
		}
		uri = strings.TrimSuffix(uri, "&")
	}

	stringToSign := req.Method + "\n" + contentSha256 + "\n" + headers + "\n" + uri
	logger.Log.Debugf("[calStringToSign] httpMethod=%s, contentSha256=%s, headers=%s, uri=%s, stringToSign=%s",
		req.Method, contentSha256, headers, uri, stringToSign)

	return stringToSign
}
