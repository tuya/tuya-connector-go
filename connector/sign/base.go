package sign

import (
	"context"
	"github.com/tuya/tuya-connector-go/connector/env"
	"github.com/tuya/tuya-connector-go/connector/utils"
	"strings"
)

const (
	TOKEN = "token"
	TS    = "ts"
)

// sign interface
// implemented this interface and supports custom signatures
type ISign interface {
	GetSign(ctx context.Context) string
}

type signWrapper struct {
	ak    string
	sk    string
	token string
	ts    string
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
	t.token = ctx.Value(TOKEN).(string)
	t.ts = ctx.Value(TS).(string)
	sign := utils.HS256Sign(t.sk, t.ak+t.token+t.ts)
	return strings.ToUpper(sign)
}
