package env

type OptionFunc func(*environment)

// app name
func WithAppName(v string) OptionFunc {
	return func(e *environment) {
		e.appName = v
	}
}

// call event message domain name
func WithMsgHost(v string) OptionFunc {
	return func(e *environment) {
		e.msgHost = v
	}
}

// call api domain name
func WithApiHost(v string) OptionFunc {
	return func(e *environment) {
		e.apiHost = v
	}
}

// set ak
func WithAccessID(v string) OptionFunc {
	return func(e *environment) {
		e.accessID = v
	}
}

// set sk
func WithAccessKey(v string) OptionFunc {
	return func(e *environment) {
		e.accessKey = v
	}
}

func WithDebugMode(v bool) OptionFunc {
	return func(e *environment) {
		e.isDebug = v
	}
}
