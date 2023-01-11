package events

import "testing"

const EventType EventId = "event a"

func TestMakesEvents(test *testing.T) {
    event := Make(EventType, "some data")

    if event.id != EventType {
        test.Fatal("The id should be `EventType` but found: " + event.id)
    }
    if event.data != "some data" {
        test.Fatal("The event is not holding the data `some data`")
    }
}

func TestReturnsEventData(test *testing.T) {
    event := Make(EventType, 7)

    data, ok := event.Get().(int)
    if !ok {
        test.Fatal("The event should contain an int")
    }
    if data != 7 {
        test.Fatal("The data should be `7`")
    }
}
