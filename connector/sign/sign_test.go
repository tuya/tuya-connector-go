package sign

import (
	"context"
	"github.com/tuya/tuya-connector-go/connector/env"
	"os"
	"testing"
)

func TestSign(t *testing.T) {
	sw := &signWrapper{}
	sw.ak = os.Getenv(env.TUYA_ACCESSID)
	sw.sk = os.Getenv(env.TUYA_ACCESSKEY)
	sw.token = "123"
	sw.ts = "123"
	t.Log(sw.GetSign(context.Background()))
}
