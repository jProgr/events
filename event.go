package events

// EventId identifies events. This should be unique in each dispatcher.
type EventId string

// Event wraps an event ID and any data that may be passed to the listeners.
// Denotes that something happened in the codebase.
type Event struct {
    id   EventId
    data any
}

// Make creates a new event. In general, data should be a pointer.
func Make(id EventId, data any) Event {
    return Event{id, data}
}

// Get returns the data contained in the event. A type assertion is needed due
// to the return type being any.
func (event Event) Get() any {
    return event.data
}

// Listener is a function that has as only argument an Event and returns nothing.
// Contains the logic that should be run when an event is triggered.
type Listener func(Event)
