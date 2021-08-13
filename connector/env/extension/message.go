package extension

type IEventMessage interface {
	// init event message client
	InitMessageClient()
	// listen event message
	SubEventMessage(f interface{})
	// stop event message receive
	Stop()
}

var (
	messages = make(map[string]IEventMessage)
)

func SetMessage(name string, v func() IEventMessage) {
	messages[name] = v()
}

func GetMessage(name string) IEventMessage {
	if messages[name] == nil {
		panic("message for " + name + " is not existing, make sure you have import the package.")
	}
	return messages[name]
}
