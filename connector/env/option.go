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

// set custom log wrapper
func WithLogWrapper(v interface{}) OptionFunc {
	return func(e *environment) {
		e.logHandler = v
	}
}

// set custom token wrapper
func WithTokenWrapper(v interface{}) OptionFunc {
	return func(e *environment) {
		e.tokenHandler = v
	}
}

// set custom sign wrapper
func WithSignWrapper(v interface{}) OptionFunc {
	return func(e *environment) {
		e.signHandler = v
	}
}

// set custom header wrapper
func WithHeaderWrapper(v interface{}) OptionFunc {
	return func(e *environment) {
		e.headerHandler = v
	}
}

// set custom event message wrapper
func WithEventMsgWrapper(v interface{}) OptionFunc {
	return func(e *environment) {
		e.eventHandler = v
	}
}
