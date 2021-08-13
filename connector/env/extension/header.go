package extension

import "context"

type IError interface {
	Process(ctx context.Context, code int, msg string)
}

// header interface
type IHeader interface {
	GetHeader(ctx context.Context) map[string]string
}

var (
	headers = make(map[string]IHeader)
)

func SetHeader(name string, v func() IHeader) {
	headers[name] = v()
}

func GetHeader(name string) IHeader {
	if headers[name] == nil {
		panic("header for " + name + " is not existing, make sure you have import the package.")
	}
	return headers[name]
}
