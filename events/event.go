package events

// Tipo de evento (usando enum/int para performance)
type EventType int

const (
	EventTreeCut EventType = iota
	// EventVillagerHungry
	// EventResourceDeposited
)

// Estrutura do evento com dados genéricos ou específicos
type GameEvent struct {
	Type    EventType
	Payload map[string]interface{}
}

// O canal global de eventos (pode ser buffered para não travar o loop)
var EventQueue = make(chan GameEvent, 100)

// TODO: create Lots of kind on events
func Emit(evt GameEvent) {
	EventQueue <- evt
}
