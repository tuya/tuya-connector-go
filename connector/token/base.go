package token

import (
	"context"
	"github.com/tuya/tuya-connector-go/connector/constant"
	"github.com/tuya/tuya-connector-go/connector/env"
	"sync"
	"time"
)

// token interface
// implemented this interface and supports custom token manage
type IToken interface {
	GetToken(ctx context.Context) (string, error)
	GetRefreshToken(ctx context.Context) (string, error)
}

var Handler IToken

type token struct {
	mu       *sync.RWMutex
	token    string
	reToken  string
	expireAt time.Time
}

func NewTokenWrapper() IToken {
	if env.Config.GetTokenHandler() != nil {
		return env.Config.GetTokenHandler().(IToken)
	}
	return &token{
		mu: &sync.RWMutex{},
	}
}

func (t *token) GetToken(ctx context.Context) (string, error) {
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
	if tk == "" || t.reToken == "" {
		_, err := t.fromAPIGetToken(tokenCtx)
		if err != nil {
			return "", err
		}
	} else {
		_, err := t.fromAPIRefreshToken(tokenCtx)
		if err != nil {
			return "", err
		}
	}
	t.mu.RLock()
	tk = t.token
	t.mu.RUnlock()
	return tk, nil
}

func (t *token) GetRefreshToken(ctx context.Context) (string, error) {
	if t.reToken != "" {
		return t.reToken, nil
	}
	_, err := t.GetToken(ctx)
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
