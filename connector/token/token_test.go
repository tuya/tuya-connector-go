package token

import (
	"context"
	"fmt"
	"github.com/tuya/tuya-connector-go/connector/env"
	"github.com/tuya/tuya-connector-go/connector/logger"
	"github.com/tuya/tuya-connector-go/connector/sign"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("init....")
	env.Config = env.NewEnv()
	env.Config.Init()
	if Handler == nil {
		Handler = NewTokenWrapper()
	}
	if sign.Handler == nil {
		sign.Handler = sign.NewSignWrapper()
	}
	if logger.Log == nil {
		logger.Log = logger.NewDefaultLogger(env.Config.GetAppName(), env.Config.DebugMode())
	}
	fmt.Println("### iot core init success ###")
	m.Run()
}

func TestToken(t *testing.T) {
	tk, err := Handler.GetToken(context.Background())
	t.Log(tk, err)
}
