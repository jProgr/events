package events

import "sync"

var facade *Dispatcher

// Dispatcher

type Dispatcher struct {
    listeners map[EventId][]Listener
    config    *Config
    waitGroup *sync.WaitGroup
}

func NewDispatcher(configurers ...func(*Config)) *Dispatcher {
    config := newDefaultConfig()
    for _, configurer := range configurers {
        configurer(config)
    }

    waitGroup := config.waitGroup
    if waitGroup == nil {
        waitGroup = new(sync.WaitGroup)
    }

    dispatcher := &Dispatcher{
        listeners: make(map[EventId][]Listener),
        config:    config,
        waitGroup: waitGroup,
    }

    if config.isFacade {
        facade = dispatcher
    }

    return dispatcher
}

func (dispatcher *Dispatcher) Register(id EventId, listener Listener) *Dispatcher {
    if listeners, ok := dispatcher.listeners[id]; ok {
        dispatcher.listeners[id] = append(listeners, listener)
        return dispatcher
    }

    dispatcher.listeners[id] = []Listener{listener}

    return dispatcher
}

func (dispatcher *Dispatcher) Dispatch(event Event) *Dispatcher {
    listeners, ok := dispatcher.listeners[event.id]
    if !ok || len(listeners) == 0 {
        return dispatcher
    }

    for _, listener := range listeners {
        if dispatcher.config.isAsync {
            dispatchAsync(dispatcher, event, listener)
            continue
        }

        listener(event)
    }

    if dispatcher.config.shouldWait {
        dispatcher.waitGroup.Wait()
    }

    return dispatcher
}

func dispatchAsync(dispatcher *Dispatcher, event Event, listener Listener) {
    dispatcher.waitGroup.Add(1)

    go func() {
        defer dispatcher.waitGroup.Done()
        listener(event)
    }()
}

// Interaction as facade

func Register(id EventId, listener Listener) *Dispatcher {
    if facade == nil {
        panic("No facade registered")
    }

    return facade.Register(id, listener)
}

func Dispatch(event Event) *Dispatcher {
    if facade == nil {
        panic("No facade registered")
    }

    return facade.Dispatch(event)
}
