package token

import (
	"context"
	"fmt"
	"testing"

	"github.com/callmegema/tuya-connector-go/connector/constant"
	"github.com/callmegema/tuya-connector-go/connector/env"
	"github.com/callmegema/tuya-connector-go/connector/env/extension"
	"github.com/callmegema/tuya-connector-go/connector/logger"
	"github.com/callmegema/tuya-connector-go/connector/sign"
)

func TestMain(m *testing.M) {
	fmt.Println("init....")
	env.Config = env.NewEnv()
	env.Config.Init()
	extension.SetToken(constant.TUYA_TOKEN, newTokenInstance)
	extension.SetSign(constant.TUYA_SIGN, sign.NewSignWrapper)
	if logger.Log == nil {
		logger.Log = logger.NewDefaultLogger(env.Config.GetAppName(), env.Config.DebugMode())
	}
	fmt.Println("### iot core init success ###")
	m.Run()
}

func TestToken(t *testing.T) {
	tk, err := extension.GetToken(constant.TUYA_TOKEN).Do(context.Background())
	t.Log(tk, err)
}
