package message

type MessageBus struct {
	data chan []byte
}

var defaultMessageBus *MessageBus

func init() {
	defaultMessageBus = &MessageBus{
		data: make(chan []byte, 1000),
	}
}

func (bus *MessageBus) push(rawData []byte) {
	bus.data <- rawData
}

func Push(rawData []byte) {
	defaultMessageBus.push(rawData)
}

func Data() <-chan []byte {
	return defaultMessageBus.data
}
