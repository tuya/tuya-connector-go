package extension

import "context"

// token interface
// implemented this interface and supports custom token manage
type IToken interface {
	GetToken(ctx context.Context) (string, error)
	GetRefreshToken(ctx context.Context) (string, error)
}

var (
	tokens = make(map[string]IToken)
)

func SetToken(name string, v func() IToken) {
	tokens[name] = v()
}

func GetToken(name string) IToken {
	if tokens[name] == nil {
		panic("token for " + name + " is not existing, make sure you have import the package.")
	}
	return tokens[name]
}
