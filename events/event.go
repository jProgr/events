package events

type EventId string

type Event struct {
    id   EventId
    data any
}

func Make(id EventId, data any) Event {
    return Event{id, data}
}

func (event Event) Get() any {
    return event.data
}

type Listener func(Event)
