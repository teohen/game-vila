package events

type EventType int

const (
	EventTreeCut EventType = iota
)

type GameEvent struct {
	Type    EventType
	Payload map[string]interface{}
}

var EventQueue = make(chan GameEvent, 100)

func Emit(evt GameEvent) {
	EventQueue <- evt
}
