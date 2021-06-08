package error_proc

import "context"

type IError interface {
	Process(ctx context.Context, code int, msg string)
}
