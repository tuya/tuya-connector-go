package sign

import (
	"context"
	"github.com/tuya/tuya-connector-go/connector/constant"
	"github.com/tuya/tuya-connector-go/connector/env"
	"testing"
)

func TestSign(t *testing.T) {
	env.Config = env.NewEnv()
	env.Config.Init()
	ctx := context.Background()
	ctx = context.WithValue(ctx, constant.TOKEN, "123")
	ctx = context.WithValue(ctx, constant.TS, "123")
	ctx = context.WithValue(ctx, constant.NONCE, "123")
	sw := &signWrapper{}
	t.Log(sw.Sign(ctx))
}
