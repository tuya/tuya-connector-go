package token

import (
	"context"
	"fmt"
	"github.com/tuya/tuya-connector-go/connector/constant"
	"github.com/tuya/tuya-connector-go/connector/env/extension"
	"sync"
	"time"
)

func init() {
	extension.SetToken(constant.TUYA_TOKEN, newTokenInstance)
	fmt.Println("init token extension......")
}

func newTokenInstance() extension.IToken {
	return NewTokenWrapper()
}

type token struct {
	mu       *sync.RWMutex
	token    string
	reToken  string
	expireAt time.Time
}

func NewTokenWrapper() extension.IToken {
	return &token{
		mu: &sync.RWMutex{},
	}
}

func (t *token) Do(ctx context.Context) (string, error) {
	t.mu.RLock()
	tk := t.token
	ttl := t.expireAt
	t.mu.RUnlock()
	if tk != "" && ttl.After(time.Now()) {
		return tk, nil
	}
	tokenCtx := context.Background()
	exeCnt := ctx.Value(constant.ExeCount)
	if exeCnt != nil && exeCnt.(int) > 0 {
		tokenCtx = context.WithValue(tokenCtx, constant.ExeCount, exeCnt)
	}
	//if tk == "" || t.reToken == "" {
	_, err := t.fromAPIGetToken(tokenCtx)
	if err != nil {
		return "", err
	}
	/*} else {
		_, err := t.fromAPIRefreshToken(tokenCtx)
		if err != nil {
			return "", err
		}
	}*/
	t.mu.RLock()
	tk = t.token
	t.mu.RUnlock()
	return tk, nil
}

func (t *token) Refresh(ctx context.Context) (string, error) {
	if t.reToken != "" {
		return t.reToken, nil
	}
	_, err := t.Do(ctx)
	if err != nil {
		return "", err
	}
	return t.reToken, nil
}

func (t *token) setToken(token, refreshToken string, expire int) {
	t.mu.Lock()
	t.token = token
	t.reToken = refreshToken
	t.expireAt = time.Now().Add(time.Duration(expire) * time.Second)
	t.mu.Unlock()
}
