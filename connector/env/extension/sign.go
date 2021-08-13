package extension

import "context"

// sign interface
// implemented this interface and supports custom signatures
type ISign interface {
	Sign(ctx context.Context) string
}

var (
	signs = make(map[string]ISign)
)

func SetSign(name string, v func() ISign) {
	signs[name] = v()
}

func GetSign(name string) ISign {
	if signs[name] == nil {
		panic("sign for " + name + " is not existing, make sure you have import the package.")
	}
	return signs[name]
}
