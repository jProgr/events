package events

import "testing"

const (
    EventA EventId = "event a"
    EventB EventId = "event b"
)

type dto struct {
    data int
}

type dto2 struct {
    listenerA bool
    listenerB bool
}

func TestBuildsNewDefaultDispatcher(test *testing.T) {
    dispatcher, err := NewDispatcher()

    if err != nil {
        test.Fatal("There shouldn't be an error when creating a default dispatcher")
    }
    if dispatcher.listeners == nil {
        test.Fatal("`dispatcher.listeners` should be a map")
    }

    config := dispatcher.config
    if config.isAsync {
        test.Fatal("The default dispatcher should not be async")
    }
    if config.isFacade {
        test.Fatal("The default dispatcher should not be a facade")
    }
}

func TestBuildsNewAsyncDispatcher(test *testing.T) {
    dispatcher, err := NewDispatcher(func(config *Config) {
        config.ShouldAsync(true)
    })

    if err != nil {
        test.Fatal("An async dispatcher shouldn't produce an error")
    }
    if !dispatcher.config.isAsync {
        test.Fatal("This dispatcher should be async")
    }
}

func TestBuildsNewFacadeDispatcher(test *testing.T) {
    dispatcher, err := NewDispatcher(func(config *Config) {
        config.AsFacade(true)
    })

    if err != nil {
        test.Fatal("A facade dispatcher shouldn't produce an error")
    }
    if !dispatcher.config.isFacade {
        test.Fatal("This dispatcher should be a facade")
    }
    if facade != dispatcher {
        test.Fatal("The dispatcher should be registered as global facade")
    }
}

func TestFailsToBuildADispatcherOnWrongConfig(test *testing.T) {
    _, err := NewDispatcher(func(config *Config) {
        config.ShouldWait(false, nil)
    })

    if err == nil {
        test.Fatal("There should be an error when creating a dispatcher that does not wait and has no sync.WaitGroup")
    }
    if _, ok := err.(*AsyncConfigError); !ok {
        test.Fatal("The error should be of type AsyncConfigError")
    }
}

func TestRegistersEvents(test *testing.T) {
    listener := func(_ Event) {}
    dispatcher, _ := NewDispatcher()
    dispatcher.Register(EventA, listener)

    listeners, ok := dispatcher.listeners[EventA]
    if !ok {
        test.Fatal("EventA should be registered in dispatcher")
    }
    if len(listeners) != 1 {
        test.Fatal("There should be exactly one listener registered")
    }

    dispatcher.Register(EventA, listener)

    listeners = dispatcher.listeners[EventA]
    if len(listeners) != 2 {
        test.Fatal("There should be exactly two listeners registered")
    }
}

func TestTriggersListeners(test *testing.T) {
    listener := func(event Event) {
        dto := event.Get().(*dto)
        dto.data = dto.data + 1
    }
    dispatcher, _ := NewDispatcher()
    dispatcher.
        Register(EventA, listener).
        Register(EventB, listener)

    eventData := &dto{7}
    dispatcher.Dispatch(Make(EventA, eventData))
    if eventData.data != 8 {
        test.Fatal("Listeners of EventA weren't dispatched")
    }

    dispatcher.Dispatch(Make(EventB, eventData))
    if eventData.data != 9 {
        test.Fatal("Listeners of EventB weren't dispatched")
    }
}

func TestTriggersListenersAsync(test *testing.T) {
    listenerA := func(event Event) {
        dto := event.Get().(*dto2)
        dto.listenerA = true
    }
    listenerB := func(event Event) {
        dto := event.Get().(*dto2)
        dto.listenerB = true
    }

    dispatcher, _ := NewDispatcher(func(config *Config) {
        config.ShouldAsync(true)
    })
    dispatcher.
        Register(EventA, listenerA).
        Register(EventA, listenerB)

    dto := &dto2{
        listenerA: false,
        listenerB: false,
    }
    dispatcher.Dispatch(Make(EventA, dto))

    if !dto.listenerA {
        test.Fatal("`listenerA` wasn't dispatched")
    }
    if !dto.listenerB {
        test.Fatal("`listenerB` wasn't dispatched")
    }
}

func TestRegistersEventsAsFacade(test *testing.T) {
    NewDispatcher(func(config *Config) {
        config.AsFacade(true)
    })
    listener := func(_ Event) {}

    Register(EventA, listener)

    listeners, ok := facade.listeners[EventA]
    if !ok {
        test.Fatal("EventA should be registered in dispatcher")
    }
    if len(listeners) != 1 {
        test.Fatal("There should be exactly one listener registered")
    }

    Register(EventA, listener)

    listeners = facade.listeners[EventA]
    if len(listeners) != 2 {
        test.Fatal("There should be exactly two listeners registered")
    }
}

func TestPanicsOnRegisteringEventsAsFacadeWhenNotAFacade(test *testing.T) {
    defer recoverAndCheckPanic(test)

    resetFacade()
    NewDispatcher()

    Register(EventB, func(_ Event) {})

    test.Fatal("This shouldn't be executed")
}

func TestTriggersListenersAsFacade(test *testing.T) {
    NewDispatcher(func(config *Config) {
        config.AsFacade(true)
    })
    Register(EventB, func(event Event) {
        dto := event.Get().(*dto)
        dto.data = dto.data + 1
    })

    eventData := &dto{7}
    Dispatch(Make(EventB, eventData))
    if eventData.data != 8 {
        test.Fatal("Listeners of EventB weren't dispatched")
    }
}

func TestPanicsOnDispatchingEventsAsFacadeWhenNotAFacade(test *testing.T) {
    defer recoverAndCheckPanic(test)

    resetFacade()
    NewDispatcher()

    Dispatch(Make(EventB, "some data"))

    test.Fatal("This shouldn't be executed")
}

func resetFacade() {
    facade = nil
}

func recoverAndCheckPanic(test *testing.T) {
    err, ok := recover().(string)
    if !ok {
        test.Fatal("It should return a string")
    }
    if err != "No facade registered" {
        test.Fatal("Wrong error. Actual error: " + err)
    }
}
