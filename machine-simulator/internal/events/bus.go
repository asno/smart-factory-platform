package events

type Event struct {
	Type    string
	Payload interface{}
}

type Bus struct {
	Events chan Event
}

func NewBus(bufferSize int) *Bus {
	return &Bus{
		Events: make(chan Event, bufferSize),
	}
}